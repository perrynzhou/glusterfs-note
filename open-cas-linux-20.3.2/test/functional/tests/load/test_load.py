#
# Copyright(c) 2019 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause-Clear
#


import pytest
from api.cas import casadm, casadm_parser
from core.test_run import TestRun
from storage_devices.disk import DiskType, DiskTypeSet, DiskTypeLowerThan
from test_utils.size import Size, Unit


@pytest.mark.require_disk("cache", DiskTypeSet([DiskType.optane, DiskType.nand]))
@pytest.mark.require_disk("core", DiskTypeLowerThan("cache"))
def test_load_occupied_id():
    """
        1. Start new cache instance (don't specify cache id)
        2. Add core to newly create cache.
        3. Stop cache instance.
        4. Start new cache instance on another device (don't specify cache id).
        5. Try to load metadata from first device.
            * Load should fail.
    """

    cache_device = TestRun.disks['cache']
    core_device = TestRun.disks['core']

    TestRun.LOGGER.info("Creating partitons for test")
    cache_device.create_partitions([Size(500, Unit.MebiByte), Size(500, Unit.MebiByte)])
    core_device.create_partitions([Size(1, Unit.GibiByte)])

    cache_device_1 = cache_device.partitions[0]
    cache_device_2 = cache_device.partitions[1]
    core_device = core_device.partitions[0]

    TestRun.LOGGER.info("Starting cache with default id and one core")
    cache1 = casadm.start_cache(cache_device_1, force=True)
    cache1.add_core(core_device)

    TestRun.LOGGER.info("Stopping cache")
    cache1.stop()

    TestRun.LOGGER.info("Starting cache with default id on different device")
    cache2 = casadm.start_cache(cache_device_2, force=True)

    TestRun.LOGGER.info("Attempt to load metadata from first cache device")
    try:
        casadm.load_cache(cache_device_1)
    except Exception:
        pass

    caches = casadm_parser.get_caches()
    assert len(caches) == 1, "Inappropirate number of caches after load!"
    assert caches[0].cache_device.system_path == cache_device_2.system_path
    assert caches[0].cache_id == 1

    cores = caches[0].get_core_devices()
    assert len(cores) == 0
