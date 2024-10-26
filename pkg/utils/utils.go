package utils

import (
	"os"
)

func GetEnv(key string) (string) {
	value, exists := os.LookupEnv(key)
    if !exists {
        return ""
    }
    return value
}