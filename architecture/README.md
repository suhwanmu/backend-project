# 🏗️ Architecture

서비스의 주요 컴포넌트 구성과 역할을 설명합니다.

---

## 1. 🧩 API Gateway

**Envoy + xDS Control Plane + Backend 서비스** 아키텍처로 구성되어 있으며,  
**Kafka 기반 메시지 트리거 + 동적 endpoint 등록**을 통해 확장성과 복원력을 확보합니다.

> 🔁 Backend ↔ Kafka ↔ Control Plane ↔ Envoy

Envoy는 **서비스 진입점** 역할을 수행하며, 동적으로 Backend endpoint를 구성합니다.

### ▶️ 주요 기능
- **EDS(Endpoint Discovery Service)** 기반 동적 엔드포인트 수신
- gRPC 기반 xDS API 통신
- HTTP, gRPC 요청 라우팅
- 통계, 필터링, 관찰성 구성 가능

### ⚙️ 기술 스택

| 구성 요소       | 기술 스택                                   |
|----------------|---------------------------------------------|
| API Gateway    | [Envoy Proxy](https://www.envoyproxy.io/)  |
| Control Plane  | Go, [go-control-plane](https://github.com/envoyproxy/go-control-plane) |
| Backend        | Go                                           |
| 메시지 브로커   | [Apache Kafka](https://kafka.apache.org/)   |
| 배포 환경      | Docker Compose                               |
