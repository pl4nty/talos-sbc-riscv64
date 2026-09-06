// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
)

func main() {
	adapter.Execute(context.Background(), &k3PicoITXInstaller{})
}

type k3PicoITXInstaller struct{}

type k3PicoITXExtraOptions struct{}

func (i *k3PicoITXInstaller) GetOptions(ctx context.Context, extra k3PicoITXExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=tty0",
		"console=ttyS0,115200",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	return overlay.Options{
		Name:       "k3-pico-itx",
		KernelArgs: kernelArgs,
	}, nil
}

func (i *k3PicoITXInstaller) Install(ctx context.Context, options overlay.InstallOptions[k3PicoITXExtraOptions]) error {
	// ROM -> vendor FSBL -> vendor U-Boot -> OpenSBI -> kernel
	// U-Boot has no SpacemiT K3 support yet, so the vendor firmware in SPI NOR
	// stays in place and nothing is written to the install disk here.

	return copy.Dir(filepath.Join(options.ArtifactsPath, "riscv64/dtb"), filepath.Join(options.MountPrefix, "/boot/EFI/dtb"))
}
