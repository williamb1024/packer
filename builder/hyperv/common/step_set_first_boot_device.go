package common

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type StepSetFirstBootDevice struct {
	Generation      uint
	FirstBootDevice string
}

func ParseBootDeviceIdentifier(deviceIdentifier string, generation uint) (string, uint, uint, error) {

	// all input strings are forced to upperCase for comparison, I believe this is
	// safe as all of our values are 7bit ASCII clean.

	lookupDeviceIdentifier := strings.ToUpper(deviceIdentifier)

	if (generation == 1) {

		// Gen1 values are a simple set of if/then/else values, which we coalesce into a map
		// here for simplicity

		lookupTable := map[string]string {
			"FLOPPY": "FLOPPY",
			"IDE": "IDE",
			"NET": "NET",
			"CD": "CD",
			"DVD": "CD",
		}

		controllerType, isDefined := lookupTable[lookupDeviceIdentifier]
		if (!isDefined) {

			return "", 0, 0, fmt.Errorf("The value %q is not a properly formatted device group identifier.", deviceIdentifier)

		}

		// success
		return controllerType, 0, 0, nil
	}

	// everything else is treated as generation 2... the first set of lookups covers
	// the simple options..

	lookupTable := map[string]string {
		"CD": "CD",
		"DVD": "CD",
		"NET": "NET",
	}

	controllerType, isDefined := lookupTable[lookupDeviceIdentifier]
	if (isDefined) {

		// these types do not require controllerNumber or controllerLocation
		return controllerType, 0, 0, nil

	}

	// not a simple option, check for a controllerType:controllerNumber:controllerLocation formatted
	// device..

	r, err := regexp.Compile("^(IDE|SCSI):(\\d+):(\\d+)$")
	if err != nil {
		return "", 0, 0, err
	}

	controllerMatch := r.FindStringSubmatch(lookupDeviceIdentifier)
	if controllerMatch != nil {

		var controllerLocation int64
		var controllerNumber int64

		// NOTE: controllerNumber and controllerLocation cannot be negative, the regex expression
		// would not have matched if either number was signed

		controllerNumber, err = strconv.ParseInt(controllerMatch[2], 10, 8)
		if err == nil {

			controllerLocation, err = strconv.ParseInt(controllerMatch[3], 10, 8)
			if err == nil {

				return controllerMatch[1], uint(controllerNumber), uint(controllerLocation), nil

			}

		}

		return "", 0, 0, err

	}

	return "", 0, 0, fmt.Errorf("The value %q is not a properly formatted device identifier.", deviceIdentifier)

	// captureExpression := "^(FLOPPY|IDE|NET)|(CD|DVD)$"
	// if generation > 1 {
	// 	captureExpression = "^((IDE|SCSI):(\\d+):(\\d+))|(DVD|CD)|(NET)$"
	// }

	// r, err := regexp.Compile(captureExpression)
	// if err != nil {
	// 	return "", 0, 0, err
	// }

	// // match against the appropriate set of values.. we force to uppercase to ensure that
	// // all devices are always in the same case

	// identifierMatches := r.FindStringSubmatch(strings.ToUpper(deviceIdentifier))
	// if identifierMatches == nil {
	// 	return "", 0, 0, fmt.Errorf("The value %q is not a properly formatted device or device group identifier.", deviceIdentifier)
	// }

	// switch {

	// // CD or DVD are always returned as "CD"
	// case ((generation == 1) && (identifierMatches[2] != "")) || ((generation > 1) && (identifierMatches[5] != "")):
	// 	return "CD", 0, 0, nil

	// // generation 1 only has FLOPPY, IDE or NET remaining..
	// case (generation == 1):
	// 	return identifierMatches[0], 0, 0, nil

	// // generation 2, check for IDE or SCSI and parse location and number
	// case (identifierMatches[2] != ""):
	// 	{

	// 		var controllerLocation int64
	// 		var controllerNumber int64

	// 		// NOTE: controllerNumber and controllerLocation cannot be negative, the regex expression
	// 		// would not have matched if either number was signed

	// 		controllerNumber, err = strconv.ParseInt(identifierMatches[3], 10, 8)
	// 		if err == nil {

	// 			controllerLocation, err = strconv.ParseInt(identifierMatches[4], 10, 8)
	// 			if err == nil {

	// 				return identifierMatches[2], uint(controllerNumber), uint(controllerLocation), nil

	// 			}

	// 		}

	// 		return "", 0, 0, err

	// 	}

	// // only "NET" left on generation 2
	// default:
	// 	return "NET", 0, 0, nil

	// }

}

func (s *StepSetFirstBootDevice) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {

	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	if s.FirstBootDevice != "" {

		controllerType, controllerNumber, controllerLocation, err := ParseBootDeviceIdentifier(s.FirstBootDevice, s.Generation)
		if err == nil {

			switch {

			case controllerType == "CD":
				{
					// the "DVD" controller is special, we only apply the setting if we actually mounted
					// an ISO and only if that was mounted as the "IsoUrl" not a secondary ISO.

					dvdControllerState := state.Get("os.dvd.properties")
					if dvdControllerState == nil {

						ui.Say("First Boot Device is DVD, but no primary ISO mounted. Ignoring.")
						return multistep.ActionContinue

					}

					ui.Say(fmt.Sprintf("Setting boot device to %q", s.FirstBootDevice))
					dvdController := dvdControllerState.(DvdControllerProperties)
					err = driver.SetFirstBootDevice(vmName, controllerType, dvdController.ControllerNumber, dvdController.ControllerLocation, s.Generation)

				}

			default:
				{
					// anything else, we just pass as is..
					ui.Say(fmt.Sprintf("Setting boot device to %q", s.FirstBootDevice))
					err = driver.SetFirstBootDevice(vmName, controllerType, controllerNumber, controllerLocation, s.Generation)
				}
			}

		}

		if err != nil {
			err := fmt.Errorf("Error setting first boot device: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt

		}

	}

	return multistep.ActionContinue
}

func (s *StepSetFirstBootDevice) Cleanup(state multistep.StateBag) {
	// do nothing
}
