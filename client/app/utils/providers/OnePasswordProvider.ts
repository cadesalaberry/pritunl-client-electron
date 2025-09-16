/// <reference path="../../References.d.ts"/>
import * as child_process from "child_process"
import { BaseKeychainProvider, KeychainItem } from "../KeychainProviders"

/*
1PASSWORD PROVIDER
==================

Client-side implementation of 1Password CLI integration for OTP retrieval.

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

IMPLEMENTATION DETAILS:
- Uses Node.js child_process to execute 'op' CLI commands
- Validates OTP format before returning (6-8 digits)
- Provides detailed error messages for troubleshooting
- Supports custom field names for flexibility

USAGE:
This provider is automatically registered and available when 1Password CLI
is installed and authenticated. Users can reference items using the
vault/item or vault/item/field format.
*/
export class OnePasswordProvider extends BaseKeychainProvider {
	readonly id = "1password"
	readonly name = "1Password"
	readonly description = "1Password CLI integration for OTP retrieval"

	async isAvailable(): Promise<boolean> {
		return new Promise((resolve) => {
			child_process.exec("which op", (error) => {
				resolve(!error)
			})
		})
	}

	async testAuthentication(): Promise<boolean> {
		return new Promise((resolve) => {
			child_process.exec("op account list", {
				timeout: 5000,
			}, (error, stdout) => {
				if (error) {
					resolve(false)
					return
				}
				
				// If we can list accounts, we're authenticated
				resolve(stdout.trim().length > 0)
			})
		})
	}

	parseReference(reference: string): KeychainItem | null {
		const parts = reference.trim().split("/")
		
		if (parts.length < 2) {
			return null
		}
		
		return {
			provider: this.id,
			vault: parts[0],
			item: parts[1],
			field: parts.length > 2 ? parts.slice(2).join("/") : "one-time password"
		}
	}

	async getOTP(item: KeychainItem): Promise<string> {
		return new Promise((resolve, reject) => {
			// Build the secret reference
			const field = item.field || "one-time password"
			const secretRef = `op://${item.vault}/${item.item}/${field}?attribute=otp`
			
			this.logInfo(`Retrieving OTP for ${secretRef}`)
			
			// Execute the 1Password CLI command
			const cmd = `op read "${secretRef}"`
			
			child_process.exec(cmd, {
				timeout: 10000, // 10 second timeout
			}, (error, stdout, stderr) => {
				if (error) {
					this.logError(`Failed to get OTP: ${error.message}`)
					if (stderr) {
						this.logError(`stderr: ${stderr}`)
					}
					reject(new Error(`Failed to retrieve OTP from 1Password: ${error.message}`))
					return
				}
				
				if (stderr) {
					this.logWarning(`stderr: ${stderr}`)
				}
				
				const otpCode = stdout.trim()
				
				// Validate that we got a numeric OTP code
				if (!this.validateOTPFormat(otpCode)) {
					this.logError(`Invalid OTP format received: ${otpCode}`)
					reject(new Error("Invalid OTP format received from 1Password"))
					return
				}
				
				this.logInfo("Successfully retrieved OTP code")
				resolve(otpCode)
			})
		})
	}

	getExamples(): string[] {
		return [
			"Private/GitHub",
			"Work/AWS/one-time password",
			"Personal/VPN Server/authenticator"
		]
	}

	getHelpText(): string {
		return "Format: vault/item or vault/item/field (e.g., Private/GitHub or Private/GitHub/one-time password)"
	}
}
