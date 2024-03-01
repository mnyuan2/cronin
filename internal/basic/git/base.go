package git

type Config struct {
	AccessToken string `json:"access_token"`
}

func (m *Config) GetAccessToken() string {
	return m.AccessToken
}
