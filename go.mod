module example.com/mod

replace github.com/google/go-sev-guest => /home/ec2-user/go-sev-guest

go 1.20

require (
	github.com/google/go-sev-guest v0.6.1
	github.com/jellydator/ttlcache/v3 v3.0.1
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/google/logger v1.1.1 // indirect
	github.com/google/uuid v1.0.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220608164250-635b8c9b7f68 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
