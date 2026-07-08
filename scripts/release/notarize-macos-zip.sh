#!/usr/bin/env bash

# Packages and notarizes the signed macOS CLI binary.
#
# Usage:
#   scripts/release/notarize-macos-zip.sh \
#     dist/embed-code-macos \
#     dist/embed-code-macos.zip
#
# The script expects these environment variables:
#   APPLE_ID
#   APPLE_TEAM_ID
#   APPLE_APP_SPECIFIC_PASSWORD
#
# Apple notarization accepts an archive/package submission, so the CLI binary is
# wrapped in a ZIP instead of being uploaded as a loose executable.
#
# The script writes the ZIP path passed as its second argument and waits for
# Apple to accept or reject the notarization submission.
set -euo pipefail

# Verifies, that secrets are set.
required_vars=(
  APPLE_ID
  APPLE_TEAM_ID
  APPLE_APP_SPECIFIC_PASSWORD
)

missing=()
for var_name in "${required_vars[@]}"; do
  if [[ -z "${!var_name:-}" ]]; then
    missing+=("$var_name")
  fi
done
if (( ${#missing[@]} > 0 )); then
  printf 'Missing required environment variable: %s\n' "${missing[@]}" >&2
  exit 1
fi

if (( $# != 2 )); then
  echo "Usage: $0 <signed-macos-binary-path> <zip-output-path>" >&2
  exit 1
fi

binary_path="$1"
zip_path="$2"

if [[ ! -f "$binary_path" ]]; then
  echo "Signed macOS binary does not exist: $binary_path" >&2
  exit 1
fi

binary_dir="$(cd "$(dirname "$binary_path")" && pwd)"
binary_name="$(basename "$binary_path")"
zip_dir="$(dirname "$zip_path")"
zip_name="$(basename "$zip_path")"

# Resolve the ZIP path before changing directories.
mkdir -p "$zip_dir"
zip_dir="$(cd "$zip_dir" && pwd)"
zip_path="$zip_dir/$zip_name"

pushd "$binary_dir" >/dev/null
ditto -c -k --keepParent "$binary_name" "$zip_path"
popd >/dev/null

# Wait for the notarization result before publishing.
# ZIP archives do not support stapling, so the release publishes the accepted archive as-is.
xcrun notarytool submit "$zip_path" \
  --apple-id "$APPLE_ID" \
  --team-id "$APPLE_TEAM_ID" \
  --password "$APPLE_APP_SPECIFIC_PASSWORD" \
  --wait
