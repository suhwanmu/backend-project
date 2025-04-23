// internal/adapter/validator/load_validator.go
package validator

import (
	"errors"
	"net/url"
)

type DefaultLoadValidator struct{}

func NewDefaultLoadValidator() *DefaultLoadValidator {
	return &DefaultLoadValidator{}
}

func (v *DefaultLoadValidator) ValidateCPULoad(duration int) error {
	if duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	return nil
}

func (v *DefaultLoadValidator) ValidateMemoryLoad(duration int, mb int) error {
	if duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	if mb <= 0 {
		return errors.New("memory size must be greater than 0")
	}
	return nil
}

func (v *DefaultLoadValidator) ValidateDiskIOLoad(duration int) error {
	if duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	return nil
}

func (v *DefaultLoadValidator) ValidateNetworkIOLoad(duration int, urlStr string) error {
	if duration <= 0 {
		return errors.New("duration must be greater than 0")
	}
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return errors.New("invalid URL format")
	}
	return nil
}

func (v *DefaultLoadValidator) ValidateMixedLoad(duration int, urlStr string, mb int) error {
	if err := v.ValidateCPULoad(duration); err != nil {
		return err
	}
	if err := v.ValidateMemoryLoad(duration, mb); err != nil {
		return err
	}
	if err := v.ValidateDiskIOLoad(duration); err != nil {
		return err
	}
	if err := v.ValidateNetworkIOLoad(duration, urlStr); err != nil {
		return err
	}
	return nil
}
