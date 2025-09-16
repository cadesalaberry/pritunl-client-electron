# pritunl-client-electron: pritunl vpn client

[![package-macOS](https://img.shields.io/badge/package-macOS-cfcfcf.svg?style=flat)](https://github.com/pritunl/pritunl-client-electron/releases)
[![package-windows](https://img.shields.io/badge/package-windows-00adef.svg?style=flat)](https://github.com/pritunl/pritunl-client-electron/releases)
[![github](https://img.shields.io/badge/github-pritunl-11bdc2.svg?style=flat)](https://github.com/pritunl)
[![twitter](https://img.shields.io/badge/twitter-pritunl-55acee.svg?style=flat)](https://twitter.com/pritunl)
[![medium](https://img.shields.io/badge/medium-pritunl-b32b2b.svg?style=flat)](https://pritunl.medium.com)
[![forum](https://img.shields.io/badge/discussion-forum-ffffff.svg?style=flat)](https://forum.pritunl.com)

[Pritunl-client-electron](https://github.com/pritunl/pritunl-client-electron)
is an open source openvpn client. Documentation and more information can be
found at the home page [client.pritunl.com](https://client.pritunl.com)

## Install From Source (macOS)

If the Pritunl package is currently installed run the uninstall command
below. Requires homebrew with git, go and node.

```bash
brew install git go node
bash <(curl -s https://raw.githubusercontent.com/pritunl/pritunl-client-electron/master/tools/install_macos.sh)
```

## Uninstall From Source (macOS)

```bash
bash <(curl -s https://raw.githubusercontent.com/pritunl/pritunl-client-electron/master/tools/uninstall_macos.sh)
```

## Keychain Provider Integration

The Pritunl client includes automatic OTP retrieval from password managers like 1Password and Bitwarden.

**📖 [Complete Documentation](keychain/README.md)**

### Quick Start
- Install a password manager CLI (e.g., `brew install --cask 1password-cli`)
- Authenticate with the CLI
- Connect to OTP-enabled VPN profiles - keychain integration works automatically!

### Control Options
```bash
# Disable all keychain providers
pritunl-cli start profile123 --disable-keychain

# Use only specific providers
pritunl-cli start profile123 --keychain-providers=1password,bitwarden
```
