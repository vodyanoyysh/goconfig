package goconfig

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

func LoadConfig(config any, path string) {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(".env file not found, using the system environment variables")
	} else {
		slog.Debug(".env file loaded")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading the file: %v", err))
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		slog.Error(fmt.Sprintf("error unmarshalling the file: %v", err))
	}

	replaceEnvVariables(config)
}

func replaceEnvVariables(config interface{}) {
	val := reflect.ValueOf(config).Elem()
	replaceEnvVariablesRecursive(val)
}

func replaceEnvVariablesRecursive(val reflect.Value) {
	switch val.Kind() {
	case reflect.Ptr:
		replaceEnvVariablesRecursive(val.Elem())
	case reflect.Interface:
		replaceEnvVariablesRecursive(val.Elem())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			replaceEnvVariablesRecursive(val.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			replaceEnvVariablesRecursive(val.Index(i))
		}
	case reflect.String:
		replaceIfEnvVar(val)
	}
}

func replaceIfEnvVar(field reflect.Value) {
	if field.Kind() != reflect.String {
		return
	}
	strVal := field.String()
	if strings.HasPrefix(strVal, "$") {
		envVar := strings.TrimPrefix(strVal, "$")
		envValue := os.Getenv(envVar)
		if envValue != "" {
			field.SetString(envValue)
		} else {
			slog.Warn(fmt.Sprintf("env variable %s not found", envVar))
		}
	}
}
