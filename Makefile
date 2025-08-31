# Copyright (C) 2025 The go-ble Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)

MODULE_ROOT=github.com/cybergarage/go-ble
PKG_NAME=ble
PKG_COVER=${PKG_NAME}-cover

PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_DIR}

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG=${MODULE_ROOT}/${TEST_PKG_DIR}

.PHONY: format vet lint clean
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SRC_DIR} && ./version.gen > version.go && popd
	-git commit ${PKG_SRC_DIR}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_DIR} ${TEST_PKG_DIR} ${BIN_SRC_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID}

lint: vet
	golangci-lint run ${PKG_SRC_DIR}/... ${TEST_PKG_DIR}/...

godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	open http://localhost:6060/pkg/${PKG_ID}/ || xdg-open http://localhost:6060/pkg/${PKG_ID}/ || gnome-open http://localhost:6060/pkg/${PKG_ID}/
	godoc -http=:6060 -play

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

cover: test
	open ${PKG_COVER}.html || xdg-open ${PKG_COVER}.html || gnome-open ${PKG_COVER}.html

build:
	go build -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINS}

install:
	go install -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINS}

clean:
	go clean -i ${PKG}
