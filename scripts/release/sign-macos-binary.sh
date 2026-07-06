#!/usr/bin/env bash

# Signs the macOS CLI binary with a Developer ID Application certificate.
#
# Usage:
#   scripts/release/sign-macos-binary.sh dist/embed-code-macos
#
# The script expects these environment variables:
#   MACOS_CERTIFICATE_P12_BASE64
#   MACOS_CERTIFICATE_PASSWORD
#   MACOS_CODESIGN_IDENTITY
#
# The certificate value must be a base64-encoded `.p12` export that contains the
# Developer ID Application certificate and private key.
# `MACOS_CODESIGN_IDENTITY` is the full identity printed by
# `security find-identity`, such as:
#
#   Developer ID Application: Company Name (TEAMID)
#
# For local diagnostics, `MACOS_KEYCHAIN_PASSWORD` can be set to reuse a known
# temporary keychain password.
#
# The script writes temporary certificate and keychain files under RUNNER_TEMP,
# or TMPDIR when RUNNER_TEMP is not set. It signs the binary in place and verifies
# the resulting signature.
set -euo pipefail

# Verifies, that secrets are set.
required_vars=(
  MACOS_CERTIFICATE_P12_BASE64
  MACOS_CERTIFICATE_PASSWORD
  MACOS_CODESIGN_IDENTITY
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

if (( $# != 1 )); then
  echo "Usage: $0 <macos-binary-path>" >&2
  exit 1
fi

binary_path="$1"
if [[ ! -f "$binary_path" ]]; then
  echo "macOS binary does not exist: $binary_path" >&2
  exit 1
fi

temp_dir="${RUNNER_TEMP:-${TMPDIR:-/tmp}}"
keychain_path="$temp_dir/embed-code-signing.keychain-db"
certificate_path="$temp_dir/embed-code-signing-certificate.p12"
keychain_password="${MACOS_KEYCHAIN_PASSWORD:-$(uuidgen)}"

# Decode the certificate at runtime so only the temporary runner filesystem ever
# contains the `.p12`.
if printf '%s' "$MACOS_CERTIFICATE_P12_BASE64" | base64 --decode > "$certificate_path" 2>/dev/null; then
  :
else
  printf '%s' "$MACOS_CERTIFICATE_P12_BASE64" | base64 -D > "$certificate_path"
fi

security create-keychain -p "$keychain_password" "$keychain_path"
security set-keychain-settings -lut 21600 "$keychain_path"
security unlock-keychain -p "$keychain_password" "$keychain_path"

# Import Apple's public certificate chain into the temporary keychain.
# Some runners do not have the current Developer ID intermediate certificates.
curl -fsSL https://www.apple.com/certificateauthority/AppleRootCA-G3.cer \
  -o "$temp_dir/AppleRootCA-G3.cer"
curl -fsSL https://www.apple.com/certificateauthority/AppleWWDRCAG6.cer \
  -o "$temp_dir/AppleWWDRCAG6.cer"
curl -fsSL https://www.apple.com/certificateauthority/DeveloperIDG2CA.cer \
  -o "$temp_dir/DeveloperIDG2CA.cer"

security import "$temp_dir/AppleRootCA-G3.cer" -k "$keychain_path"
security import "$temp_dir/AppleWWDRCAG6.cer" -k "$keychain_path"
security import "$temp_dir/DeveloperIDG2CA.cer" -k "$keychain_path"

# Import the private signing identity and allow Apple signing tools to access
# the key non-interactively.
security import "$certificate_path" \
  -P "$MACOS_CERTIFICATE_PASSWORD" \
  -A \
  -t cert \
  -f pkcs12 \
  -k "$keychain_path"
security list-keychains -d user -s "$keychain_path"
security set-key-partition-list \
  -S apple-tool:,apple: \
  -s \
  -k "$keychain_password" \
  "$keychain_path"
security find-identity -v -p codesigning "$keychain_path"

# Use the hardened runtime and a trusted timestamp so Gatekeeper can evaluate
# the signature after the Developer ID certificate eventually expires.
codesign \
  --force \
  --options runtime \
  --timestamp \
  --sign "$MACOS_CODESIGN_IDENTITY" \
  "$binary_path"
codesign --verify --verbose=4 "$binary_path"
