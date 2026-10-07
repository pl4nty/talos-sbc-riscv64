# Sipeed LicheePi 4A

Overlay: `licheepi-4a` (TH1520, 8 GB RAM model only).

The board's U-Boot lives in the eMMC `boot0` hardware partition, and Talos
takes the whole eMMC user area. Both are written from a host PC over the
USB-C port, using the vendor U-Boot as a USB mass-storage (UMS) gadget. Our
U-Boot has no USB gadget support, so it can't flash itself.

## Prerequisites

- `talosctl`, `fastboot` (android-tools) and `xz` on the host.
- A USB-C cable from the host to the board's USB-C port.
- The UART0 serial console at 115200 8N1. You need it to type U-Boot
  commands during flashing.
- Two downloads from the
  [releases page](https://github.com/pl4nty/talos-sbc-riscv64/releases):
  - `metal-licheepi-4a-riscv64.raw.xz`: the Talos disk image.
  - `licheepi-4a-u-boot-with-spl.bin`: our U-Boot, with SPL, DDR firmware
    and OpenSBI.
- The vendor U-Boot, used only as the flasher:
  `u-boot-with-spl-lpi4a.bin` from
  [revyos/thead-u-boot releases](https://github.com/revyos/thead-u-boot/releases).
  For a 16 GB board, use `u-boot-with-spl-lpi4a-16g.bin`.

Check the downloads against `sha256sum.txt` from the same release.

## Load the vendor U-Boot into RAM

The boot ROM's fastboot loads into RAM only. It never writes the eMMC.

1. Hold `BOOT`, connect the USB-C cable to the host, then release `BOOT`.
   `lsusb` shows `2345:7654`.
2. On the host:

   ```bash
   fastboot flash ram u-boot-with-spl-lpi4a.bin
   fastboot reboot
   ```

3. On the serial console, press a key to stop autoboot.

On every start, the vendor U-Boot rewrites the eMMC GPT with the RevyOS
layout. Always write the Talos image last.

Use one `ums` target per boot. The vendor U-Boot hangs if you change from
`boot0` to the user area in the same session. Do steps 1-3 again before each
of the two writes below.

## Write U-Boot to eMMC boot0

Do this once, and again only when you update U-Boot. `talosctl upgrade` does
not write `boot0`.

```text
=> ums 0 mmc 0.1
```

A 4 MiB disk (`/dev/sdX`) comes up on the host. Then, on the host:

```bash
sudo dd if=licheepi-4a-u-boot-with-spl.bin of=/dev/sdX bs=1M conv=fsync
```

## Write Talos to the eMMC

Load the vendor U-Boot again (steps 1-3 above), then:

```text
=> ums 0 mmc 0
```

On the host:

```bash
xz -dc metal-licheepi-4a-riscv64.raw.xz | sudo dd of=/dev/sdX bs=4M conv=fsync status=progress
```

UMS on the vendor U-Boot can drop off the bus during a write. If the write
fails, start again from step 1.

## Boot

Disconnect the USB-C cable and power the board from its normal supply. The
serial console shows Talos boot to maintenance mode. Then:

```bash
talosctl apply-config --insecure --mode=interactive --nodes <node IP>
```

## Upgrade

Each release pushes an installer with the overlay included, the same as
`factory.talos.dev/installer/<schematic>` upstream:

```bash
talosctl upgrade --nodes <node IP> --image ghcr.io/pl4nty/installer-licheepi-4a:<release tag>
```

The installer updates the DTB and AON firmware on the EFI partition. It
does not write U-Boot to `boot0`.

## Troubleshooting

- Serial console: UART0, 115200 8N1. The overlay adds `console=ttyS0,115200`.
- A board with no output after the boot ROM usually has a missing or old
  U-Boot in `boot0`. Write it again.
- Ethernet and the eMMC need the kernel DTB. Our U-Boot loads it from
  `EFI/dtb/thead/th1520-lichee-pi-4a.dtb` on the EFI partition. A vendor
  U-Boot in `boot0` passes its own old DTB instead, and then GMAC is missing
  and the eMMC gives ADMA errors.
