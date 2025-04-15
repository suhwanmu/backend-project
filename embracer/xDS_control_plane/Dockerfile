# 1단계: 빌드 스테이지
FROM golang:1.24-alpine AS builder

# 필수 도구 설치
RUN apk --no-cache add git

# 작업 디렉토리 설정
WORKDIR /app

# go mod 파일 복사 및 의존성 설치
COPY go.mod go.sum ./
RUN go mod tidy

# 소스 코드 복사 및 빌드
COPY . .
RUN go build -o controle_plane .

# 2단계: 런타임 스테이지
FROM alpine:latest

# 인증서 관련 패키지 (gRPC용)
RUN apk --no-cache add ca-certificates

# 작업 디렉토리 생성
WORKDIR /app

# 빌드한 바이너리 복사
COPY --from=builder /app/controle_plane .

# 포트 오픈 (gRPC용)
EXPOSE 18000

# 실행 커맨드
CMD ["./controle_plane"]
