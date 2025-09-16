package cmd

var (
	mode                     string
	password                 string
	passwordPrompt           bool
	disableKeychain          bool
	enabledKeychainProviders string // Comma-separated list of provider IDs
	jsonFormat               bool
	jsonFormated             bool
)

func init() {
	StartCmd.Flags().StringVarP(
		&mode,
		"mode",
		"m",
		"",
		"VPN mode (ovpn, wg)",
	)
	StartCmd.Flags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"VPN password",
	)
	StartCmd.Flags().BoolVarP(
		&passwordPrompt,
		"password-read",
		"r",
		false,
		"Prompt for VPN password",
	)
	StartCmd.Flags().BoolVar(
		&disableKeychain,
		"disable-keychain",
		false,
		"Disable ALL keychain provider integrations",
	)
	StartCmd.Flags().StringVar(
		&enabledKeychainProviders,
		"keychain-providers",
		"",
		"Comma-separated list of keychain providers to enable (e.g., 1password,bitwarden)",
	)

	ListCmd.Flags().BoolVarP(
		&jsonFormat,
		"json",
		"j",
		false,
		"Format output in JSON",
	)

	ListCmd.Flags().BoolVarP(
		&jsonFormated,
		"json-formatted",
		"f",
		false,
		"Format output in indented JSON",
	)
}
