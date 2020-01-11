# Versioning variables
GIT_DIGEST = $(shell git rev-parse --verify HEAD)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG = $(shell git --no-pager tag --points-at HEAD)

# Output variables
BINARY = "muffin"

ifeq ($(strip $(GIT_TAG)),)
	APP_VERSION=$(subst /,_,$(GIT_BRANCH))
else
	APP_VERSION=$(GIT_TAG)
endif

# Docker
DOCKER_IMAGE = "amitizle/muffin"
DOCKER_TAG = $(APP_VERSION)

ifeq ($(GIT_BRANCH),master)
	DOCKER_TAG=latest
endif

all: build

build: export CGO_ENABLED=0
build:
	@echo "building muffin"
	@go build -o $(BINARY) cmd/muffin/muffin.go

test:
	@go test -cover ./...

clean:
	@rm -f $(BINARY)

modules-tidy:
	@go mod tidy

modules-update:
	@go get -u ./...

docker-build: export DOCKER_BUILDKIT=1
docker-build:
	@echo "Building docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@docker build . --progress=plain -f Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG)
