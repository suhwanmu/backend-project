all: tracer

control_plane:
	docker build -t control_plane:latest ./xDS_control_plane/.

tracer: control_plane
	docker build -t tracer:latest .