package utils

import (
	"errors"
	"net/url"
	"os"
	"strconv"
)

func GetString(envar string) string {
	return os.Getenv(envar)
}

func GetInt(envar string) int {
	val, err := strconv.Atoi(os.Getenv(envar))

	if err != nil {
		panic(errors.New("Could not parse string to integer"))
	}

	return val
}

func GetURL(envar string) *url.URL {
	val, err := url.Parse(os.Getenv(envar))

	if err != nil {
		panic(errors.New("Could not parse URL"))
	}

	return val
}
