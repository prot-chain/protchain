package function

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func ComputeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Load is a helper function that unfolds a json string into a struct.
// It is important that dest is a pointer to a struct.
func Load(src string, dest interface{}) error {
	return json.Unmarshal([]byte(src), dest)
}
