SED=$(shell which sed)
.DEFAULT_GOAL := help
background_option=-d
nickname=
repository_name=$(shell basename $(PWD))

DOCKER_IMAGE     := $(repository_name)
DOCKER_CONTAINER := $(repository_name)

.PHONY: setup/* docker/*

setup/mac: $(nickname)
	$(MAKE) setup/bsd

setup/bsd: $(nickname) ## for mac
	$(SED) -i '' -e 's/original/$(nickname)/g' ./$(nickname)/*.go
	$(SED) -i '' -e 's/original/$(nickname)/g' ./$(nickname)/**/*.go

setup/gnu: $(nickname) ## for linux
	$(SED) --in-place 's/original/$(nickname)/g' ./$(nickname)/*.go
	$(SED) --in-place 's/original/$(nickname)/g' ./$(nickname)/**/*.go

$(nickname):
	cp -rf original $(nickname)

docker/build:
	docker build -t $(DOCKER_IMAGE) .

docker/deps:
	docker run --rm --name $(DOCKER_CONTAINER) -v $(CURDIR):/go/src/github.com/VG-Tech-Dojo/vg-1day-2018 -it $(DOCKER_IMAGE) -C original deps

docker/run:
	docker run --rm --name $(DOCKER_CONTAINER) -p 8080:8080 -v $(CURDIR):/go/src/github.com/VG-Tech-Dojo/vg-1day-2018 -it $(DOCKER_IMAGE)

docker/deps/%: $(@F)
	docker run --rm --name $(DOCKER_CONTAINER) -v $(CURDIR):/go/src/github.com/VG-Tech-Dojo/vg-1day-2018 -it $(DOCKER_IMAGE) -C $(@F) deps

docker/run/%: $(@F)
	docker run --rm --name $(DOCKER_CONTAINER) -p 8080:8080 -v $(CURDIR):/go/src/github.com/VG-Tech-Dojo/vg-1day-2018 -it $(DOCKER_IMAGE) -C $(@F) run
