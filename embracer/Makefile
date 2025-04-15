.PHONY: all kafka control_plane embracer infra

all: embracer kafka
	@echo "✅ all build done (embracer, controle_plane, kafka)"

infra: kafka
	@echo "✅ infra build done (kafka)"

kafka:
	docker build -t embracer_kafka:latest ./kafka/.

control_plane:
	docker build -t embracer_control_plane:latest ./xDS_control_plane/.

embracer: control_plane
	docker build -t embracer:latest .