package config

type ServerStore struct {
	URL      *string                `json:"_url"`
	Info     *PublicSystemInfo      `json:"info"`
	Branding *ServerBrandingOptions `json:"config"`
}

type PublicSystemInfo struct {
	LocalAddress           *string
	ServerName             *string
	Version                *string
	ProductName            *string
	OperatingSystem        *string
	Id                     *string
	StartupWizardCompleted bool
}

type ServerBrandingOptions struct {
	LoginDisclaimer     *string
	CustomCss           *string
	SplashscreenEnabled bool
}
