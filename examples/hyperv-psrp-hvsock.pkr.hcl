packer {
  required_plugins {
    hyperv = {
      version = ">= 0.0.1"
      source  = "github.com/Geogboe/hyperv"
    }
  }
}

variable "iso_url" {
  type    = string
  default = "https://software-static.download.prss.microsoft.com/sg/download/888969d5-f34g-4e03-ac9d-1f9786c66749/SERVER_EVAL_x64FRE_en-us.iso"
}

variable "iso_checksum" {
  type    = string
  default = "sha256:888969d5f34g4e03ac9d1f9786c66749" // Placeholder
}

source "hyperv-iso" "hvsock" {
  vm_name           = "psrp-hvsock-test"
  iso_url           = var.iso_url
  iso_checksum      = var.iso_checksum
  communicator      = "psrp"
  psrp_username     = "Administrator"
  psrp_password     = "Password123!"
  psrp_transport    = "hvsock"
  // Note: psrp_vmid is auto-detected by the plugin!
  
  cpus              = 2
  memory            = 4096
  disk_size         = 20000
  enable_secure_boot = false
  generation        = 2
  switch_name       = "Default Switch"

  // Boot command to bypass "Press any key to boot from CD or DVD..."
  boot_command      = ["<spacebar><wait><spacebar><wait><enter>"]
  boot_wait         = "1s"
  
  // Provide autounattend.xml for automated installation
  // Using cd_content as cd_files seems to have mapping issues
  cd_content        = {
    "Autounattend.xml" = file("./Autounattend.xml")
  }
  
  // Basic shutdown command
  shutdown_command  = "shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\""
}

build {
  sources = ["source.hyperv-iso.hvsock"]

  provisioner "file" {
    source      = "./test_file.txt"
    destination = "C:\\Windows\\Temp\\test_file.txt"
  }

  provisioner "powershell" {
    inline = [
      "Write-Host 'Connected via PSRP HvSocket (PowerShell Direct)!'",
      "Write-Host '--- System Info ---'",
      "$env:COMPUTERNAME",
      "[System.Environment]::OSVersion.ToString()",
      "Write-Host '--- File Verification ---'",
      "if (Test-Path C:\\Windows\\Temp\\test_file.txt) { Write-Host 'File copy successful!' } else { Write-Error 'File copy failed!' }",
      "Get-Content C:\\Windows\\Temp\\test_file.txt"
    ]
  }
}
