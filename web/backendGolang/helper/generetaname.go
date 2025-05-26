package helper

import (
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateRandomFileName generates a random file name with the format YYYYMMDD_randomString_HHMMSS.ext
func GenerateRandomFileName(ext string) string {
	now := time.Now()
	
	datePrefix := now.Format("20060102")

	randomString := generateRandomString(5)

	timeSuffix := now.Format("150405")

	return fmt.Sprintf("%s_%s_%s%s", datePrefix, randomString, timeSuffix, ext)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}