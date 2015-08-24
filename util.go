package martiniframework

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var (
	// IsRunning specifies whether or not the application is running.
	// It governs the SleepFor function.
	IsRunning = true
)

// Md5hash is a convenience function creating and returing an MD5 hash
// for a specified string.
func Md5hash(orig string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(orig)))
}

func SleepFor(sec int64) {
	for i := 0; i < int(sec); i++ {
		if !IsRunning {
			return
		}
		time.Sleep(time.Second)
	}
}

// JsonEncode is a convenience function wrapping the serialization of an
// interface{} to its JSON byte equivalent. If an error occurs, "false"
// is returned.
func JsonEncode(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Print(err.Error())
		return []byte("false")
	}
	return b
}
