export VERSION ?= $(shell git show -q --format=%h)
export DIRTY ?= $(shell git diff --quiet || echo 'dirty')
export APPS ?= $(shell cd app && ls -d */ | grep -oh "[^\/]\+")
export GITLAB = registry.gitlab.com/dekape

IMAGE = ${GITLAB}/shadow/$(app)
PROJ = shadow
OSFLAG = linux

ifeq (${ENV},production)
CFLAGS = arm-linux-gnueabihf-gcc
ARMVERSION = 7
ARCHFLAG = arm
RETAG=latest
else
RETAG=devel
CFLAGS=
endif

ifeq (${ARCHFLAG},)
ifeq (${ENV},production)
ARCHFLAG = armv7
endif
ifeq (${ENV},development)
ARCHFLAG = amd64
endif
endif

all: compile

version:
	@sed -i 's/Version = "\([0-9\.]\+\)-.*"/Version = "\1-${VERSION}_${DIRTY}"/g' env/env.go

dependencies:
	@echo "Make sure golang-glide is installed"
	glide install --strip-vendor

checkenv:
ifeq (${ENV},)
	@echo "WARNING! You must set your ENV variable!"
endif

compile: version checkenv
ifeq (${ENV},production)
	#dep ensure -v
endif
	@$(foreach app, $(APPS), \
		echo compiling "$(app)" for "$(ENV)" in arch=$(ARCHFLAG); \
		GOOS=$(OSFLAG) CGO_ENABLED=0 GOARCH=$(ARCHFLAG) GOARM=$(ARMVERSION) CC=$(CFLAGS) go build -v -o $(PROJ)-$(app) app/$(app)/main.go;)

build: checkenv
	@$(foreach app, $(APPS), \
		echo building image for "$(app) $(IMAGE):$(VERSION)" for "$(ENV)"; \
		docker build -t $(IMAGE):$(VERSION) --build-arg ENV=$(ENV) --build-arg ARCHFLAG=$(ARCHFLAG) -f ./deploy/$(app)/Dockerfile.$(ENV) .;)

latest: 
	@$(foreach app, $(APPS), \
		echo tagging image for "$(app) $(IMAGE):$(VERSION)" for "$(ENV)" as $(RETAG); \
		docker tag $(IMAGE):$(VERSION) $(IMAGE):$(RETAG);)

push: latest
	@$(foreach app, $(APPS), \
		echo pushing image for "$(app) $(IMAGE):$(VERSION)"; \
		docker push $(IMAGE):$(RETAG) && docker push $(IMAGE):$(VERSION);)

start:
ifeq (${ENV},development)
	@$(foreach app, $(APPS), \
		echo running docker-compose for $(app) in ${ENV} environment; \
		docker-compose -f deploy/$(app)/compose.${ENV}.yml up;)
else
	@echo "You can only start during development"
endif

clean:
	rm -rf $(PROJ)-backend

.PHONY: all dependencies checkenv compile build latest push start clean
