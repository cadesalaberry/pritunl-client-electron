/// <reference path="../../References.d.ts"/>
import { Registry } from "../KeychainProviders"
import { OnePasswordProvider } from "./OnePasswordProvider"
import { BitwardenProvider } from "./BitwardenProvider"

/*
KEYCHAIN PROVIDER REGISTRATION
==============================

This file registers all available keychain providers with the global registry
and exports the main keychain functionality for use throughout the client.

CURRENT PROVIDERS:
- OnePasswordProvider: 1Password CLI integration (vault/item format)
- BitwardenProvider: Bitwarden CLI integration (item-name format)

ADDING NEW PROVIDERS:
1. Create provider file implementing KeychainProvider interface
2. Import and register here: Registry.register(new YourProvider())
3. Export for direct access: export { YourProvider } from "./YourProvider"
4. Provider is automatically available throughout the client

USAGE:
	import { getAvailableProviders, getOTPFromAny } from "../utils/providers"
	
	// Get available providers
	const providers = await getAvailableProviders()
	
	// Get OTP from any provider
	const result = await getOTPFromAny("Private/GitHub")

The providers are automatically registered when this module is imported,
ensuring they're available throughout the client application.
*/

// Register all available providers with the global registry
Registry.register(new OnePasswordProvider())
Registry.register(new BitwardenProvider())

// Export providers for direct access if needed
export { OnePasswordProvider } from "./OnePasswordProvider"
export { BitwardenProvider } from "./BitwardenProvider"

// Re-export the registry and main keychain functionality
export { Registry, getAvailableProviders, getOTPFromAny, parseAnyReference } from "../KeychainProviders"
export type { KeychainProvider, KeychainItem } from "../KeychainProviders"
