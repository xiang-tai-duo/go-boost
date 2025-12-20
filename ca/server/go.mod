module server

go 1.26

require (
	common v0.0.0
	github.com/xiang-tai-duo/go-boost v0.0.0
)

require (
	github.com/mitchellh/go-ps v1.0.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
)

replace (
	common => ../common
	github.com/xiang-tai-duo/go-boost => ../..
)
