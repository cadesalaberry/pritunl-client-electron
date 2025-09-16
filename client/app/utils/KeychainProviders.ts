/// <reference path="../References.d.ts"/>
import * as Logger from "../Logger"

/**
 * KEYCHAIN PROVIDER SYSTEM
 * ========================
 * 
 * This system provides automatic OTP (One-Time Password) retrieval from multiple
 * password managers and keychain services like 1Password, Bitwarden, etc.
 * 
 * FEATURES:
 * - Automatic OTP retrieval from supported password managers
 * - Multi-provider support with automatic detection
 * - Provider-specific reference formats and help text
 * - Granular control: disable all providers or select specific ones
 * - Automatic reference storage for future use
 * - Graceful fallback to manual entry
 * 
 * USAGE:
 * 1. Install supported password manager CLI (e.g., 'op' for 1Password)
 * 2. Authenticate with the CLI
 * 3. Connect to OTP-enabled VPN profile
 * 4. System automatically detects providers and retrieves OTP
 * 
 * CONTROL OPTIONS:
 * - Settings: "Disable ALL keychain providers" toggle
 * - CLI: --disable-keychain flag
 * - CLI: --keychain-providers=provider1,provider2 flag
 * 
 * ADDING NEW PROVIDERS:
 * 1. Implement KeychainProvider interface
 * 2. Register in providers/index.ts
 * 3. Provider automatically integrates with entire system
 */

export interface KeychainItem {
	provider: string
	vault: string
	item: string
	field?: string
}

/**
 * KeychainProvider interface that all password manager integrations must implement.
 * 
 * IMPLEMENTATION GUIDE:
 * - Extend BaseKeychainProvider for common functionality
 * - Implement all required methods with proper error handling
 * - Use provider-specific CLI commands in getOTP()
 * - Provide clear examples and help text for users
 * - Validate OTP format (6-8 digits) before returning
 * 
 * EXAMPLE PROVIDERS:
 * - 1Password: Uses 'op read "op://vault/item/field?attribute=otp"'
 * - Bitwarden: Uses 'bw get totp "item-name"'
 * - LastPass: Uses 'lpass show --field=totp "item-name"'
 */
export interface KeychainProvider {
	/** Unique identifier for this provider (e.g., "1password", "bitwarden") */
	readonly id: string
	
	/** Human-readable name (e.g., "1Password", "Bitwarden") */
	readonly name: string
	
	/** Brief description of the provider */
	readonly description: string
	
	/**
	 * Check if this provider's CLI tool is available on the system.
	 * Should check if the CLI command exists in PATH.
	 * 
	 * @returns Promise<boolean> true if CLI is available
	 */
	isAvailable(): Promise<boolean>
	
	/**
	 * Test if the provider is authenticated and ready to use.
	 * Should verify the user is logged in to the provider.
	 * 
	 * @returns Promise<boolean> true if authenticated
	 */
	testAuthentication(): Promise<boolean>
	
	/**
	 * Parse a provider-specific reference string into structured components.
	 * Should validate format and return null for invalid references.
	 * 
	 * @param reference Provider-specific reference string
	 * @returns KeychainItem | null Parsed item or null if invalid
	 */
	parseReference(reference: string): KeychainItem | null
	
	/**
	 * Retrieve OTP code from the provider using its CLI.
	 * Should execute provider CLI and return current OTP code.
	 * Must validate OTP format (6-8 digits) before returning.
	 * 
	 * @param item Parsed keychain item reference
	 * @returns Promise<string> Current OTP code
	 * @throws Error if retrieval fails or OTP format is invalid
	 */
	getOTP(item: KeychainItem): Promise<string>
	
	/**
	 * Format a reference string for display/storage.
	 * Default implementation handles common vault/item/field format.
	 * 
	 * @param item Keychain item to format
	 * @returns string Formatted reference
	 */
	formatReference(item: KeychainItem): string
	
	/**
	 * Get example reference formats for this provider.
	 * Used in UI to show users valid reference examples.
	 * 
	 * @returns string[] Array of example references
	 */
	getExamples(): string[]
	
	/**
	 * Get help text explaining the reference format for this provider.
	 * Shown in UI to guide users on proper reference format.
	 * 
	 * @returns string Help text for reference format
	 */
	getHelpText(): string
}

/**
 * Base class for keychain providers with common functionality.
 * 
 * USAGE:
 * Extend this class and implement the abstract methods:
 * 
 * ```typescript
 * export class MyProvider extends BaseKeychainProvider {
 *   readonly id = "myprovider"
 *   readonly name = "My Provider"
 *   readonly description = "My Provider CLI integration"
 *   
 *   async isAvailable(): Promise<boolean> {
 *     // Check if CLI is installed: which my-cli
 *   }
 *   
 *   async getOTP(item: KeychainItem): Promise<string> {
 *     // Execute CLI and return OTP: my-cli get-otp item.item
 *   }
 *   
 *   // ... implement other methods
 * }
 * ```
 * 
 * COMMON FUNCTIONALITY PROVIDED:
 * - formatReference(): Standard vault/item/field formatting
 * - validateOTPFormat(): Validates 6-8 digit OTP codes
 * - Logging helpers: logInfo(), logWarning(), logError()
 */
export abstract class BaseKeychainProvider implements KeychainProvider {
	abstract readonly id: string
	abstract readonly name: string
	abstract readonly description: string
	
	abstract isAvailable(): Promise<boolean>
	abstract testAuthentication(): Promise<boolean>
	abstract parseReference(reference: string): KeychainItem | null
	abstract getOTP(item: KeychainItem): Promise<string>
	abstract getExamples(): string[]
	abstract getHelpText(): string
	
	/**
	 * Default reference formatting for vault/item/field structure.
	 * Override if your provider uses a different format.
	 */
	formatReference(item: KeychainItem): string {
		if (item.field && item.field !== "one-time password") {
			return `${item.vault}/${item.item}/${item.field}`
		}
		return `${item.vault}/${item.item}`
	}
	
	/**
	 * Validate that OTP code is in correct format (6-8 digits).
	 * Call this in your getOTP() implementation before returning.
	 */
	protected validateOTPFormat(otpCode: string): boolean {
		return /^\d{6,8}$/.test(otpCode)
	}
	
	/** Log info message with provider name prefix */
	protected logInfo(message: string): void {
		Logger.info(`${this.name}: ${message}`)
	}
	
	/** Log warning message with provider name prefix */
	protected logWarning(message: string): void {
		Logger.warning(`${this.name}: ${message}`)
	}
	
	/** Log error message with provider name prefix */
	protected logError(message: string): void {
		Logger.error(`${this.name}: ${message}`)
	}
}

/**
 * Registry for keychain providers with automatic discovery and filtering.
 * 
 * FUNCTIONALITY:
 * - Registers all available keychain providers
 * - Filters providers by availability and authentication
 * - Supports provider selection via configuration
 * - Auto-detects which provider can handle a reference
 * 
 * USAGE:
 * ```typescript
 * // Register a new provider
 * Registry.register(new MyProvider())
 * 
 * // Get all available providers
 * const providers = await Registry.getAvailable()
 * 
 * // Get filtered providers
 * const filtered = await Registry.getAvailableFiltered(["1password"])
 * ```
 */
class KeychainProviderRegistry {
	private providers: Map<string, KeychainProvider> = new Map()
	
	/** Register a new keychain provider */
	register(provider: KeychainProvider): void {
		this.providers.set(provider.id, provider)
	}
	
	/** Get a specific provider by ID */
	get(id: string): KeychainProvider | undefined {
		return this.providers.get(id)
	}
	
	/** Get all registered providers */
	getAll(): KeychainProvider[] {
		return Array.from(this.providers.values())
	}
	
	getAvailable(): Promise<KeychainProvider[]> {
		return Promise.all(
			this.getAll().map(async (provider) => {
				const available = await provider.isAvailable()
				return available ? provider : null
			})
		).then(providers => providers.filter(p => p !== null) as KeychainProvider[])
	}
	
	getFiltered(enabledProviderIds?: string[]): KeychainProvider[] {
		if (!enabledProviderIds || enabledProviderIds.length === 0) {
			return this.getAll()
		}
		
		const enabledSet = new Set(enabledProviderIds)
		return this.getAll().filter(provider => enabledSet.has(provider.id))
	}
	
	getAvailableFiltered(enabledProviderIds?: string[]): Promise<KeychainProvider[]> {
		const filtered = this.getFiltered(enabledProviderIds)
		return Promise.all(
			filtered.map(async (provider) => {
				const available = await provider.isAvailable()
				return available ? provider : null
			})
		).then(providers => providers.filter(p => p !== null) as KeychainProvider[])
	}
	
	findByReference(reference: string): KeychainProvider | null {
		// Try to determine provider from reference format
		for (const provider of this.getAll()) {
			const item = provider.parseReference(reference)
			if (item) {
				return provider
			}
		}
		return null
	}
}

/** Global registry instance - automatically populated by provider imports */
export const Registry = new KeychainProviderRegistry()

/**
 * Get available keychain providers, optionally filtered by enabled list.
 * 
 * @param enabledProviderIds Optional array of provider IDs to filter by
 * @returns Promise<KeychainProvider[]> Available providers
 * 
 * USAGE:
 * ```typescript
 * // Get all available providers
 * const providers = await getAvailableProviders()
 * 
 * // Get only specific providers
 * const filtered = await getAvailableProviders(["1password", "bitwarden"])
 * ```
 */
export async function getAvailableProviders(enabledProviderIds?: string[]): Promise<KeychainProvider[]> {
	return Registry.getAvailableFiltered(enabledProviderIds)
}

/**
 * Automatically retrieve OTP from any provider based on reference format.
 * 
 * @param reference Provider-specific reference string
 * @returns Promise<{provider, otpCode}> Provider used and OTP code
 * @throws Error if no provider found, not available, not authenticated, or retrieval fails
 * 
 * USAGE:
 * ```typescript
 * try {
 *   const result = await getOTPFromAny("Private/GitHub")
 *   console.log(`Got OTP ${result.otpCode} from ${result.provider.name}`)
 * } catch (error) {
 *   console.error(`Failed to get OTP: ${error.message}`)
 * }
 * ```
 */
export async function getOTPFromAny(reference: string): Promise<{provider: KeychainProvider, otpCode: string}> {
	const provider = Registry.findByReference(reference)
	if (!provider) {
		throw new Error("No provider found that can handle this reference format")
	}
	
	const item = provider.parseReference(reference)
	if (!item) {
		throw new Error("Invalid reference format")
	}
	
	const available = await provider.isAvailable()
	if (!available) {
		throw new Error(`${provider.name} is not available`)
	}
	
	const authenticated = await provider.testAuthentication()
	if (!authenticated) {
		throw new Error(`${provider.name} is not authenticated`)
	}
	
	const otpCode = await provider.getOTP(item)
	return { provider, otpCode }
}

/**
 * Parse any reference format to determine provider and structured item.
 * 
 * @param reference Provider-specific reference string
 * @returns {provider, item} | null Parsed result or null if no provider can handle it
 * 
 * USAGE:
 * ```typescript
 * const result = parseAnyReference("Private/GitHub")
 * if (result) {
 *   console.log(`Provider: ${result.provider.name}, Item: ${result.item.item}`)
 * }
 * ```
 */
export function parseAnyReference(reference: string): {provider: KeychainProvider, item: KeychainItem} | null {
	const provider = Registry.findByReference(reference)
	if (!provider) {
		return null
	}
	
	const item = provider.parseReference(reference)
	if (!item) {
		return null
	}
	
	return { provider, item }
}
