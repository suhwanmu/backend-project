package utils

import (
	"embracer/utils/log"
	"fmt"
	"time"

	"golang.org/x/mod/semver"
)

var (
	Version   string
	Commit    string
	BuildTime string
)

func init() {
	// 기본값 설정
	if Version == "" {
		Version = "v0.0.0"
	}
	if Version[0] != 'v' {
		Version = fmt.Sprintf("v%s", Version)
	}
	if !semver.IsValid(Version) {
		log.Fatal().Msgf("Provided version %s is not valid", Version)
	}

	if Commit == "" {
		Commit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = time.Now().UTC().Format(time.RFC3339)
	}
}
