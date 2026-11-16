#!/usr/bin/env python3
"""Build the portable plugin ZIP from an explicit, secret-free allowlist."""
import argparse
import json
from pathlib import Path
from zipfile import ZipFile, ZIP_DEFLATED, ZipInfo

root = Path(__file__).resolve().parents[1] / 'plugins/reminder'
parser = argparse.ArgumentParser()
parser.add_argument('output', type=Path, help='Output ZIP path')
args = parser.parse_args()
files = ['plugin.json', 'mcp.json', 'assets/icon.png']
for name in files[:2]:
    json.loads((root / name).read_text())
args.output.parent.mkdir(parents=True, exist_ok=True)
with ZipFile(args.output, 'w', compression=ZIP_DEFLATED) as archive:
    for name in files:
        info = ZipInfo(name, date_time=(2026, 1, 1, 0, 0, 0))
        info.external_attr = 0o100644 << 16
        info.compress_type = ZIP_DEFLATED
        archive.writestr(info, (root / name).read_bytes())
print(args.output)
