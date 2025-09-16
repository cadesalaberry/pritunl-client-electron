/*
Package keychain provides a generic, extensible system for integrating with
password managers and keychain services to automatically retrieve OTP codes.

KEYCHAIN PROVIDER SYSTEM
========================

This package enables automatic OTP (One-Time Password) retrieval from multiple
password managers like 1Password, Bitwarden, LastPass, etc.

FEATURES:
- Multi-provider support with automatic detection
- Provider-specific reference formats and CLI integration
- Granular control: disable all providers or select specific ones
- Automatic provider registration and discovery
- Graceful fallback to manual entry

USAGE:
1. Install supported password manager CLI (e.g., 'op' for 1Password)
2. Authenticate with the CLI
3. Use Pritunl CLI with OTP-enabled profiles
4. System automatically detects providers and retrieves OTP

CONTROL OPTIONS:
- CLI: --disable-keychain flag
- CLI: --keychain-providers=provider1,provider2 flag

ADDING NEW PROVIDERS:
1. Implement Provider interface
2. Register in providers/init.go
3. Provider automatically integrates with entire system

Example:

	provider, otpCode, err := keychain.GetOTPFromAny("Private/GitHub")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got OTP %s from %s\n", otpCode, provider.Name())
*/
package keychain

import (
	"fmt"
	"regexp"
	"strings"
)

// ItemRef represents a structured keychain item reference
type ItemRef struct {
	Provider string // Provider ID (e.g., "1password", "bitwarden")
	Vault    string // Vault/folder name (empty for providers without vaults)
	Item     string // Item/entry name
	Field    string // Field name containing OTP (e.g., "one-time password", "totp")
}

/*
Provider interface that all keychain providers must implement.

IMPLEMENTATION GUIDE:
- Embed BaseProvider for common functionality
- Implement all required methods with proper error handling
- Use provider-specific CLI commands in GetOTP()
- Provide clear examples and help text for users
- Validate OTP format (6-8 digits) before returning

EXAMPLE PROVIDERS:
- 1Password: Uses 'op read "op://vault/item/field?attribute=otp"'
- Bitwarden: Uses 'bw get totp "item-name"'
- LastPass: Uses 'lpass show --field=totp "item-name"'

TESTING:
- IsAvailable() should check if CLI tool exists in PATH
- TestAuthentication() should verify user is logged in
- ParseReference() should validate format and return error for invalid input
- GetOTP() should execute CLI and validate returned OTP format
*/
type Provider interface {
	// ID returns the unique identifier for this provider (e.g., "1password")
	ID() string
	
	// Name returns the human-readable name (e.g., "1Password")
	Name() string
	
	// Description returns a brief description of the provider
	Description() string
	
	// IsAvailable checks if this provider's CLI tool is available on the system
	IsAvailable() bool
	
	// TestAuthentication tests if the provider is authenticated and ready to use
	TestAuthentication() bool
	
	// ParseReference parses a provider-specific reference string into structured format
	ParseReference(reference string) (*ItemRef, error)
	
	// GetOTP retrieves an OTP code from the provider using its CLI
	GetOTP(ref *ItemRef) (string, error)
	
	// FormatReference formats a reference for display/storage
	FormatReference(ref *ItemRef) string
	
	// GetExamples returns example reference formats for this provider
	GetExamples() []string
	
	// GetHelpText returns help text explaining the reference format
	GetHelpText() string
}

/*
BaseProvider provides common functionality for keychain providers.

USAGE:
Embed this struct in your provider implementation:

	type MyProvider struct {
		keychain.BaseProvider
	}

COMMON FUNCTIONALITY PROVIDED:
- FormatReference(): Standard vault/item/field formatting
- ValidateOTPFormat(): Validates 6-8 digit OTP codes
*/
type BaseProvider struct{}

// FormatReference provides standard formatting for vault/item/field references.
// Override this method if your provider uses a different format.
func (p *BaseProvider) FormatReference(ref *ItemRef) string {
	if ref.Field != "" && ref.Field != "one-time password" {
		return fmt.Sprintf("%s/%s/%s", ref.Vault, ref.Item, ref.Field)
	}
	return fmt.Sprintf("%s/%s", ref.Vault, ref.Item)
}

// ValidateOTPFormat validates that an OTP code is in the correct format (6-8 digits).
// Call this in your GetOTP() implementation before returning the code.
func (p *BaseProvider) ValidateOTPFormat(otpCode string) bool {
	matched, _ := regexp.MatchString(`^\d{6,8}$`, otpCode)
	return matched
}

/*
Registry manages all keychain providers with automatic discovery and filtering.

FUNCTIONALITY:
- Registers all available keychain providers
- Filters providers by availability and authentication
- Auto-detects which provider can handle a reference format
- Provides helper functions for common operations

USAGE:
	// Register a new provider
	GlobalRegistry.Register(&MyProvider{})
	
	// Get all available providers
	providers := GetAvailableProviders()
	
	// Get OTP from any provider
	provider, otp, err := GetOTPFromAny("Private/GitHub")
*/
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(provider Provider) {
	r.providers[provider.ID()] = provider
}

// Get retrieves a specific provider by ID
func (r *Registry) Get(id string) Provider {
	return r.providers[id]
}

// GetAll returns all registered providers
func (r *Registry) GetAll() []Provider {
	providers := make([]Provider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}

// GetAvailable returns only providers that are available on the system
func (r *Registry) GetAvailable() []Provider {
	var available []Provider
	for _, provider := range r.providers {
		if provider.IsAvailable() {
			available = append(available, provider)
		}
	}
	return available
}

// FindByReference finds the first provider that can parse the given reference
func (r *Registry) FindByReference(reference string) Provider {
	for _, provider := range r.providers {
		_, err := provider.ParseReference(reference)
		if err == nil {
			return provider
		}
	}
	return nil
}

// GlobalRegistry is the global instance used throughout the application.
// Providers are automatically registered via init() functions in provider packages.
var GlobalRegistry = NewRegistry()

// DefaultReferences provides sensible defaults for each provider.
// These can be overridden via configuration.
var DefaultReferences = map[string]string{
	"1password": "Private/Pritunl",  // Default 1Password reference
	"bitwarden":  "Pritunl",         // Default Bitwarden reference
}

/*
Helper functions for common keychain operations.

EXAMPLES:
	// Get all available providers
	providers := GetAvailableProviders()
	
	// Automatically get OTP from any provider
	provider, otpCode, err := GetOTPFromAny("Private/GitHub")
	if err != nil {
		log.Printf("Failed to get OTP: %v", err)
		return
	}
	fmt.Printf("Got OTP %s from %s\n", otpCode, provider.Name())
	
	// Parse reference to determine provider
	provider, ref, err := ParseAnyReference("Private/GitHub")
	if err != nil {
		log.Printf("Invalid reference: %v", err)
		return
	}
	fmt.Printf("Reference %s will use %s\n", reference, provider.Name())
*/

// GetAvailableProviders returns all providers that are available on the system
func GetAvailableProviders() []Provider {
	return GlobalRegistry.GetAvailable()
}

/*
GetOTPFromAny automatically retrieves OTP from any provider based on reference format.

The function:
1. Determines which provider can handle the reference
2. Validates the reference format
3. Checks provider availability and authentication
4. Retrieves and validates the OTP code

USAGE:
	provider, otpCode, err := GetOTPFromAny("Private/GitHub")
	if err != nil {
		// Handle error - provider not available, not authenticated, or retrieval failed
		return err
	}
	// Use otpCode for authentication

ERRORS:
- "no provider found": No provider can parse the reference format
- "invalid reference format": Reference format is invalid for the detected provider
- "[Provider] is not available": Provider CLI is not installed
- "[Provider] is not authenticated": User is not logged in to provider
- Provider-specific errors: CLI execution or OTP retrieval failed
*/
func GetOTPFromAny(reference string) (Provider, string, error) {
	provider := GlobalRegistry.FindByReference(reference)
	if provider == nil {
		return nil, "", fmt.Errorf("no provider found that can handle this reference format")
	}
	
	ref, err := provider.ParseReference(reference)
	if err != nil {
		return nil, "", fmt.Errorf("invalid reference format: %v", err)
	}
	
	if !provider.IsAvailable() {
		return nil, "", fmt.Errorf("%s is not available", provider.Name())
	}
	
	if !provider.TestAuthentication() {
		return nil, "", fmt.Errorf("%s is not authenticated", provider.Name())
	}
	
	otpCode, err := provider.GetOTP(ref)
	if err != nil {
		return nil, "", err
	}
	
	return provider, otpCode, nil
}

/*
ParseAnyReference parses any reference format to determine provider and structured item.

USAGE:
	provider, ref, err := ParseAnyReference("Private/GitHub")
	if err != nil {
		log.Printf("Invalid reference: %v", err)
		return
	}
	fmt.Printf("Provider: %s, Vault: %s, Item: %s\n", 
		provider.Name(), ref.Vault, ref.Item)
*/
func ParseAnyReference(reference string) (Provider, *ItemRef, error) {
	provider := GlobalRegistry.FindByReference(reference)
	if provider == nil {
		return nil, nil, fmt.Errorf("no provider found that can handle this reference format")
	}
	
	ref, err := provider.ParseReference(reference)
	if err != nil {
		return nil, nil, err
	}
	
	return provider, ref, nil
}

/*
TryAutoGetOTP attempts to automatically retrieve OTP from available providers.

This function tries each available provider with its default reference or
configured reference, returning the first successful OTP retrieval.

USAGE:
	provider, otpCode, err := TryAutoGetOTP(configDefaultRefs, enabledProviderIds)
	if err != nil {
		// No automatic OTP available, fall back to manual entry
		return manual_prompt()
	}
	// Use otpCode for authentication

PARAMETERS:
- configDefaultRefs: Map of provider_id -> default_reference from configuration
- enabledProviderIds: List of enabled provider IDs (empty = all providers)

RETURNS:
- Provider that provided the OTP
- OTP code
- Error if no providers available or all failed
*/
func TryAutoGetOTP(configDefaultRefs map[string]string, enabledProviderIds []string) (Provider, string, error) {
	providers := GetAvailableProviders()
	
	// Filter providers based on enabled list if specified
	if len(enabledProviderIds) > 0 {
		enabledSet := make(map[string]bool)
		for _, id := range enabledProviderIds {
			enabledSet[strings.TrimSpace(id)] = true
		}
		
		var filteredProviders []Provider
		for _, provider := range providers {
			if enabledSet[provider.ID()] {
				filteredProviders = append(filteredProviders, provider)
			}
		}
		providers = filteredProviders
	}
	
	// Try each available provider
	for _, provider := range providers {
		// Get reference for this provider (config override or default)
		var reference string
		if configDefaultRefs != nil {
			reference = configDefaultRefs[provider.ID()]
		}
		if reference == "" {
			reference = DefaultReferences[provider.ID()]
		}
		if reference == "" {
			continue // No reference available for this provider
		}
		
		// Try to get OTP from this provider
		providerUsed, otpCode, err := GetOTPFromAny(reference)
		if err == nil {
			return providerUsed, otpCode, nil
		}
		// Continue to next provider if this one failed
	}
	
	return nil, "", fmt.Errorf("no keychain providers could automatically provide OTP")
}
