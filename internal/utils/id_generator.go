package utils

import "github.com/google/uuid"

// GenerateID generates a new unique string ID
func GenerateID() string {
	return uuid.NewString()
}
