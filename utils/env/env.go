package env

import (
	"log"
	"os"
)

func RequireEnv(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.Fatal("cannot find ENV: ", env)
	}

	return str
}

func GetEnv(env string, callback string) string {
	if str, ok := os.LookupEnv(env); ok {
		return str
	}

	return callback
}
