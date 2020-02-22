package common

import (
	"testing"

	"github.com/hashicorp/packer/helper/multistep"
)

func TestStepSetFirstBootDevice_impl(t *testing.T) {
	var _ multistep.Step = new(StepSetFirstBootDevice)
}

func TestStepSetFirstBootDevice(t *testing.T) {
//	t.Fatal("Fail IT!")
}

type parseBootDeviceIdentifierTest struct {
	generation         uint
	deviceIdentifier   string
	controllerType     string
	controllerNumber   uint
	controllerLocation uint
	shouldError        bool
}

func TestStepSetFirstBootDevice_ParseIdentifier(t *testing.T) {

	identifierTests := [...]parseBootDeviceIdentifierTest{
		{1, "IDE", "IDE", 0, 0, true},
	}

	for _, identifierTest := range identifierTests {

		controllerType, controllerNumber, controllerLocation, err := ParseBootDeviceIdentifier(
			identifierTest.deviceIdentifier,
			identifierTest.generation)

		if (err != nil) != identifierTest.shouldError {

			t.Fatalf("Test %q (gen %v) has shouldError: %v but err: %v", identifierTest.deviceIdentifier, 
				identifierTest.generation, identifierTest.shouldError, err)
			
		}

		if controllerType == "" || controllerNumber == 0 || controllerLocation == 0 {
			t.Fatal("Bah")
		}

		t.Fatal(identifierTest.deviceIdentifier)
		
	}

//	func ParseBootDeviceIdentifier(deviceIdentifier string, generation uint) (string, uint, uint, error) {

	t.Fatal("Fail Parsing!")

}