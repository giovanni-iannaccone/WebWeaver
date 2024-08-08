package data

// struct that holds configuration info
type Config struct {
	Algorithm  string   `json:"algorithm"`
	Servers    []string `json:"servers"`
	HealtCheck bool     `json:"healthCheck"`
	Logs       bool     `json:"logs"`
}
