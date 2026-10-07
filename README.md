# sbc-riscv64

This repo provides overlays for RISC-V Talos images. It's based on a fork so please [report issues here](https://github.com/pl4nty/talos) instead of upstream.

## Supported Overlays

| Overlay Name | Board                  | SoC         | Description                                                                                          |
| -------------| ---------------------- | ----------- | ---------------------------------------------------------------------------------------------------- |
| licheepi-4a  | Sipeed LicheePi 4A 8GB | TH1520      | Overlay for Sipeed LicheePi 4A 8GB model. 16GB model isn't supported - open an issue if you want it. |
| licheepi-3a  | Sipeed LicheePi 3A 8GB | SpacemiT K1 | WIP. Builds a signed FSBL and an SPL with DDR init and SD/eMMC boot, but there's still no upstream LicheePi 3A DTS, so the fdt has to come from elsewhere. |
| k3-pico-itx  | Sipeed K3 Pico-ITX     | SpacemiT K3 | WIP. DTB is upstream and the K3 clock and reset drivers are carried as patches, but U-Boot has no K3 board support yet - needs vendor firmware to boot.                  |

## Installation

There is no Image Factory for riscv64, so each tagged
[release](https://github.com/pl4nty/talos-sbc-riscv64/releases) takes its place.
A release carries:

- `metal-<board>-riscv64.raw.xz`: a Talos disk image with the overlay, the same as `factory.talos.dev/image/<schematic>/<version>/metal-arm64.raw.xz` upstream.
- `<board>-<file>`: the bootloader files that you flash by hand.
- `sbc-riscv64-artifacts.tar.gz`: every bootloader and firmware file in the overlay.
- `sha256sum.txt`: checksums for the files above.

Each release also pushes `ghcr.io/pl4nty/installer-<board>:<tag>` for `talosctl upgrade`, and the overlay itself as `ghcr.io/pl4nty/sbc-riscv64:<tag>` for use with `imager --overlay-image`.

Flashing guides:

- [LicheePi 4A](docs/licheepi-4a.md)
- [LicheePi 3A](docs/licheepi-3a.md) (WIP)

## Artifacts

CI uploads images and U-Boot, SPL and FSBL files from every run as workflow artifacts. U-Boot's own flashing docs:
[K1](https://docs.u-boot.org/en/latest/board/spacemit/k1-spl.html),
[TH1520](https://docs.u-boot.org/en/latest/board/thead/lpi4a.html),
[JH7110](https://docs.u-boot.org/en/latest/board/starfive/visionfive2.html).
