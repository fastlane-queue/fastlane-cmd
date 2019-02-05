# Licensed under the MIT license:
# http://www.opensource.org/licenses/MIT-license
# Copyright (c) 2018, Backstage <backstage@corp.globo.com>

COMPOSE := $(shell command -v docker-compose 2> /dev/null)

setup:
	@go get -u github.com/wadey/gocovmerge
	@go get -u github.com/onsi/ginkgo
	@dep ensure

build:
	@go build $(PACKAGES)
	@go build -o ./bin/fastlane-cmd main.go

test: unit test-coverage-func

clear-coverage-profiles:
	@find . -name '*.coverprofile' -delete

unit: clear-coverage-profiles unit-run gather-unit-profiles

unit-run:
	#@LOGXI="*=ERR,dat:sqlx=OFF,dat=OFF" ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}
	@ginkgo -cover -r -randomizeAllSpecs -randomizeSuites -skipMeasurements ${TEST_PACKAGES}

gather-unit-profiles:
	@mkdir -p _build
	@echo "mode: count" > _build/coverage-unit.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/coverage-unit.out; done'

merge-profiles:
	@mkdir -p _build
	@gocovmerge _build/*.out > _build/coverage-all.out

test-coverage-func coverage-func: merge-profiles
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions NOT COVERED by Tests"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@go tool cover -func=_build/coverage-all.out | egrep -v "100.0[%]"

test-coverage-html coverage-html: merge-profiles
	@go tool cover -html=_build/coverage-all.out

enqueue:
	@go run main.go enqueue -v3

.PHONY: static

static:
	@cd static && go-bindata -o ../core/static.go -pkg core ./*
	@cd third-party/backstage-stack && go-bindata -o ../../stack/stack.go -pkg stack ./*

clean: clear-coverage-profiles

build-linux-64:
	@mkdir -p ./bin
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/fastlane-cmd-linux-x86_64
	@chmod +x bin/*

cross:
	@mkdir -p ./bin
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ./bin/fastlane-cmd-linux-i386
	$(MAKE) build-linux-64
	@echo "Building for darwin-i386..."
	@env GOOS=darwin GOARCH=386 go build -o ./bin/fastlane-cmd-darwin-i386
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ./bin/fastlane-cmd-darwin-x86_64
	@chmod +x bin/*

deps-build:
ifdef COMPOSE
	@echo "Starting dependencies (rebuilding)..."
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml rm -f
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml pull
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml up --build -d --remove-orphans
	@echo "Dependencies started successfully."
endif

deps:
ifdef COMPOSE
	@echo "Starting dependencies..."
	@git submodule update --remote
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml rm -f
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml up -d --remove-orphans
	@echo "Dependencies started successfully."
endif

stop-deps:
ifdef COMPOSE
	@echo "Stopping dependencies..."
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml stop
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml rm -f
endif

ps:
ifdef COMPOSE
	@docker-compose --project-name fastlane-cmd -f ./docker-compose.yml ps
endif
