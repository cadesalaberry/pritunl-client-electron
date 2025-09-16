/*
Package providers contains implementations of keychain providers for various
password managers and keychain services.

1PASSWORD PROVIDER
==================

Integrates with 1Password CLI to retrieve OTP codes.

PREREQUISITES:
- 1Password CLI installed: brew install --cask 1password-cli
- 1Password desktop app running and authenticated
- CLI authenticated: op account add

REFERENCE FORMAT:
- vault/item (uses default "one-time password" field)
- vault/item/field (uses specific field)

EXAMPLES:
- Private/GitHub
- Work/AWS/one-time password
- Personal/VPN Server/authenticator

CLI COMMANDS USED:
- op account list (check authentication)
- op read "op://vault/item/field?attribute=otp" (get OTP)
*/
package providers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pritunl/pritunl-client-electron/cli/keychain"
)

// OnePasswordProvider implements 1Password CLI integration for OTP retrieval
type OnePasswordProvider struct {
	keychain.BaseProvider
}

func (p *OnePasswordProvider) ID() string {
	return "1password"
}

func (p *OnePasswordProvider) Name() string {
	return "1Password"
}

func (p *OnePasswordProvider) Description() string {
	return "1Password CLI integration for OTP retrieval"
}

func (p *OnePasswordProvider) IsAvailable() bool {
	_, err := exec.LookPath("op")
	return err == nil
}

func (p *OnePasswordProvider) TestAuthentication() bool {
	cmd := exec.Command("op", "account", "list")
	err := cmd.Run()
	return err == nil
}

func (p *OnePasswordProvider) ParseReference(reference string) (*keychain.ItemRef, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return nil, fmt.Errorf("empty reference")
	}

	parts := strings.Split(reference, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid reference format, use: vault/item or vault/item/field")
	}

	ref := &keychain.ItemRef{
		Provider: p.ID(),
		Vault:    parts[0],
		Item:     parts[1],
	}

	if len(parts) > 2 {
		ref.Field = strings.Join(parts[2:], "/")
	} else {
		ref.Field = "one-time password"
	}

	return ref, nil
}

func (p *OnePasswordProvider) GetOTP(ref *keychain.ItemRef) (string, error) {
	// Build the secret reference
	secretRef := fmt.Sprintf("op://%s/%s/%s?attribute=otp", ref.Vault, ref.Item, ref.Field)

	// Execute the 1Password CLI command
	cmd := exec.Command("op", "read", secretRef)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("1Password CLI error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute 1Password CLI: %v", err)
	}

	otpCode := strings.TrimSpace(string(output))

	// Validate that we got a numeric OTP code
	if !p.ValidateOTPFormat(otpCode) {
		return "", fmt.Errorf("invalid OTP format received: %s", otpCode)
	}

	return otpCode, nil
}

func (p *OnePasswordProvider) GetExamples() []string {
	return []string{
		"Private/GitHub",
		"Work/AWS/one-time password",
		"Personal/VPN Server/authenticator",
	}
}

func (p *OnePasswordProvider) GetHelpText() string {
	return "Format: vault/item or vault/item/field (e.g., Private/GitHub or Private/GitHub/one-time password)"
}
