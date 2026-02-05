See [Releases](https://github.com/hashicorp/packer-plugin-hyperv/releases) for latest CHANGELOG information.

## Unreleased

### Features

* **PSRP Communicator Support:** Added full support for the PSRP (PowerShell Remoting Protocol) communicator.
* **HvSocket Support:** Added `psrp_transport = "hvsock"` support, allowing PSRP connections directly to the VM via Hyper-V sockets without networking.
* **Auto-detect VMID:** The plugin now automatically detects the VM's GUID for HvSocket connections.

### Improvements

* **Automated Installation:** Added `cd_content` examples and `Autounattend.xml` support for fully automated Windows installation.
* **Boot Command:** Improved boot command timing and key sequences to bypass "Press any key" prompts on UEFI Windows builds.

### Bug Fixes

* **HCL2 Specs:** Fixed generated HCL2 specs for embedded communicator configuration.

## 1.0.0 (June 14, 2021)

* Update packer-plugin-sdk to version 0.2.3. [GH-29]
* Add disable_shutdown option to Hyper-V builders. [GH-23]

## 0.0.1 (April 21, 2021)

* Hyper-V Plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md).
