package app

import (
	"embracer/internal/domain/port"
	"embracer/utils/log"
)

type LoadTestUsecase struct {
	loader port.Loader
	validator port.LoadValidator
}

func NewLoadTestUsecase(loader port.Loader, validator port.LoadValidator) *LoadTestUsecase {
	return &LoadTestUsecase{loader: loader, validator: validator}
}

func (u *LoadTestUsecase) ExecuteCPULoad(durationSec int) error {
	log.Info().Msg("Executing CPU Load")
	if err := u.validator.ValidateCPULoad(durationSec); err != nil {
		return err
	}

	return u.loader.RunCPULoad(durationSec)
}

func (u *LoadTestUsecase) ExecuteMemoryLoad(durationSec int, mb int) error {
	log.Info().Msg("Executing Memory Load")
	if err := u.validator.ValidateMemoryLoad(durationSec, mb); err != nil {
		return err
	}
	return u.loader.RunMemoryLoad(durationSec , mb )
}

func (u *LoadTestUsecase) ExecuteDiskIOLoad(durationSec int) error {
	log.Info().Msg("Executing Disk IO Load")
	if err := u.validator.ValidateDiskIOLoad(durationSec); err != nil {
		return err
	}
	return u.loader.RunDiskIOLoad(durationSec)
}

func (u *LoadTestUsecase) ExecuteNetworkIOLoad(durationSec int, url string) error {
	log.Info().Msg("Executing Network IO Load")
	if err := u.validator.ValidateNetworkIOLoad(durationSec, url); err != nil {
		return err
	}
	return u.loader.RunNetworkIOLoad(durationSec , url )
}

func (u *LoadTestUsecase) ExecuteMixedLoad(durationSec int, url string, mb int) error {
	log.Info().Msg("Executing Mixed Load")
	if err := u.validator.ValidateMixedLoad(durationSec, url, mb); err != nil {
		return err
	}
	return u.loader.RunMixedLoad(durationSec , url ,mb)
}


