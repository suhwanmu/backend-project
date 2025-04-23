package loader

import (
	"embracer/utils/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type LoaderAdapter struct{}

func NewLoaderAdapter() *LoaderAdapter {
	return &LoaderAdapter{}
}

func (l *LoaderAdapter) RunCPULoad(durationSec int) error {
	log.Info().Msg("🔥 CPU 부하 시작")
	n := 100000000
	res := 0
	for i := 0; i < n; i++ {
		res += i % 7
	}
	log.Info().Msgf("CPU 부하 결과: %d", res)
	return nil
}

func (l *LoaderAdapter) RunMemoryLoad(durationSec int, mb int) error {
	log.Info().Msg("🧠 메모리 부하 시작")
	size := 500 * 1024 * 1024 // 500MB
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	log.Info().Msgf("메모리 할당 완료: %d 바이트", len(b))
	return nil
}

func (l *LoaderAdapter) RunDiskIOLoad(durationSec int) error {
	log.Info().Msg("💾 디스크 IO 부하 시작")
	content := strings.Repeat("disk-io-test-", 1024*1024) // ~13MB
	tmpFile, err := ioutil.TempFile("", "loadtest-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	for i := 0; i < 20; i++ {
		if _, err := tmpFile.WriteString(content); err != nil {
			return err
		}
	}
	log.Info().Msg("디스크 쓰기 완료")
	return nil
}

func (l *LoaderAdapter) RunNetworkIOLoad(durationSec int, url string) error {
	log.Info().Msg("🌐 네트워크 IO 부하 시작")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Info().Msgf("네트워크 응답 수신 (%d bytes)", len(body))
	return nil
}

func (l *LoaderAdapter) RunMixedLoad(durationSec int, url string, mb int) error {
	log.Info().Msg("🌀 Mixed 부하 시작")
	go l.RunCPULoad(durationSec)
	go l.RunMemoryLoad(durationSec, mb)
	go l.RunDiskIOLoad(durationSec)
	go l.RunNetworkIOLoad(durationSec, url)
	time.Sleep(5 * time.Second)
	log.Info().Msg("🌀 Mixed 부하 종료")
	return nil
}