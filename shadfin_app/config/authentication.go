package config

type AuthenticationStore struct {
	Header   AuthenticationHeader `json:"_header"`
	LoggedIn bool                 `json:"loggedIn"`
}

type AuthenticationHeader struct {
	Client        string  `json:"client"`
	Device        string  `json:"device"`
	DeviceID      string  `json:"deviceID"`
	Version       string  `json:"version"`
	Authorization *string `json:"authorization"`
}
