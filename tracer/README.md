# 🚀 Tracer Project

**Tracer**는 Envoy + xDS Control Plane + Kafka 기반의 **동적 서비스 디스커버리 아키텍처**를 구성하는 Go 기반 프로젝트입니다.

이 디렉토리는 다음과 같은 핵심 컴포넌트들을 포함하고 있습니다:

- 🔹 **Envoy Proxy** 설정 (`envoy.yaml`)
- 🔹 **xDS Control Plane** (Go + gRPC)
- 🔹 **Backend Service (tracer)**: 자체 등록 기능 내장
- 🔹 **Kafka 연동 구조**
- 🔹 **Docker Compose** 기반의 로컬 실행 환경

---

## 🗂️ 디렉토리 구조

```plaintext
tracer/
├── main.go                # Backend 서비스 (tracer)
├── Dockerfile             # Backend 서비스용 Dockerfile
├── Makefile               # 빌드 및 실행 스크립트
├── envoy.yaml             # Envoy 프록시 설정
├── docker-compose.yaml    # 전체 구성 실행 스크립트
├── go.mod / go.sum        # 모듈 관리
└── xDS_control_plane/     # Control Plane (Go + gRPC xDS)
    ├── main.go
    ├── go.mod / go.sum
    └── Dockerfile