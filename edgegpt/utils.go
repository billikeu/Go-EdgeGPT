package edgegpt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

var DELIMITER string = "\x1e"

// Appends special character to end of message to identify end of message
func appendIdentifier(msg map[string]interface{}) (string, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(b) + DELIMITER, nil
}

// Returns random hex string
func GetRandomHex(n int, allowedChars ...[]rune) string {
	var letters []rune
	if len(allowedChars) == 0 {
		letters = []rune("0123456789abcdef")
	} else {
		letters = allowedChars[0]
	}
	b := make([]rune, n)
	for i := range b {
		rand.Seed(time.Now().UTC().UnixNano() + int64(i<<20))
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Generate random IP between range 13.104.0.0/14
func GetRandomIp() string {
	ip := fmt.Sprintf("13.%d.%d.%d", 104+rand.Intn(3), rand.Intn(255), rand.Intn(255))
	return ip
}

func GetUuidV4() string {
	return uuid.NewV4().String()
}
