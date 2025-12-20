//go:build linux

package service

type (
	SERVICE_HANDLER struct {
		OnStart SERVICE_START_HANDLER
		OnStop  SERVICE_STOP_HANDLER
	}
	SERVICE_START_HANDLER func()
	SERVICE_STOP_HANDLER  func()
)

const (
	SERVICE_STOP_POLL_INTERVAL = 500
	SERVICE_STOP_TIMEOUT       = 60
)

func EnableService(serviceName string) error {
	return nil
}

func GetServiceBinaryPath(serviceName string) (string, error) {
	return "", nil
}

func GetServiceStatus(serviceName string) (uint32, error) {
	return 0, nil
}

func IsServiceAutomatic(serviceName string) (bool, error) {
	return false, nil
}

func IsServiceExists(serviceName string) bool {
	return false
}

func IsWindowsService() (bool, error) {
	return false, nil
}

func RegisterService(serviceName string, displayName string, description string) error {
	return nil
}

func Run(serviceName string, start SERVICE_START_HANDLER, stop SERVICE_STOP_HANDLER) error {
	start()
	select {}
}

func StartService(serviceName string) error {
	return nil
}

func StopService(serviceName string) error {
	return nil
}

func UnregisterService(serviceName string) bool {
	return false
}
