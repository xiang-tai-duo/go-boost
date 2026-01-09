module github.com/xiang-tai-duo/go-boost

go 1.25.0

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/gorilla/websocket v1.5.3
	github.com/mitchellh/go-ps v1.0.0
	github.com/mochi-mqtt/server/v2 v2.7.9
	github.com/mutecomm/go-sqlcipher/v4 v4.4.2
	golang.org/x/crypto v0.46.0
	golang.org/x/net v0.48.0
)

require (
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
)

replace github.com/xiang-tai-duo/go-boost => ./
