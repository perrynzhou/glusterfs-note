#
# Copyright(c) 2019 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause-Clear
#

import time

import pytest

from storage_devices.disk import DiskType, DiskTypeSet, DiskTypeLowerThan
from test_tools.dd import Dd
from test_utils.os_utils import sync, Udev
from .io_class_common import *


@pytest.mark.require_disk("cache", DiskTypeSet([DiskType.optane, DiskType.nand]))
@pytest.mark.require_disk("core", DiskTypeLowerThan("cache"))
def test_ioclass_process_name():
    """Check if data generated by process with particular name is cached"""
    cache, core = prepare()

    ioclass_id = 1
    dd_size = Size(4, Unit.KibiByte)
    dd_count = 1
    iterations = 100

    ioclass_config.add_ioclass(
        ioclass_id=ioclass_id,
        eviction_priority=1,
        allocation=True,
        rule=f"process_name:dd&done",
        ioclass_config_path=ioclass_config_path,
    )
    casadm.load_io_classes(cache_id=cache.cache_id, file=ioclass_config_path)

    cache.flush_cache()

    Udev.disable()

    TestRun.LOGGER.info(f"Check if all data generated by dd process is cached.")
    for i in range(iterations):
        dd = (
            Dd()
            .input("/dev/zero")
            .output(core.system_path)
            .count(dd_count)
            .block_size(dd_size)
            .seek(i)
        )
        dd.run()
        sync()
        time.sleep(0.1)
        dirty = cache.get_io_class_statistics(io_class_id=ioclass_id).usage_stats.dirty
        if dirty.get_value(Unit.Blocks4096) != (i + 1) * dd_count:
            TestRun.LOGGER.error(f"Wrong amount of dirty data ({dirty}).")


@pytest.mark.require_disk("cache", DiskTypeSet([DiskType.optane, DiskType.nand]))
@pytest.mark.require_disk("core", DiskTypeLowerThan("cache"))
def test_ioclass_pid():
    cache, core = prepare()

    ioclass_id = 1
    iterations = 20
    dd_count = 100
    dd_size = Size(4, Unit.KibiByte)

    Udev.disable()

    # Since 'dd' has to be executed right after writing pid to 'ns_last_pid',
    # 'dd' command is created and is appended to 'echo' command instead of running it
    dd_command = str(
        Dd()
        .input("/dev/zero")
        .output(core.system_path)
        .count(dd_count)
        .block_size(dd_size)
    )

    for i in range(iterations):
        cache.flush_cache()

        output = TestRun.executor.run("cat /proc/sys/kernel/ns_last_pid")
        if output.exit_code != 0:
            raise Exception(
                f"Failed to retrieve pid. stdout: {output.stdout} \n stderr :{output.stderr}"
            )

        # Few pids might be used by system during test preparation
        pid = int(output.stdout) + 50

        ioclass_config.add_ioclass(
            ioclass_id=ioclass_id,
            eviction_priority=1,
            allocation=True,
            rule=f"pid:eq:{pid}&done",
            ioclass_config_path=ioclass_config_path,
        )
        casadm.load_io_classes(cache_id=cache.cache_id, file=ioclass_config_path)

        TestRun.LOGGER.info(f"Running dd with pid {pid}")
        # pid saved in 'ns_last_pid' has to be smaller by one than target dd pid
        dd_and_pid_command = (
            f"echo {pid-1} > /proc/sys/kernel/ns_last_pid && {dd_command}"
        )
        output = TestRun.executor.run(dd_and_pid_command)
        if output.exit_code != 0:
            raise Exception(
                f"Failed to run dd with target pid. "
                f"stdout: {output.stdout} \n stderr :{output.stderr}"
            )
        sync()
        dirty = cache.get_io_class_statistics(io_class_id=ioclass_id).usage_stats.dirty
        if dirty.get_value(Unit.Blocks4096) != dd_count:
            TestRun.LOGGER.error(f"Wrong amount of dirty data ({dirty}).")
        ioclass_config.remove_ioclass(ioclass_id)
