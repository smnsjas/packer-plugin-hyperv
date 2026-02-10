# Hyper-V Plugin

The Hyper-V Packer Plugin creates
[Hyper-V](https://www.microsoft.com/en-us/server-cloud/solutions/virtualization.aspx)
virtual machines and exports them. It supports SSH, WinRM, and PSRP
communicators, with PSRP over HvSocket (PowerShell Direct) recommended for
Windows VMs as it requires no network configuration on the guest.

## Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    hyperv = {
      source  = "github.com/Geogboe/hyperv"
      version = "~> 1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/Geogboe/hyperv
```

### Components

#### Builders

- [hyperv-iso](builders/iso.mdx) - Creates a new Hyper-V VM from an ISO
  image, installs an OS, provisions software, then exports the machine.
  Best for building images from scratch.

- [hyperv-vmcx](builders/vmcx.mdx) - Clones an existing Hyper-V virtual
  machine or imports an exported VM, provisions software, then exports the
  machine. Best for customizing existing base images.

#### Communicators

- [PSRP](communicators/psrp.mdx) - PowerShell Remoting Protocol communicator
  with WSMan and HvSocket transport support. Recommended for Windows VMs.

### Running from WSL2

This plugin supports being run from WSL2 provided it is run from a Windows
filesystem and `PACKER_CACHE_DIR` is set to a path on the Windows filesystem.

For example, assuming a Windows username of `user`:

    /mnt/c/Users/user/$ PACKER_CACHE_DIR=/mnt/c/Users/user/.packer packer build ...
