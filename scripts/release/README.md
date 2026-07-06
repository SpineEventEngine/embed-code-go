# Release Scripts

These scripts support the `release-binaries` GitHub workflow. 
They are intended for the macOS release job.

## Signing

Use `sign-macos-binary.sh` to import the Developer ID certificate into a
temporary keychain and sign the macOS binary.

Required environment variables:

- `MACOS_CERTIFICATE_P12_BASE64`: base64-encoded `.p12` export of the Developer
  ID Application certificate and private key.
- `MACOS_CERTIFICATE_PASSWORD`: password used when exporting the `.p12` file.
- `MACOS_CODESIGN_IDENTITY`: full Developer ID Application identity, such as
  `Developer ID Application: Company Name (TEAMID)`.

Optional environment variable:

- `MACOS_KEYCHAIN_PASSWORD`: password for the temporary keychain. When omitted,
  the script generates one.

Example:

```bash
scripts/release/sign-macos-binary.sh dist/embed-code-macos
```

## Notarization

Use `notarize-macos-zip.sh` to package the signed CLI binary as a ZIP archive
and submit that archive to Apple notarization.

Required environment variables:

- `APPLE_ID`: Apple Account email used for notarization.
- `APPLE_TEAM_ID`: 10-character Apple Developer Team ID.
- `APPLE_APP_SPECIFIC_PASSWORD`: app-specific password for the Apple Account.

Example:

```bash
scripts/release/notarize-macos-zip.sh \
  dist/embed-code-macos \
  dist/embed-code-macos.zip
```
