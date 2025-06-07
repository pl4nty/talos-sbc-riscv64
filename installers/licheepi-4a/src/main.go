// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
)

func main() {
	adapter.Execute(&licheePi4AInstaller{})
}

type licheePi4AInstaller struct{}

type licheePi4AExtraOptions struct{}

func (i *licheePi4AInstaller) GetOptions(extra licheePi4AExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=ttyS0,115200n8",
		"console=ttyS1,115200n8",
		"console=ttyS2,115200n8",
		"console=ttyS3,115200n8",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	return overlay.Options{
		Name:       "licheepi-4a",
		KernelArgs: kernelArgs,
	}, nil
}

func (i *licheePi4AInstaller) Install(options overlay.InstallOptions[licheePi4AExtraOptions]) error {
	var f *os.File

	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", options.InstallDisk, err)
	}

	defer f.Close() //nolint:errcheck
	
	// ROM -> U-Boot w/ SPL and DDR -> OpenSBI -> kernel
	// https://patchwork.ozlabs.org/project/uboot/patch/20250530094851.57198-6-ziyao@disroot.org/
	// vendor images use ROM -> vendor U-Boot -> OpenSBI -> kernel
	// https://wiki.sipeed.com/hardware/en/lichee/th1520/lpi4a/4_burn_image.html#Board-Boot-Process

	// can find the offsets from their "secboot" - 0 for SPL, 0x1c00000 for U-Boot
	// https://github.com/revyos/thead-u-boot/blob/93ff49d9f5bbe7942f727ab93311346173506d27/board/thead/light-c910/boot.c#L679
	// try a random offset
	// uboot, err := os.ReadFile(filepath.Join(options.ArtifactsPath, "riscv64/u-boot/licheepi-4a/u-boot-with-spl.bin"))
	// if err != nil {
	// 	return err
	// }
	// if _, err = f.WriteAt(uboot, 1024 * 8); err != nil {
	// 	return err
	// }

	// NB: In the case that the block device is a loopback device, we sync here
	// to esure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	// allows to copy a directory from the overlay to the target
	return copy.Dir(filepath.Join(options.ArtifactsPath, "riscv64/dtb"), filepath.Join(options.MountPrefix, "/boot/EFI/dtb"))
}
