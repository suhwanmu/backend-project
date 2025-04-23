VERSION?=$(shell git describe --tags)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u +%FT%TZ)

PKG_PATH=embracer/utils

BUILD_FLAGS=-ldflags "\
  -X '$(PKG_PATH).Version=$(VERSION)' \
  -X '$(PKG_PATH).Commit=$(COMMIT)' \
  -X '$(PKG_PATH).BuildTime=$(BUILDTIME)'"

.PHONY: all kafka control_plane embracer infra traffic_ui

all: embracer control_plane traffic_ui kafka 
	@echo "✅ all build done (embracer, controle_plane, traffic_ui, kafka)"

infra: kafka
	@echo "✅ infra build done (kafka)"

kafka:
	docker build -t embracer_kafka:latest ./kafka/.

traffic_ui:
	docker build -t traffic_ui:latest ./traffic-testing/

control_plane: control_plane_build
	docker build -t embracer_control_plane:latest ./xDS_control_plane/.

control_plane_build :
	cd xDS_control_plane && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${BUILD_FLAGS} -o control_plane .

embracer:
	docker build -t embracer:latest .