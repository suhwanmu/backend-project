// internal/domain/port/validator.go
package port

type LoadValidator interface {
	ValidateCPULoad(duration int) error
	ValidateMemoryLoad(duration int, mb int) error
	ValidateDiskIOLoad(duration int) error
	ValidateNetworkIOLoad(duration int, url string) error
	ValidateMixedLoad(duration int, url string, mb int) error
}
