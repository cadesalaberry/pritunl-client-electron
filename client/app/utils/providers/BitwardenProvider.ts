/// <reference path="../../References.d.ts"/>
import * as child_process from "child_process"
import { BaseKeychainProvider, KeychainItem } from "../KeychainProviders"

/*
BITWARDEN PROVIDER
==================

Client-side implementation of Bitwarden CLI integration for OTP retrieval.

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

IMPLEMENTATION DETAILS:
- Uses Node.js child_process to execute 'bw' CLI commands
- Checks unlock status via 'bw status' JSON response
- Uses simpler reference format than 1Password (no vault hierarchy)
- Validates OTP format before returning (6-8 digits)

USAGE:
This provider is automatically registered and available when Bitwarden CLI
is installed and the user is logged in with an unlocked session.
*/
export class BitwardenProvider extends BaseKeychainProvider {
	readonly id = "bitwarden"
	readonly name = "Bitwarden"
	readonly description = "Bitwarden CLI integration for OTP retrieval"

	async isAvailable(): Promise<boolean> {
		return new Promise((resolve) => {
			child_process.exec("which bw", (error) => {
				resolve(!error)
			})
		})
	}

	async testAuthentication(): Promise<boolean> {
		return new Promise((resolve) => {
			child_process.exec("bw status", {
				timeout: 5000,
			}, (error, stdout) => {
				if (error) {
					resolve(false)
					return
				}
				
				try {
					const status = JSON.parse(stdout)
					resolve(status.status === "unlocked")
				} catch {
					resolve(false)
				}
			})
		})
	}

	parseReference(reference: string): KeychainItem | null {
		// Bitwarden reference format: "item-name" or "item-id"
		// For this example, we'll use a simple format
		const trimmed = reference.trim()
		
		if (!trimmed) {
			return null
		}
		
		return {
			provider: this.id,
			vault: "", // Bitwarden doesn't use vault concept in the same way
			item: trimmed,
			field: "totp"
		}
	}

	async getOTP(item: KeychainItem): Promise<string> {
		return new Promise((resolve, reject) => {
			this.logInfo(`Retrieving OTP for item: ${item.item}`)
			
			// Execute the Bitwarden CLI command to get TOTP
			const cmd = `bw get totp "${item.item}"`
			
			child_process.exec(cmd, {
				timeout: 10000,
			}, (error, stdout, stderr) => {
				if (error) {
					this.logError(`Failed to get OTP: ${error.message}`)
					if (stderr) {
						this.logError(`stderr: ${stderr}`)
					}
					reject(new Error(`Failed to retrieve OTP from Bitwarden: ${error.message}`))
					return
				}
				
				if (stderr) {
					this.logWarning(`stderr: ${stderr}`)
				}
				
				const otpCode = stdout.trim()
				
				// Validate that we got a numeric OTP code
				if (!this.validateOTPFormat(otpCode)) {
					this.logError(`Invalid OTP format received: ${otpCode}`)
					reject(new Error("Invalid OTP format received from Bitwarden"))
					return
				}
				
				this.logInfo("Successfully retrieved OTP code")
				resolve(otpCode)
			})
		})
	}

	getExamples(): string[] {
		return [
			"GitHub",
			"AWS Account",
			"VPN Server"
		]
	}

	getHelpText(): string {
		return "Format: item-name (e.g., GitHub or AWS Account)"
	}
}
