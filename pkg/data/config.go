package data

type Config struct {
    Algorithm string `json:"algorithm"`
	Servers []string `json:"servers"`
	HealtCheck bool `json:"healthCheck"`
} 