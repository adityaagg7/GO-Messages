package utils

import (
	"fmt"
	namesgenerator "github.com/dillonstreator/go-unique-name-generator"
	"math/rand"
)

// GenerateRoomName generates a random, unique room name by combining a base name with a random three-digit number.
func GenerateRoomName() string {
	base := namesgenerator.NewUniqueNameGenerator().Generate()
	number := rand.Intn(900) + 100
	return fmt.Sprintf("%s-%d", base, number)
}
