// domain/port/loader.go
package port

type Loader interface {
	RunCPULoad(durationSec int) error
	RunMemoryLoad(durationSec int, mb int) error
	RunDiskIOLoad(durationSec int) error
	RunNetworkIOLoad(durationSec int, url string) error
	RunMixedLoad(durationSec int, url string, mb int) error
}
