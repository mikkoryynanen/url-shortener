package utils

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func GenerateId() string {
	// Define the character set from which you want to generate random characters
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	var result string
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < 32; i++ {
		// Generate a random index within the character set
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			panic(err)
		}
		result += string(charset[index.Int64()])
	}

	return result
}

func LoadEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return os.Getenv(key)
}