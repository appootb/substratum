## Shipping package
ifeq ($(shell uname), Darwin)
PROTO_ROOT_DIR := $(shell brew --prefix)/Cellar/protobuf/*
else
PROTO_ROOT_DIR := /usr/local
endif

SUBDIRS = $(shell ls ./appootb)

default: $(SUBDIRS)

$(SUBDIRS):
	# generating go/$@
	@rm -rf go/$@
	@mkdir -p go/$@
	@protoc -Iappootb/$@ -I. \
		--go_out=paths=source_relative:go/$@ \
		appootb/$@/*.proto
