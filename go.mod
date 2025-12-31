module github.com/xiang-tai-duo/go-boost

go 1.25.0

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/gorilla/websocket v1.5.3
	github.com/mitchellh/go-ps v1.0.0
	github.com/mutecomm/go-sqlcipher/v4 v4.4.2
	golang.org/x/crypto v0.42.0
)

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

replace github.com/xiang-tai-duo/go-boost => ./
