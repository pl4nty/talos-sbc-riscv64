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

// https://wiki.sipeed.com/hardware/en/lichee/th1520/lpi4a/7_develop.html#U-boot
const off int64 = 29360128 // 0x1c00000

func main() {
	adapter.Execute(&licheePi4AInstaller{})
}

type licheePi4AInstaller struct{}

type licheePi4AExtraOptions struct{}

func (i *licheePi4AInstaller) GetOptions(extra licheePi4AExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=ttyS0,115200n8",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	return overlay.Options{
		Name:       "licheepi-4a",
		KernelArgs: kernelArgs,
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *licheePi4AInstaller) Install(options overlay.InstallOptions[licheePi4AExtraOptions]) error {
	var f *os.File

	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", options.InstallDisk, err)
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(filepath.Join(options.ArtifactsPath, "riscv64/u-boot/licheepi-4a/u-boot-dtb.bin"))
	if err != nil {
		return err
	}

	// we need an offset so can't use copy.File
	if _, err = f.WriteAt(uboot, off); err != nil {
		return err
	}

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
