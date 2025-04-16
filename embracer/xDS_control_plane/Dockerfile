# 1단계: 런타임 스테이지
FROM alpine:latest

# 인증서 관련 패키지 (gRPC용)
RUN apk --no-cache add ca-certificates

# 작업 디렉토리 생성
WORKDIR /app

# 빌드한 바이너리 복사
COPY ./control_plane .

# 포트 오픈 (gRPC용)
EXPOSE 18000

# 실행 커맨드
CMD ["./control_plane"]
