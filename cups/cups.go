//go:build linux || darwin

// Package bootstrap
// File:        cups.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/cups/cups.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: CUPS wrapper for Linux printing operations.
// --------------------------------------------------------------------------------
package bootstrap

// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -lcups
// #include "cups.h"
// #include <stdlib.h>
import "C"
import (
	"errors"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	CUPS struct{}

	CUPS_OPTION struct {
		Name  string
		Value string
	}

	CUPS_DEST struct {
		Name                string
		Instance            string
		IsDefault           bool
		Info                string
		Location            string
		MakeAndModel        string
		DeviceURI           string
		PrinterURI          string
		PrinterURISupported string
		StateReasons        string
		StateMessage        string
		AuthInfoRequired    string
		MediaDefault        string
		SidesDefault        string
		ColorModeDefault    string
		FinishingsDefault   string
		PrintQualityDefault string
		OrientationDefault  string
		CopiesDefault       string
		NumberUpDefault     string
		JobSheetsDefault    string
		State               int
		PrinterType         int
		IsAcceptingJobs     bool
		IsShared            bool
		IsTemporary         bool
		Options             []CUPS_OPTION
	}

	CUPS_JOB struct {
		ID             int
		Dest           string
		Title          string
		User           string
		Format         string
		State          int
		Size           int
		Priority       int
		CompletedTime  int64
		CreationTime   int64
		ProcessingTime int64
	}

	CUPS_ATTRIBUTE struct {
		Name     string
		Value    string
		GroupTag int
		ValueTag int
		IntValue int
	}

	CUPS_IPP_RESPONSE struct {
		Status        int
		StatusMessage string
		Attributes    []CUPS_ATTRIBUTE
	}
)

var Cups CUPS

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	CUPS_WHICH_JOBS_ALL        = -1
	CUPS_WHICH_JOBS_ACTIVE     = 0
	CUPS_WHICH_JOBS_COMPLETED  = 1
	IPP_GET_PRINTER_ATTRIBUTES = 0x000B
	IPP_GET_JOB_ATTRIBUTES     = 0x0009
	IPP_GET_JOBS               = 0x000A
	MODULE_NAME_CUPS           = "cups"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_CUPS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_CUPS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_CUPS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_CUPS, logger.SKIP_STACK_FRAMES_BASE)
}

func (cups CUPS) CancelJob(printer string, jobID int) error {
	var result error
	cPrinter := cStringOrNil(printer)
	defer freeCString(cPrinter)
	if int(C.cups_cancel_job(cPrinter, C.int(jobID))) == 0 {
		result = errors.New(cups.LastErrorString())
	}
	return result
}

func (cups CUPS) CancelJob2(printer string, jobID int, purge bool) error {
	var result error
	cPrinter := cStringOrNil(printer)
	defer freeCString(cPrinter)
	cPurge := 0
	if purge {
		cPurge = 1
	}
	status := int(C.cups_cancel_job2(cPrinter, C.int(jobID), C.int(cPurge)))
	if status >= 0x0400 {
		result = errors.New(cups.LastErrorString())
	}
	return result
}

func (cups CUPS) DoRequest(operation int, resource string, uriAttrName string, uri string, requestedAttrs []string) *CUPS_IPP_RESPONSE {
	cResource := cStringOrNil(resource)
	cURIAttrName := cStringOrNil(uriAttrName)
	cURI := cStringOrNil(uri)
	defer freeCString(cResource)
	defer freeCString(cURIAttrName)
	defer freeCString(cURI)

	cAttrs, freeAttrs := cStringArray(requestedAttrs)
	defer freeAttrs()

	response := C.cups_do_request(C.int(operation), cResource, cURIAttrName, cURI, cAttrs, C.int(len(requestedAttrs)))
	return copyIPPResponse(response)
}

func (cups CUPS) Encryption() int {
	return int(C.cups_encryption())
}

func (cups CUPS) GetDefault() string {
	value := C.cups_get_default()
	defer C.cups_free_string(value)
	return C.GoString(value)
}

func (cups CUPS) GetDefault2() string {
	value := C.cups_get_default2()
	defer C.cups_free_string(value)
	return C.GoString(value)
}

func (cups CUPS) GetDests() []CUPS_DEST {
	var result []CUPS_DEST
	var cDests *C.cups_go_dest_t
	count := int(C.cups_get_dests(&cDests))
	if count > 0 && cDests != nil {
		result = copyDests(cDests, count)
		C.cups_free_dests(C.int(count), cDests)
	}
	return result
}

func (cups CUPS) GetJobAttributes(printer string, jobID int, requestedAttrs []string) *CUPS_IPP_RESPONSE {
	cPrinter := C.CString(printer)
	defer C.free(unsafe.Pointer(cPrinter))
	cAttrs, freeAttrs := cStringArray(requestedAttrs)
	defer freeAttrs()

	response := C.cups_get_job_attributes(cPrinter, C.int(jobID), cAttrs, C.int(len(requestedAttrs)))
	return copyIPPResponse(response)
}

func (cups CUPS) GetJobs(printer string, myJobs bool, whichJobs int) []CUPS_JOB {
	var result []CUPS_JOB
	cPrinter := cStringOrNil(printer)
	defer freeCString(cPrinter)

	cMyJobs := 0
	if myJobs {
		cMyJobs = 1
	}

	var cJobs *C.cups_go_job_t
	count := int(C.cups_get_jobs(cPrinter, C.int(cMyJobs), C.int(whichJobs), &cJobs))
	if count > 0 && cJobs != nil {
		result = make([]CUPS_JOB, 0, count)
		for i := 0; i < count; i++ {
			cJob := (*C.cups_go_job_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cJobs)) + uintptr(i)*unsafe.Sizeof(C.cups_go_job_t{})))
			result = append(result, CUPS_JOB{
				ID:             int(cJob.id),
				Dest:           C.GoString(cJob.dest),
				Title:          C.GoString(cJob.title),
				User:           C.GoString(cJob.user),
				Format:         C.GoString(cJob.format),
				State:          int(cJob.state),
				Size:           int(cJob.size),
				Priority:       int(cJob.priority),
				CompletedTime:  int64(cJob.completed_time),
				CreationTime:   int64(cJob.creation_time),
				ProcessingTime: int64(cJob.processing_time),
			})
		}
		C.cups_free_jobs(C.int(count), cJobs)
	}
	return result
}

func (cups CUPS) GetJobsIPP(printer string, requestedAttrs []string) *CUPS_IPP_RESPONSE {
	cPrinter := C.CString(printer)
	defer C.free(unsafe.Pointer(cPrinter))
	cAttrs, freeAttrs := cStringArray(requestedAttrs)
	defer freeAttrs()

	response := C.cups_get_jobs_ipp(cPrinter, cAttrs, C.int(len(requestedAttrs)))
	return copyIPPResponse(response)
}

func (cups CUPS) GetNamedDest(name string, instance string) *CUPS_DEST {
	var result *CUPS_DEST
	cName := cStringOrNil(name)
	cInstance := cStringOrNil(instance)
	defer freeCString(cName)
	defer freeCString(cInstance)

	cDest := C.cups_get_named_dest(cName, cInstance)
	if cDest != nil {
		dest := copyDest(cDest)
		C.cups_free_dest(cDest)
		result = &dest
	}
	return result
}

func (cups CUPS) GetPrinterAttributes(printer string, requestedAttrs []string) *CUPS_IPP_RESPONSE {
	cPrinter := C.CString(printer)
	defer C.free(unsafe.Pointer(cPrinter))
	cAttrs, freeAttrs := cStringArray(requestedAttrs)
	defer freeAttrs()

	response := C.cups_get_printer_attributes(cPrinter, cAttrs, C.int(len(requestedAttrs)))
	return copyIPPResponse(response)
}

func (cups CUPS) LastError() int {
	return int(C.cups_last_error())
}

func (cups CUPS) LastErrorString() string {
	value := C.cups_last_error_string()
	defer C.cups_free_string(value)
	return C.GoString(value)
}

func (cups CUPS) PrintFile(printer string, filename string, title string, options string) (int, error) {
	cPrinter := cStringOrNil(printer)
	cFilename := C.CString(filename)
	cTitle := cStringOrNil(title)
	cOptions := cStringOrNil(options)
	defer freeCString(cPrinter)
	defer C.free(unsafe.Pointer(cFilename))
	defer freeCString(cTitle)
	defer freeCString(cOptions)

	var err error
	jobID := int(C.cups_print_file(cPrinter, cFilename, cTitle, cOptions))
	if jobID == 0 {
		err = errors.New(cups.LastErrorString())
	}
	return jobID, err
}

func (cups CUPS) PrintFiles(printer string, files []string, title string, options string) (int, error) {
	jobID := 0
	var err error
	if len(files) == 0 {
		err = errors.New("files is empty")
	} else {
		cPrinter := cStringOrNil(printer)
		cTitle := cStringOrNil(title)
		cOptions := cStringOrNil(options)
		defer freeCString(cPrinter)
		defer freeCString(cTitle)
		defer freeCString(cOptions)

		cFiles := make([]*C.char, len(files))
		for i, file := range files {
			cFiles[i] = C.CString(file)
			defer C.free(unsafe.Pointer(cFiles[i]))
		}

		jobID = int(C.cups_print_files(cPrinter, C.int(len(cFiles)), (**C.char)(unsafe.Pointer(&cFiles[0])), cTitle, cOptions))
		if jobID == 0 {
			err = errors.New(cups.LastErrorString())
		}
	}
	return jobID, err
}

func (cups CUPS) Server() string {
	value := C.cups_server()
	defer C.cups_free_string(value)
	return C.GoString(value)
}

func (cups CUPS) SetEncryption(encryption int) {
	C.cups_set_encryption(C.int(encryption))
}

func (cups CUPS) SetServer(server string) {
	cServer := C.CString(server)
	defer C.free(unsafe.Pointer(cServer))
	C.cups_set_server(cServer)
}

func (cups CUPS) SetUser(user string) {
	cUser := C.CString(user)
	defer C.free(unsafe.Pointer(cUser))
	C.cups_set_user(cUser)
}

func (cups CUPS) User() string {
	value := C.cups_user()
	defer C.cups_free_string(value)
	return C.GoString(value)
}

func cStringArray(values []string) (**C.char, func()) {
	var result **C.char
	free := func() {}
	if len(values) > 0 {
		cValues := make([]*C.char, len(values))
		for i, value := range values {
			cValues[i] = C.CString(value)
		}
		result = (**C.char)(unsafe.Pointer(&cValues[0]))
		free = func() {
			for _, value := range cValues {
				C.free(unsafe.Pointer(value))
			}
		}
	}
	return result, free
}

func cStringOrNil(value string) *C.char {
	var result *C.char
	if value != "" {
		result = C.CString(value)
	}
	return result
}

func copyDest(cDest *C.cups_go_dest_t) CUPS_DEST {
	dest := CUPS_DEST{
		Name:                C.GoString(cDest.name),
		Instance:            C.GoString(cDest.instance),
		IsDefault:           int(cDest.is_default) != 0,
		Info:                C.GoString(cDest.info),
		Location:            C.GoString(cDest.location),
		MakeAndModel:        C.GoString(cDest.make_and_model),
		DeviceURI:           C.GoString(cDest.device_uri),
		PrinterURI:          C.GoString(cDest.printer_uri),
		PrinterURISupported: C.GoString(cDest.printer_uri_supported),
		StateReasons:        C.GoString(cDest.state_reasons),
		StateMessage:        C.GoString(cDest.state_message),
		AuthInfoRequired:    C.GoString(cDest.auth_info_required),
		MediaDefault:        C.GoString(cDest.media_default),
		SidesDefault:        C.GoString(cDest.sides_default),
		ColorModeDefault:    C.GoString(cDest.color_mode_default),
		FinishingsDefault:   C.GoString(cDest.finishings_default),
		PrintQualityDefault: C.GoString(cDest.print_quality_default),
		OrientationDefault:  C.GoString(cDest.orientation_default),
		CopiesDefault:       C.GoString(cDest.copies_default),
		NumberUpDefault:     C.GoString(cDest.number_up_default),
		JobSheetsDefault:    C.GoString(cDest.job_sheets_default),
		State:               int(cDest.state),
		PrinterType:         int(cDest.printer_type),
		IsAcceptingJobs:     int(cDest.is_accepting_jobs) != 0,
		IsShared:            int(cDest.is_shared) != 0,
		IsTemporary:         int(cDest.is_temporary) != 0,
	}
	for i := 0; i < int(cDest.num_options); i++ {
		cOpt := (*C.cups_go_option_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cDest.options)) + uintptr(i)*unsafe.Sizeof(C.cups_go_option_t{})))
		dest.Options = append(dest.Options, CUPS_OPTION{
			Name:  C.GoString(cOpt.name),
			Value: C.GoString(cOpt.value),
		})
	}
	return dest
}

func copyDests(cDests *C.cups_go_dest_t, count int) []CUPS_DEST {
	dests := make([]CUPS_DEST, 0, count)
	for i := 0; i < count; i++ {
		cDest := (*C.cups_go_dest_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cDests)) + uintptr(i)*unsafe.Sizeof(C.cups_go_dest_t{})))
		dests = append(dests, copyDest(cDest))
	}
	return dests
}

func copyIPPResponse(cResponse *C.cups_go_ipp_response_t) *CUPS_IPP_RESPONSE {
	var response *CUPS_IPP_RESPONSE
	if cResponse != nil {
		defer C.cups_free_ipp_response(cResponse)

		response = &CUPS_IPP_RESPONSE{
			Status:        int(cResponse.status),
			StatusMessage: C.GoString(cResponse.status_message),
		}
		for i := 0; i < int(cResponse.num_attrs); i++ {
			cAttr := (*C.cups_go_attr_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cResponse.attrs)) + uintptr(i)*unsafe.Sizeof(C.cups_go_attr_t{})))
			response.Attributes = append(response.Attributes, CUPS_ATTRIBUTE{
				Name:     C.GoString(cAttr.name),
				Value:    C.GoString(cAttr.value),
				GroupTag: int(cAttr.group_tag),
				ValueTag: int(cAttr.value_tag),
				IntValue: int(cAttr.int_value),
			})
		}
	}
	return response
}

func freeCString(value *C.char) {
	if value != nil {
		C.free(unsafe.Pointer(value))
	}
}
