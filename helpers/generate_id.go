package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

// permette di avere prefissi custom per ogni tipologia (es. products = p-, materials = m- ecc)
func GenerateId(prefix string) string {
	const randomChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	id := prefix

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	id += fmt.Sprintf("%d", timestamp)

	randomValue := rand.Intn(1000)
	id += fmt.Sprintf("%03d", randomValue)

	for i := 0; i < 3; i++ {
		id += string(randomChars[rand.Intn(len(randomChars))])
	}

	return id
}
