.PHONY: all kafka control_plane embracer infra

all: embracer control_plane kafka
	@echo "✅ all build done (embracer, controle_plane, kafka)"

infra: kafka
	@echo "✅ infra build done (kafka)"

kafka:
	docker build -t embracer_kafka:latest ./kafka/.

control_plane: control_plane_build
	docker build -t embracer_control_plane:latest ./xDS_control_plane/.

control_plane_build :
	cd xDS_control_plane && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o control_plane .


embracer:
	docker build -t embracer:latest .