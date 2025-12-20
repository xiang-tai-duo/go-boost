//go:build linux || darwin

package main

import (
	"fmt"

	bootstrap "github.com/xiang-tai-duo/go-boost/cups"
)

func main() {
	dests := bootstrap.Cups.GetDests()
	fmt.Printf("found %d printer(s)\n", len(dests))
	for i, dest := range dests {
		fmt.Printf("\n[%d] ------------------------------\n", i+1)
		fmt.Printf("    name                    : %s\n", dest.Name)
		fmt.Printf("    instance                : %s\n", dest.Instance)
		fmt.Printf("    is_default              : %t\n", dest.IsDefault)
		fmt.Printf("    is_accepting_jobs       : %t\n", dest.IsAcceptingJobs)
		fmt.Printf("    is_shared               : %t\n", dest.IsShared)
		fmt.Printf("    is_temporary            : %t\n", dest.IsTemporary)
		fmt.Printf("    state                   : %d\n", dest.State)
		fmt.Printf("    printer_type            : %d\n", dest.PrinterType)
		fmt.Printf("    info                    : %s\n", dest.Info)
		fmt.Printf("    location                : %s\n", dest.Location)
		fmt.Printf("    make_and_model          : %s\n", dest.MakeAndModel)
		fmt.Printf("    device_uri              : %s\n", dest.DeviceURI)
		fmt.Printf("    printer_uri             : %s\n", dest.PrinterURI)
		fmt.Printf("    printer_uri_supported   : %s\n", dest.PrinterURISupported)
		fmt.Printf("    state_reasons           : %s\n", dest.StateReasons)
		fmt.Printf("    state_message           : %s\n", dest.StateMessage)
		fmt.Printf("    auth_info_required      : %s\n", dest.AuthInfoRequired)
		fmt.Printf("    media_default           : %s\n", dest.MediaDefault)
		fmt.Printf("    sides_default           : %s\n", dest.SidesDefault)
		fmt.Printf("    color_mode_default      : %s\n", dest.ColorModeDefault)
		fmt.Printf("    finishings_default      : %s\n", dest.FinishingsDefault)
		fmt.Printf("    print_quality_default   : %s\n", dest.PrintQualityDefault)
		fmt.Printf("    orientation_default     : %s\n", dest.OrientationDefault)
		fmt.Printf("    copies_default          : %s\n", dest.CopiesDefault)
		fmt.Printf("    number_up_default       : %s\n", dest.NumberUpDefault)
		fmt.Printf("    job_sheets_default      : %s\n", dest.JobSheetsDefault)
		fmt.Printf("    options                 : %d\n", len(dest.Options))
		for _, opt := range dest.Options {
			fmt.Printf("      - %s = %s\n", opt.Name, opt.Value)
		}
	}

	fmt.Printf("\ndefault printer: %s\n", bootstrap.Cups.GetDefault())
}
