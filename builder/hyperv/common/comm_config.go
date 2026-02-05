// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type CommConfig

package common

import (
	"time"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	psrp "github.com/smnsjas/packer-psrp-communicator/communicator/psrp"
)

// CommConfig contains configuration for all communicator types (SSH, WinRM, PSRP)
type CommConfig struct {
	// Standard communicator configuration (SSH/WinRM)
	Comm communicator.Config `mapstructure:",squash"`

	// PSRP communicator configuration (internal use)
	PSRP psrp.Config `mapstructure:"-"`

	// --- PSRP Connection Settings ---
	// The hostname or IP address of the remote host for PSRP connection.
	// Only used with wsman transport; for hvsock, this is ignored.
	PSRPHost string `mapstructure:"psrp_host" required:"false"`
	// The port to connect to. Defaults to 5985 (HTTP) or 5986 (HTTPS).
	PSRPPort int `mapstructure:"psrp_port" required:"false"`
	// The username for PSRP authentication.
	PSRPUsername string `mapstructure:"psrp_username" required:"false"`
	// The password for PSRP authentication.
	PSRPPassword string `mapstructure:"psrp_password" required:"false"`
	// Connection timeout. Defaults to "5m".
	PSRPTimeout string `mapstructure:"psrp_timeout" required:"false"`

	// --- PSRP Transport Settings ---
	// Transport type: "wsman" (HTTP/HTTPS) or "hvsock" (PowerShell Direct).
	// Defaults to "wsman".
	PSRPTransport string `mapstructure:"psrp_transport" required:"false"`
	// The Hyper-V Virtual Machine GUID for hvsock transport. If not specified,
	// it will be auto-detected from the VM being built.
	PSRPVMID string `mapstructure:"psrp_vmid" required:"false"`
	// PowerShell configuration name for hvsock. Defaults to "Microsoft.PowerShell".
	PSRPConfigurationName string `mapstructure:"psrp_configuration_name" required:"false"`

	// --- PSRP TLS Settings ---
	// Use TLS (HTTPS) for the connection. Defaults to false.
	PSRPUseTLS bool `mapstructure:"psrp_use_tls" required:"false"`
	// Skip TLS certificate verification. Defaults to false.
	PSRPInsecure bool `mapstructure:"psrp_insecure" required:"false"`

	// --- PSRP Authentication Settings ---
	// Authentication type: "basic", "ntlm", "kerberos", or "negotiate".
	// Defaults to "negotiate".
	PSRPAuthType string `mapstructure:"psrp_auth_type" required:"false"`
	// Domain for NTLM/Negotiate authentication.
	PSRPDomain string `mapstructure:"psrp_domain" required:"false"`
	// Kerberos realm. Optional on Windows (uses SSPI).
	PSRPRealm string `mapstructure:"psrp_realm" required:"false"`
}

// Prepare validates the communicator configuration.
// For "psrp" type, we skip SDK's validation since it only accepts known types.
func (c *CommConfig) Prepare(ctx *interpolate.Context) []error {
	// If using PSRP communicator, populate and validate PSRP config
	if c.Comm.Type == "psrp" {
		c.populatePSRPConfig()
		errs := c.PSRP.Prepare(ctx)

		// Filter out "psrp_vmid is required for hvsock transport" error for hvsock transport
		// because the VMID is unknown at validate time and will be auto-detected at runtime.
		if c.PSRP.PSRPTransport == psrp.TransportHvSocket {
			var filteredErrs []error
			for _, err := range errs {
				if err.Error() != "psrp_vmid is required for hvsock transport" {
					filteredErrs = append(filteredErrs, err)
				}
			}
			return filteredErrs
		}

		return errs
	}

	// For SSH/WinRM, use SDK's validation
	return c.Comm.Prepare(ctx)
}

// populatePSRPConfig copies field values from the HCL-exposed fields to the internal psrp.Config
func (c *CommConfig) populatePSRPConfig() {
	c.PSRP.PSRPHost = c.PSRPHost
	c.PSRP.PSRPPort = c.PSRPPort
	c.PSRP.PSRPUsername = c.PSRPUsername
	c.PSRP.PSRPPassword = c.PSRPPassword

	if c.PSRPTimeout != "" {
		if d, err := time.ParseDuration(c.PSRPTimeout); err == nil {
			c.PSRP.PSRPTimeout = d
		}
	}

	switch c.PSRPTransport {
	case "wsman":
		c.PSRP.PSRPTransport = psrp.TransportWSMan
	case "hvsock":
		c.PSRP.PSRPTransport = psrp.TransportHvSocket
	}

	c.PSRP.PSRPVMID = c.PSRPVMID
	c.PSRP.PSRPConfigurationName = c.PSRPConfigurationName
	c.PSRP.PSRPUseTLS = c.PSRPUseTLS
	c.PSRP.PSRPInsecureSkipVerify = c.PSRPInsecure

	switch c.PSRPAuthType {
	case "basic":
		c.PSRP.PSRPAuthType = psrp.AuthBasic
	case "ntlm":
		c.PSRP.PSRPAuthType = psrp.AuthNTLM
	case "kerberos":
		c.PSRP.PSRPAuthType = psrp.AuthKerberos
	case "negotiate":
		c.PSRP.PSRPAuthType = psrp.AuthNegotiate
	}

	c.PSRP.PSRPDomain = c.PSRPDomain
	c.PSRP.PSRPRealm = c.PSRPRealm
}
