package main

type DapnetNewsSettings struct {
	CallsignNames         []string `yaml:"callsignNames"`
	TransmitterGroupNames []string `yaml:"transmitterGroupNames"`

	NewsEndpoint  string  `yaml:"news-endpoint"`
	TTL           float64 `yaml:"cache-ttl"`
	DeliveryDelay int     `yaml:"delivery-delay"`
	CleanInterval string  `yaml:"clean-interval"`
	CheckInterval string  `yaml:"check-interval"`

	DapnetUsername string `yaml:"dapnet-username"`
	DapnetPassword string `yaml:"dapnet-password"`
	DapnetCallsign string `yaml:"dapnet-callsign"`
}
