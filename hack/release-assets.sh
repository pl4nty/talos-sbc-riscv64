#!/usr/bin/env bash

# Collects the bootloader files a user flashes by hand into flat release
# assets: <board>-<file> per board, plus a tarball of the whole artifacts tree.
# Upstream SBC overlays leave this to `crane export` of the overlay image;
# release assets save riscv64 users from installing crane.
#
# Usage: hack/release-assets.sh <artifacts dir> <out dir>

set -euo pipefail

src="${1:?artifacts dir}"
out="${2:?out dir}"

mkdir -p "${out}"

declare -A files=(
  [licheepi-4a]="riscv64/u-boot/licheepi-4a/u-boot-with-spl.bin"
  [licheepi-3a]="riscv64/fsbl/k1/boot0.img"
)

for board in "${!files[@]}"; do
  f="${src}/${files[${board}]}"

  # a single-board (BOARDS=) build only has its own board's files
  if [ -f "${f}" ]; then
    cp -v "${f}" "${out}/${board}-$(basename "${f}")"
  fi
done

tar -C "${src}" -czf "${out}/sbc-riscv64-artifacts.tar.gz" riscv64
