package data

// struct that holds configuration info
type Config struct {
	Algorithm   string   `json:"algorithm"`
	Servers     []string `json:"servers"`
	HealthCheck bool     `json:"healthCheck"`
	Logs        string   `json:"logs"`
	Prohibited  []string `json:"prohibited"`
}
