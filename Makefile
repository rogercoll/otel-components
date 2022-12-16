include ./Makefile.Common

GROUP ?= all
FOR_GROUP_TARGET=for-$(GROUP)-target

FIND_MOD_ARGS=-type f -name "go.mod"
TO_MOD_DIR=dirname {} \; | sort | grep -E '^./'
EX_COMPONENTS=-not -path "./receiver/*" -not -path "./processor/*" -not -path "./exporter/*" -not -path "./extension/*"
EX_INTERNAL=-not -path "./internal/*"

# NONROOT_MODS includes ./* dirs (excludes . dir)
NONROOT_MODS := $(shell find . $(FIND_MOD_ARGS) -exec $(TO_MOD_DIR) )

RECEIVER_MODS_0 := $(shell find ./receiver/[a-z]* $(FIND_MOD_ARGS) -exec $(TO_MOD_DIR) )
RECEIVER_MODS := $(RECEIVER_MODS_0)
ALL_MODS := $(RECEIVER_MODS)

.PHONY: for-all
for-all:
	@echo "running $${CMD} in root"
	@$${CMD}
	@set -e; for dir in $(NONROOT_MODS); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done


TOOLS_MOD_DIR := ./internal/tools
.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && $(GOCMD) install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen

.PHONY: generate
generate:
	$(MAKE) for-all CMD="$(GOCMD) generate ./..."


# Define a delegation target for each module
.PHONY: $(ALL_MODS)
$(ALL_MODS):
	@echo "Running target '$(TARGET)' in module '$@' as part of group '$(GROUP)'"
	$(MAKE) -C $@ $(TARGET)

# Trigger each module's delegation target
.PHONY: for-all-target
for-all-target: $(ALL_MODS)

.PHONY: gotest
gotest:
	$(MAKE) $(FOR_GROUP_TARGET) TARGET="test"

