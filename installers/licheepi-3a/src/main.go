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
	adapter.Execute(context.Background(), &licheePi3AInstaller{})
}

type licheePi3AInstaller struct{}

type licheePi3AExtraOptions struct{}

func (i *licheePi3AInstaller) GetOptions(ctx context.Context, extra licheePi3AExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=tty0",
		"console=ttyS0,115200",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	return overlay.Options{
		Name:       "licheepi-3a",
		KernelArgs: kernelArgs,
	}, nil
}

func (i *licheePi3AInstaller) Install(ctx context.Context, options overlay.InstallOptions[licheePi3AExtraOptions]) error {
	// ROM -> vendor FSBL -> U-Boot (see artifacts/u-boot/k1) -> OpenSBI -> kernel
	// The K1 boot ROM loads the vendor FSBL from SPI NOR, so nothing is written
	// to the install disk here.

	// The LicheePi 3A has no upstream DTS, so riscv64/dtb carries no DTB for it
	// yet - U-Boot has to supply the fdt. The copy still runs so the board picks
	// up its DTB automatically once one lands in the kernel package.
	return copy.Dir(filepath.Join(options.ArtifactsPath, "riscv64/dtb"), filepath.Join(options.MountPrefix, "/boot/EFI/dtb"))
}
