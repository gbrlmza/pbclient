package pbclient

import (
	"bytes"
	"regexp"
)

const (
	pbTimeRegex = `"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(.\d{1,3})?Z"`
)

// parseFromPB converts dates in pocketbase format(2006-01-02 15:04:05.000Z) to RFC3339.
func parseFromPB(b []byte) []byte {
	re := regexp.MustCompile(pbTimeRegex)
	return re.ReplaceAllFunc(b, func(b []byte) []byte {
		return bytes.Replace(b, []byte(" "), []byte("T"), 1)
	})
}

// parseFromPB converts dates in RFC3339 format to pocketbase format(2006-01-02 15:04:05.000Z).
func parseToPB(b []byte) []byte {
	// TODO:
	return b
}
