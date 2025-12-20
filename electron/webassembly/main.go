package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ncruces/zenity"
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_TITLE              = "Message"
	ERROR_MESSAGE_MISSING_PID  = "missing pid argument"
	ERROR_MESSAGE_INVALID_PID  = "invalid pid: %s"
	ERROR_MESSAGE_JSON_MARSHAL = "json marshal failed: %s"
	NONE_PROCESS               = 0
	EXIT_CODE_SUCCESS          = 0
	EXIT_CODE_ERROR            = 1
	FIRST_ARGUMENT             = 1
	MIN_ARG_COUNT              = 2
	MESSAGE_BOX_WIDTH          = 700
	MESSAGE_BOX_HEIGHT         = 300
)

type (
	RESULT struct {
		Alive bool   `json:"alive"`
		Error string `json:"error"`
		Pid   int    `json:"pid"`
	}
)

func main() {
	result := RESULT{
		Alive: false,
		Error: "",
		Pid:   NONE_PROCESS,
	}
	exitCode := EXIT_CODE_ERROR
	pid := NONE_PROCESS
	err := error(nil)
	pidStr := ""
	if len(os.Args) < MIN_ARG_COUNT {
		showUsageTip()
		result.Error = ERROR_MESSAGE_MISSING_PID
	} else {
		pidStr = os.Args[FIRST_ARGUMENT]
		if pid, err = strconv.Atoi(pidStr); err == nil {
			result.Pid = pid
			alive := false
			if alive, err = isProcessAlive(pid); err == nil {
				result.Alive = alive
				exitCode = EXIT_CODE_SUCCESS
			} else {
				result.Error = err.Error()
			}
		} else {
			result.Error = fmt.Sprintf(ERROR_MESSAGE_INVALID_PID, err.Error())
		}
	}
	outputResult(result)
	os.Exit(exitCode)
}

func getTitle(title string) string {
	result := title
	if result == "" {
		result = DEFAULT_TITLE
	}
	return result
}

func info(title string, message string) error {
	result := zenity.Info(message, zenity.Title(getTitle(title)), zenity.Width(MESSAGE_BOX_WIDTH), zenity.Height(MESSAGE_BOX_HEIGHT))
	return result
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection
func showUsageTip() {
	title := "WebAssembly Process Checker"
	message := "Usage: webassembly <pid>\n\n" +
		"This tool checks if a process is alive by its PID.\n" +
		"Example: webassembly 12345\n" +
		"Output format: JSON with fields: alive(bool), error(string), pid(int)"
	info(title, message)
}

func outputResult(result RESULT) {
	data := []byte(nil)
	err := error(nil)
	if data, err = json.Marshal(result); err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Printf("{\"alive\":false,\"error\":\"%s\",\"pid\":%d}\n", fmt.Sprintf(ERROR_MESSAGE_JSON_MARSHAL, err.Error()), result.Pid)
	}
}
