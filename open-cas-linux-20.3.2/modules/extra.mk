#
# Copyright(c) 2012-2019 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause-Clear
#
ifneq ($(KERNELRELEASE),)

ifeq ($(CAS_EXT_EXP),1)
EXTRA_CFLAGS += -DWI_AVAILABLE
endif

else #KERNELRELEASE

.PHONY: sync distsync

sync:
	@cd $(OCFDIR) && $(MAKE) inc O=$(PWD)
	@cd $(OCFDIR) && $(MAKE) src O=$(PWD)/cas_cache

distsync:
	@cd $(OCFDIR) && $(MAKE) distclean O=$(PWD)
	@cd $(OCFDIR) && $(MAKE) distclean O=$(PWD)/cas_cache

endif
