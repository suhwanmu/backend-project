module utils

go 1.24.2

replace embracer/utils => ../utils

require (
	embracer/utils v0.0.0-00010101000000-000000000000
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/rs/zerolog v1.34.0
	github.com/twmb/franz-go v1.18.1
	github.com/twmb/franz-go/pkg/kadm v1.16.0
	golang.org/x/mod v0.24.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.9.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
