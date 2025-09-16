# Keychain Provider System

**Fully automatic** OTP retrieval from password managers (1Password, Bitwarden, etc.)

## Quick Start

1. Install password manager CLI: `brew install --cask 1password-cli`
2. Authenticate: `op account add`
3. Connect to OTP-enabled VPN profile - **OTP retrieved automatically!**

## Automatic Behavior

- **CLI**: Tries keychain providers automatically, falls back to manual entry if none work
- **Client**: Auto-retrieval when dialog opens, manual option available as backup
- **No Prompts**: Uses sensible defaults, no user interaction required
- **Default References**: `Private/Pritunl` (1Password), `Pritunl` (Bitwarden)

## Control Options

```bash
# Disable all keychain providers
pritunl-cli start profile123 --disable-keychain

# Use specific providers only  
pritunl-cli start profile123 --keychain-providers=1password
```

**Settings**: Advanced Settings → "Disable ALL keychain providers"

## Configuration

Override default references via config:
```json
{
  "keychain_default_refs": {
    "1password": "Work/VPN",
    "bitwarden": "VPN Server"
  }
}
```

## Documentation

📖 **Complete documentation is in the code** - see source files:

- **System overview**: `client/app/utils/KeychainProviders.ts`
- **Core implementation**: `cli/keychain/keychain.go`
- **Provider examples**: `cli/keychain/providers/*.go`

## Adding Providers

See code documentation for complete implementation examples. Adding a new provider requires ~50 lines of code total.
