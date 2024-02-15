package types

func NewProviderConfig(iss string, jwkEndpoint string) ProviderConfig {
	return ProviderConfig{
		Iss:            iss,
		JwkEndpoint:    jwkEndpoint,
		FetchIntervals: 60,
	}
}

func GetConfig(provider OidcProvider) ProviderConfig {
	switch provider {
	case Google:
		return NewProviderConfig("https://accounts.google.com", "https://www.googleapis.com/oauth2/v3/certs")
	default:
		panic("unexpected provider")
	}
}
