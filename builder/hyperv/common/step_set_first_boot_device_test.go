package common

import (
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
)

type parseBootDeviceIdentifierTest struct {
	generation         uint
	deviceIdentifier   string
	controllerType     string
	controllerNumber   uint
	controllerLocation uint
	shouldError        bool
}

var parseIdentifierTests = [...]parseBootDeviceIdentifierTest{
	{1, "IDE", "IDE", 0, 0, false},
	{1, "idE", "IDE", 0, 0, false},
	{1, "CD", "CD", 0, 0, false},
	{1, "cD", "CD", 0, 0, false},
	{1, "DVD", "CD", 0, 0, false},
	{1, "Dvd", "CD", 0, 0, false},
	{1, "FLOPPY", "FLOPPY", 0, 0, false},
	{1, "FloppY", "FLOPPY", 0, 0, false},
	{1, "NET", "NET", 0, 0, false},
	{1, "net", "NET", 0, 0, false},
	{1, "", "", 0, 0, true},
	{1, "bad", "", 0, 0, true},
	{1, "IDE:0:0", "", 0, 0, true},
	{1, "SCSI:0:0", "", 0, 0, true},
	{2, "IDE", "", 0, 0, true},
	{2, "idE", "", 0, 0, true},
	{2, "CD", "CD", 0, 0, false},
	{2, "cD", "CD", 0, 0, false},
	{2, "DVD", "CD", 0, 0, false},
	{2, "Dvd", "CD", 0, 0, false},
	{2, "FLOPPY", "", 0, 0, true},
	{2, "FloppY", "", 0, 0, true},
	{2, "NET", "NET", 0, 0, false},
	{2, "net", "NET", 0, 0, false},
	{2, "", "", 0, 0, true},
	{2, "bad", "", 0, 0, true},
	{2, "IDE:0:0", "IDE", 0, 0, false},
	{2, "SCSI:0:0", "SCSI", 0, 0, false},
	{2, "Ide:0:0", "IDE", 0, 0, false},
	{2, "sCsI:0:0", "SCSI", 0, 0, false},
	{2, "IDEscsi:0:0", "", 0, 0, true},
	{2, "SCSIide:0:0", "", 0, 0, true},
	{2, "IDE:0", "", 0, 0, true},
	{2, "SCSI:0", "", 0, 0, true},
	{2, "IDE:0:a", "", 0, 0, true},
	{2, "SCSI:0:a", "", 0, 0, true},
	{2, "IDE:0:653", "", 0, 0, true},
	{2, "SCSI:-10:0", "", 0, 0, true},
}

func TestStepSetFirstBootDevice_impl(t *testing.T) {
	var _ multistep.Step = new(StepSetFirstBootDevice)
}

func TestStepSetFirstBootDevice_ParseIdentifier(t *testing.T) {

	for _, identifierTest := range parseIdentifierTests {

		controllerType, controllerNumber, controllerLocation, err := ParseBootDeviceIdentifier(
			identifierTest.deviceIdentifier,
			identifierTest.generation)

		if (err != nil) != identifierTest.shouldError {

			t.Fatalf("Test %q (gen %v): shouldError: %v but err: %v", identifierTest.deviceIdentifier,
				identifierTest.generation, identifierTest.shouldError, err)

		}

		switch {

		case controllerType != identifierTest.controllerType:
			t.Fatalf("Test %q (gen %v): controllerType: %q != %q", identifierTest.deviceIdentifier, identifierTest.generation,
				identifierTest.controllerType, controllerType)

		case controllerNumber != identifierTest.controllerNumber:
			t.Fatalf("Test %q (gen %v): controllerNumber: %v != %v", identifierTest.deviceIdentifier, identifierTest.generation,
				identifierTest.controllerNumber, controllerNumber)

		case controllerLocation != identifierTest.controllerLocation:
			t.Fatalf("Test %q (gen %v): controllerLocation: %v != %v", identifierTest.deviceIdentifier, identifierTest.generation,
				identifierTest.controllerLocation, controllerLocation)

		}
	}
}

func TestStepSetFirstBootDevice(t *testing.T) {

	state := testState(t)
	step := new(StepSetFirstBootDevice)
	driver := state.Get("driver").(*DriverMock)

	for _, identifierTest := range parseIdentifierTests {

		driver.SetFirstBootDevice_Called = false
		driver.SetFirstBootDevice_VmName = ""
		driver.SetFirstBootDevice_ControllerType = ""
		driver.SetFirstBootDevice_ControllerNumber = 0
		driver.SetFirstBootDevice_ControllerLocation = 0
		driver.SetFirstBootDevice_Generation = 0

		step.Generation = identifierTest.generation
		step.FirstBootDevice = identifierTest.deviceIdentifier

		action := step.Run(context.Background(), state)
		if (action != multistep.ActionContinue) != identifierTest.shouldError {

			t.Fatalf("Test %q (gen %v): shouldError: %v but err: %v", identifierTest.deviceIdentifier,
				identifierTest.generation, identifierTest.shouldError, err)

		}

	}

}
