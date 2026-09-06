# sbc-riscv64

This repo provides overlays for RISC-V Talos images. It's based on a fork so please [report issues here](https://github.com/pl4nty/talos) instead of upstream.

## Supported Overlays

| Overlay Name | Board                  | SoC         | Description                                                                                          |
| -------------| ---------------------- | ----------- | ---------------------------------------------------------------------------------------------------- |
| licheepi-4a  | Sipeed LicheePi 4A 8GB | TH1520      | Overlay for Sipeed LicheePi 4A 8GB model. 16GB model isn't supported - open an issue if you want it. |
| licheepi-3a  | Sipeed LicheePi 3A 8GB | SpacemiT K1 | WIP. Builds a signed FSBL and an SPL with DDR init and SD/eMMC boot, but there's still no upstream LicheePi 3A DTS, so the fdt has to come from elsewhere. |
| k3-pico-itx  | Sipeed K3 Pico-ITX     | SpacemiT K3 | WIP. DTB is upstream and the K3 clock and reset drivers are carried as patches, but U-Boot has no K3 board support yet - needs vendor firmware to boot.                  |

## Artifacts

CI publishes U-Boot, SPL and FSBL images for some boards. Flashing:
[K1](https://docs.u-boot.org/en/latest/board/spacemit/k1-spl.html),
[TH1520](https://docs.u-boot.org/en/latest/board/thead/lpi4a.html),
[JH7110](https://docs.u-boot.org/en/latest/board/starfive/visionfive2.html).
