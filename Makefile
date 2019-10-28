export APPS ?= $(shell cd app && ls -d */ | grep -oh "[^\/]\+")
export VERSION ?= $(shell git show -q --format=%h)
export GITLAB = registry.gitlab.com/sangkuriang-dev

IMAGE = ${GITLAB}/shadow/$(app)
PROJ = shadow
OSFLAG = linux

ifeq (${ENV},production)
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

dep:
	dep ensure -v --vendor-only

checkenv:
ifeq (${ENV},)
	@echo "WARNING! You must set your ENV variable!"
endif

compile: checkenv
ifeq (${ENV},production)
	#dep ensure -v
endif
	@$(foreach app, $(APPS), \
		echo compiling "$(app)" for "$(ENV)" in arch=$(ARCHFLAG); \
		GOOS=$(OSFLAG) CGO_ENABLED=0 GOARCH=$(ARCHFLAG) GOARM=$(ARMVERSION) go build -v -o $(PROJ)-$(app) app/$(app)/main.go;)

build: checkenv
	@$(foreach app, $(APPS), \
		echo building image for "$(app) $(IMAGE):$(VERSION)" for "$(ENV)"; \
		docker build -t $(IMAGE):$(VERSION) --build-arg ENV=$(ENV) --build-arg ARCHFLAG=$(ARCHFLAG) -f ./deploy/$(app)/Dockerfile.$(ENV) .;)

latest: 
	@$(foreach app, $(APPS), \
		echo tagging image for "$(app) $(IMAGE):$(VERSION)" for "$(ENV)" as latest; \
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

.PHONY: app
