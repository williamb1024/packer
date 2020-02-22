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

func TestStepSetFirstBootDevice_ParseIdentifier(t *testing.T) {
	t.Fatal("Fail Parsing!")

}