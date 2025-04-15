package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	httpPort = ":18080"
)

// 등록 요청 형식
type RegisterRequest struct {
	Service string `json:"service"`
	Addr    string `json:"addr"`
}

func main() {
	// 환경 변수 또는 기본값 설정
	service := os.Getenv("embracer_SERVICE")
	if service == "" {
		service = "test_dynamic_service"
	}
	addr := os.Getenv("embracer_ADDR")
	if addr == "" {
		addr = "test-embracer:18080"
	}

	if service == "" || addr == "" {
		log.Fatal("❌ embracer_SERVICE and embracer_ADDR env vars must be set")
	}

	// 요청 본문 생성
	reqBody := RegisterRequest{
		Service: service,
		Addr:    addr,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("❌ Failed to marshal JSON: %v", err)
	}

	// control-plane 주소
	controlPlaneHost := os.Getenv("CONTROL_PLANE_HOST")
	if controlPlaneHost == "" {
		controlPlaneHost = "control-plane"
	}
	controlPlaneURL := fmt.Sprintf("http://%s:2222/register", controlPlaneHost)

	// 등록 요청 전송
	resp, err := http.Post(controlPlaneURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("❌ Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("❌ Registration failed with status: %s", resp.Status)
	}

	fmt.Printf("✅ Registered to control-plane as [%s] -> %s\n", service, addr)

	// 간단한 HTTP 서버 시작 (ping/pong)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	log.Println("📡 embracer is listening on ", httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatalf("❌ Failed to start embracer HTTP server: %v", err)
	}
}
