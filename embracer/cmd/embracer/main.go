package main

import (
	"embracer/internal/adapter/controlplane"
	ahttp "embracer/internal/adapter/http"
	"embracer/internal/adapter/loader"
	"embracer/internal/adapter/validator"
	"embracer/internal/app"
	"embracer/internal/config"
	"embracer/utils/log"
	"net/http"
	"time"
)

// 등록 요청 형식
type RegisterRequest struct {
	Service string `json:"service"`
	Addr    string `json:"addr"`
}

func main() {
	cfg := config.Load()
	log.Info().Msgf("✅ Embracer starting on %s", cfg.HttpPort)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	cpClient := controlplane.NewClient(httpClient, "http://"+cfg.ControlPlaneHost)
	registerUsecase := app.NewRegisterUsecase(cpClient)

	if err := registerUsecase.Register(cfg.ClusterName, cfg.ServiceAddr); err != nil {
		log.Fatal().Msgf("failed to register: %v", err)
	}

	log.Info().Msgf("✅ Registered successfully")

	// adapter/loader, validator 인스턴스 생성
	loaderAdapter := loader.NewLoaderAdapter() // 구현체
	validatorAdapter := validator.NewDefaultLoadValidator() // 구현체
	loadUsecase := app.NewLoadTestUsecase(loaderAdapter, validatorAdapter)
	
	// gin 라우터 등록 및 실행
	r := ahttp.NewRouter(loadUsecase)
	r.Run(cfg.HttpPort)

}
