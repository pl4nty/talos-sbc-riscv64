# sbc-riscv64

This repo provides overlays for RISC-V Talos images. It's based on a fork so please [report issues here](https://github.com/pl4nty/talos) instead of upstream.

## Supported Overlays

| Overlay Name | Board                  | SoC         | Description                                                                                          |
| -------------| ---------------------- | ----------- | ---------------------------------------------------------------------------------------------------- |
| licheepi-4a  | Sipeed LicheePi 4A 8GB | TH1520      | Overlay for Sipeed LicheePi 4A 8GB model. 16GB model isn't supported - open an issue if you want it. |
| licheepi-3a  | Sipeed LicheePi 3A 8GB | SpaceMIT K1 | WIP. No upstream DTS yet, so U-Boot has to supply the fdt. U-Boot SPL is built from `spacemit_k1_defconfig` and still needs signing into an FSBL with SpacemiT's vendor tooling. |
| k3-pico-itx  | Sipeed K3 Pico-ITX     | SpaceMIT K3 | WIP. DTB is upstream, but U-Boot has no K3 support - needs vendor firmware to boot.                  |
