package config

type EnvVar struct {
	Name        string
	Value       any
	Description string
}

func AsMap() map[string]EnvVar {
	return map[string]EnvVar{
		"OMCP_HOST": {
			Name:        "OMCP_HOST",
			Value:       "http://localhost:8080",
			Description: "The HOST of the OMCP server",
		},
	}
}

func Host() string {
	return "http://localhost:8080"
}
