//go:build darwin

// File:        winspool_darwin.go
// Description: Darwin stub for winspool package (no-op)
// --------------------------------------------------------------------------------

package winspool

import "fmt"

//goland:noinspection GoSnakeCaseUsage
type (
	PAPER_SIZE struct {
		Width  int32
		Length int32
	}
	PAPER_INFO struct {
		Id        uint16
		PaperName string
		Size      PAPER_SIZE
	}
	PRINTER_INFO struct {
		Name      string
		PortName  string
		IsDefault bool
	}
	PRINTER_INFO_2 struct {
		pServerName         *uint16
		pPrinterName        *uint16
		pShareName          *uint16
		pPortName           *uint16
		pDriverName         *uint16
		pComment            *uint16
		pLocation           *uint16
		pDevMode            *byte
		pSepFile            *uint16
		pPrintProcessor     *uint16
		pDatatype           *uint16
		pParameters         *uint16
		pSecurityDescriptor *byte
		Attributes          uint32
		Priority            uint32
		DefaultPriority     uint32
		StartTime           uint32
		UntilTime           uint32
		Status              uint32
		cJobs               uint32
		AveragePPM          uint32
	}
	DRIVER_INFO_3 struct {
		cVersion         uint32
		pName            *uint16
		pEnvironment     *uint16
		pDriverPath      *uint16
		pDataFile        *uint16
		pConfigFile      *uint16
		pHelpFile        *uint16
		pDependentFiles  *uint16
		pMonitorName     *uint16
		pDefaultDataType *uint16
	}
	PORT_INFO_2 struct {
		pPortName    *uint16
		pMonitorName *uint16
		pDescription *uint16
		fPortType    uint32
		Reserved     uint32
	}
	JOB_INFO struct {
		JobID        uint32
		Document     string
		Status       uint32
		StatusText   string
		PagesPrinted uint32
		TotalPages   uint32
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,SpellCheckingInspection
const (
	WCHAR_SIZE                   = 2
	PRINTER_ENUM_LOCAL           = 0x00000002
	PRINTER_ENUM_CONNECTIONS     = 0x00000004
	PRINTER_ATTRIBUTE_DEFAULT    = 0x00000004
	PORT_TYPE_TCPIP_MONITOR      = 0x00000001
	MONITOR_STANDARD_TCPIP       = "Standard TCP/IP Port"
	MONITOR_TCPMON_DLL           = "TCPMON.DLL"
	MONITOR_TCPIP_KEYWORD        = "tcpip"
	MONITOR_TCP_KEYWORD          = "tcp"
	JOB_STATUS_PAUSED            = 0x00000001
	JOB_STATUS_ERROR             = 0x00000002
	JOB_STATUS_DELETING          = 0x00000004
	JOB_STATUS_SPOOLING          = 0x00000008
	JOB_STATUS_PRINTING          = 0x00000010
	JOB_STATUS_OFFLINE           = 0x00000020
	JOB_STATUS_PAPEROUT          = 0x00000040
	PROTOCOL_RAW                 = 1
	PROTOCOL_LPR                 = 2
)

func CString(s string) *uint16 {
	return nil
}

func FreeCString(ptr *uint16) {
}

func GoString(lpwsz interface{}) string {
	return ""
}

func GetDefaultPrinter() string {
	return ""
}

func GetJob(printerName string, jobName *string) []JOB_INFO {
	return nil
}

func GetJobByID(printerName string, jobID uint32) []JOB_INFO {
	return nil
}

func GetJobs(printerName string) []JOB_INFO {
	return nil
}

func GetPapersInfoW(printerName string) []PAPER_INFO {
	return nil
}

func GetPorts() map[string]PORT_INFO_2 {
	return nil
}

func GetPrinterDriverFiles(printerName string) []string {
	return nil
}

func GetPrinterIP(printerName string) (string, error) {
	return "", fmt.Errorf("GetPrinterIP not implemented on macOS")
}

func GetPrinterPortProtocol(printerName string) (uint32, uint32, error) {
	return PROTOCOL_RAW, 0, nil
}

func GetPrinters() []PRINTER_INFO {
	return nil
}

func GetTcpIpPrinters() []PRINTER_INFO {
	return nil
}

func IsJobCompleted(jobStatus uint32) bool {
	return false
}

func IsPrinterExist(printerName string) bool {
	return false
}

func IsPrinterIdle(printerName string) bool {
	return false
}

func PausePrinter(printerName string) error {
	return nil
}

func ResumePrinter(printerName string) error {
	return nil
}

func SetDefaultPrinter(printerName string) error {
	return nil
}

func WaitJobByID(printerName string, jobID uint32) error {
	return nil
}

func WaitJobCompleted(printerName string, jobName *string) error {
	return nil
}
