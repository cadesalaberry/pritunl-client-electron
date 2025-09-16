/*
BITWARDEN PROVIDER
==================

Integrates with Bitwarden CLI to retrieve OTP codes.

PREREQUISITES:
- Bitwarden CLI installed: npm install -g @bitwarden/cli
- User logged in: bw login
- Session unlocked: bw unlock

REFERENCE FORMAT:
- item-name (simple item name format)

EXAMPLES:
- GitHub
- AWS Account
- VPN Server

CLI COMMANDS USED:
- bw status (check authentication and unlock status)
- bw get totp "item-name" (get OTP)

NOTE: Bitwarden uses a simpler reference format compared to 1Password
since it doesn't have the same vault/item/field hierarchy.
*/
package providers

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pritunl/pritunl-client-electron/cli/keychain"
)

// BitwardenProvider implements Bitwarden CLI integration for OTP retrieval
type BitwardenProvider struct {
	keychain.BaseProvider
}

func (p *BitwardenProvider) ID() string {
	return "bitwarden"
}

func (p *BitwardenProvider) Name() string {
	return "Bitwarden"
}

func (p *BitwardenProvider) Description() string {
	return "Bitwarden CLI integration for OTP retrieval"
}

func (p *BitwardenProvider) IsAvailable() bool {
	_, err := exec.LookPath("bw")
	return err == nil
}

func (p *BitwardenProvider) TestAuthentication() bool {
	cmd := exec.Command("bw", "status")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	var status struct {
		Status string `json:"status"`
	}
	
	err = json.Unmarshal(output, &status)
	if err != nil {
		return false
	}
	
	return status.Status == "unlocked"
}

func (p *BitwardenProvider) ParseReference(reference string) (*keychain.ItemRef, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return nil, fmt.Errorf("empty reference")
	}

	return &keychain.ItemRef{
		Provider: p.ID(),
		Vault:    "", // Bitwarden doesn't use vault concept in the same way
		Item:     reference,
		Field:    "totp",
	}, nil
}

func (p *BitwardenProvider) GetOTP(ref *keychain.ItemRef) (string, error) {
	// Execute the Bitwarden CLI command to get TOTP
	cmd := exec.Command("bw", "get", "totp", ref.Item)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("Bitwarden CLI error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute Bitwarden CLI: %v", err)
	}

	otpCode := strings.TrimSpace(string(output))

	// Validate that we got a numeric OTP code
	if !p.ValidateOTPFormat(otpCode) {
		return "", fmt.Errorf("invalid OTP format received: %s", otpCode)
	}

	return otpCode, nil
}

func (p *BitwardenProvider) GetExamples() []string {
	return []string{
		"GitHub",
		"AWS Account", 
		"VPN Server",
	}
}

func (p *BitwardenProvider) GetHelpText() string {
	return "Format: item-name (e.g., GitHub or AWS Account)"
}
