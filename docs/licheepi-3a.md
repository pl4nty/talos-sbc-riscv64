# Sipeed LicheePi 3A

Overlay: `licheepi-3a` (SpacemiT K1, 8 GB RAM model).

> **Work in progress.** These steps write a working FSBL and SPL to the
> board, but Talos has not yet booted end to end on the 3A.

The K1 boot ROM tries the SD card first, then the eMMC. The 3A has no SPI
NOR. The eMMC `boot0` hardware partition holds the boot chain, and Talos
takes the whole eMMC user area. The 3A is the K1 reference design, so it
uses the upstream Banana Pi BPI-F3 DTB, because no upstream 3A DTS exists
yet.

`licheepi-3a-boot0.img` is one image for `boot0`:

| Offset  | Contents                    |
| ------- | --------------------------- |
| 0       | bootinfo header             |
| 0x200   | FSBL (SPL + DDR init)       |
| 1 MiB   | `u-boot.itb` (OpenSBI + U-Boot) |

## Prerequisites

- `talosctl` on your workstation.
- The UART serial console at 115200 8N1.
- The vendor OS (Bianbu) running from the eMMC, as the board ships.
- Two downloads from the
  [releases page](https://github.com/pl4nty/talos-sbc-riscv64/releases),
  copied to the board:
  - `metal-licheepi-3a-riscv64.raw.xz`
  - `licheepi-3a-boot0.img`

Check the downloads against `sha256sum.txt` from the same release.

## Flash from Bianbu

This overwrites the OS you are running from. Keep the files in RAM
(`/dev/shm`), and power-cycle as soon as the write ends. Find the eMMC with
`lsblk` (shown here as `mmcblkN`).

```bash
cp metal-licheepi-3a-riscv64.raw.xz licheepi-3a-boot0.img /dev/shm/
cd /dev/shm

# boot chain
echo 0 | sudo tee /sys/block/mmcblkNboot0/force_ro
sudo dd if=licheepi-3a-boot0.img of=/dev/mmcblkNboot0 conv=fsync

# Talos
xz -dc metal-licheepi-3a-riscv64.raw.xz | sudo dd of=/dev/mmcblkN bs=4M conv=fsync status=progress
```

Then remove power and connect it again.

## Recovery

Hold `BOOT` and connect the power USB-C port (the K1 OTG port) to a host.
This starts the boot ROM's USB download mode. From there, the vendor U-Boot
can be loaded into RAM with `fastboot` and used to write the eMMC again.
Upstream U-Boot has no K1 USB support yet. See
[U-Boot's K1 SPL docs](https://docs.u-boot.org/en/latest/board/spacemit/k1-spl.html).

## Boot

The serial console shows Talos boot to maintenance mode. Then:

```bash
talosctl apply-config --insecure --mode=interactive --nodes <node IP>
```

## Upgrade

```bash
talosctl upgrade --nodes <node IP> --image ghcr.io/pl4nty/installer-licheepi-3a:<release tag>
```

The installer updates the DTBs on the EFI partition. It does not write
`boot0`.
