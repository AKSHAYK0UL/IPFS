package helperfunc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GenerateCID(payload string) string {
	godotenv.Load()
	salt := os.Getenv("SALT")
	fmt.Println("SALT", salt)
	data := payload + salt
	sum := sha256.Sum256([]byte(data))

	return hex.EncodeToString(sum[:])
}
