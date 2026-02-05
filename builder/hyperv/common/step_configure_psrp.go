// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/smnsjas/packer-psrp-communicator/communicator/psrp"
)

// StepConfigurePSRP configures PSRP communicator settings that can only be determined
// at runtime, such as the VM ID for HvSocket transport.
type StepConfigurePSRP struct {
	CommConfig *CommConfig
}

func (s *StepConfigurePSRP) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// We only need to do this if we are using PSRP communicator
	if s.CommConfig.Comm.Type != "psrp" {
		return multistep.ActionContinue
	}

	// We only need to detect VMID if using HvSocket transport
	// Check the internal PSRP config which is populated in Prepare()
	if s.CommConfig.PSRP.PSRPTransport != psrp.TransportHvSocket {
		return multistep.ActionContinue
	}

	// If VMID is already specified by user, don't overwrite it
	if s.CommConfig.PSRP.PSRPVMID != "" {
		return multistep.ActionContinue
	}

	driver := state.Get("driver").(Driver)
	vmName := state.Get("vmName").(string)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Auto-detecting VM ID for PSRP HvSocket connection...")
	vmid, err := driver.GetVMId(vmName)
	if err != nil {
		err := fmt.Errorf("error getting VM ID: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Detected VM ID: %s", vmid))

	// Update the configuration
	// This modification persists because s.CommConfig is a pointer to the builder's config
	s.CommConfig.PSRP.PSRPVMID = vmid

	return multistep.ActionContinue
}

func (s *StepConfigurePSRP) Cleanup(state multistep.StateBag) {}
