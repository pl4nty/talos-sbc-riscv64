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
	// ROM -> FSBL (U-Boot SPL) -> OpenSBI -> U-Boot -> GRUB (EFI) -> kernel
	// The boot ROM reads bootinfo and the FSBL from the eMMC boot0 hardware
	// partition, and the SPL reads U-Boot from 1 MiB into boot0, so the
	// install disk (the eMMC user area) carries nothing outside the Talos
	// partitions. artifacts/riscv64/fsbl/k1/boot0.img is written to boot0
	// separately, through the boot ROM's USB download mode.

	// There's no upstream LicheePi 3A DTS. U-Boot sets fdtfile to the BPI-F3
	// DTB, which matches the 3A (both are the K1 reference design), and the
	// EFI loader reads it from dtb/ on the ESP.
	return copy.Dir(filepath.Join(options.ArtifactsPath, "riscv64/dtb"), filepath.Join(options.MountPrefix, "/boot/EFI/dtb"))
}
