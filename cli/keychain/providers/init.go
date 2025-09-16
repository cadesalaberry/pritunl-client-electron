/*
Package providers init registers all available keychain providers.

PROVIDER REGISTRATION
=====================

This file automatically registers all implemented keychain providers
with the global registry when the package is imported.

CURRENT PROVIDERS:
- OnePasswordProvider: 1Password CLI integration
- BitwardenProvider: Bitwarden CLI integration

ADDING NEW PROVIDERS:
1. Create new provider file (e.g., lastpass.go)
2. Implement the Provider interface
3. Add registration line below
4. Provider is automatically available throughout the system

EXAMPLE:

	keychain.GlobalRegistry.Register(&NewProvider{})

The init() function runs automatically when the package is imported,
ensuring all providers are registered before they're needed.
*/
package providers

import (
	"github.com/pritunl/pritunl-client-electron/cli/keychain"
)

func init() {
	// Register all available keychain providers
	// Providers are automatically available throughout the system once registered
	keychain.GlobalRegistry.Register(&OnePasswordProvider{})
	keychain.GlobalRegistry.Register(&BitwardenProvider{})
	
	// To add a new provider:
	// 1. Implement the Provider interface in a new file
	// 2. Add registration line here: keychain.GlobalRegistry.Register(&NewProvider{})
}
