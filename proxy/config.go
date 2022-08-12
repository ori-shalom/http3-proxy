package proxy

import (
	"fmt"
	"os"
)

const (
	portVarName       = "PORT"
	targetHostVarName = "TARGET_HOST"
)

var variables = map[string]envVarOptions{
	targetHostVarName: {required: true},
	portVarName:       {defaultValue: "8080"},
}

type Config struct {
	Port       string
	TargetHost string
}

type envVarOptions struct {
	defaultValue string
	required     bool
}

func LoadConfig() (Config, error) {
	vars, err := loadEnvironmentVariables(variables)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Port:       vars[portVarName],
		TargetHost: vars[targetHostVarName],
	}, nil
}

// loadEnvironmentVariables load environment variables, validate they are not empty and return them as map.
// if one of the variables is missing or empty return an error.
func loadEnvironmentVariables(envVars map[string]envVarOptions) (map[string]string, error) {
	values := make(map[string]string, len(envVars))
	for name, options := range envVars {
		value, found := os.LookupEnv(name)
		if options.required && !found {
			return nil, fmt.Errorf("missing required environment variable %q", name)
		}
		if !found {
			value = options.defaultValue
		}
		values[name] = value
	}
	return values, nil
}
