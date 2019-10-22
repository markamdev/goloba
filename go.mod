module github.com/markamdev/goloba

go 1.13

require github.com/markamdev/goloba/cmd/goloba v0.0.0

require github.com/markamdev/goloba/cmd/dummyserver v0.0.0

require github.com/markamdev/goloba/pkg/balancer v0.0.0

replace github.com/markamdev/goloba/cmd/goloba => ./cmd/goloba

replace github.com/markamdev/goloba/cmd/dummyserver => ./cmd/dummyserver

replace github.com/markamdev/goloba/pkg/balancer => ./pkg/balancer
