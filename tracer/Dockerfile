# 1단계: Go 빌드 스테이지
FROM golang:1.24-alpine AS builder

# 빌드 도구 설치
RUN apk --no-cache add git

WORKDIR /app

# go mod 파일 복사 후 의존성 설치
COPY go.mod ./
RUN go mod tidy

# 소스 코드 복사 및 빌드
COPY . .
RUN go build -o tracer .

# 2단계: 실행용 이미지 (경량 Alpine 사용)
FROM alpine:latest

# 필수 인증서 설치 (HTTPS 요청 등 위해 필요)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 빌드된 실행 파일 복사
COPY --from=builder /app/tracer .

# 포트 노출 (옵션, 실제 tracer가 서버라면 필요)
EXPOSE 8080

# 환경 변수 입력 예시:
# TRACER_SERVICE=my_dynamic_service
# TRACER_ADDR=tracer-1:8080
# 실행 명령
CMD ["./tracer"]
