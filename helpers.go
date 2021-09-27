package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func createHash(title string, date *time.Time) string {
	str := fmt.Sprintf("%s:%s", title, date)
	str = strings.ReplaceAll(str, " ", "")
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}
