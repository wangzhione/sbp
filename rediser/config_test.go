package rediser

import (
	"testing"
)

func TestParseRedisCLI(t *testing.T) {
	rediscommand := "redis-cli -h 1.0.0.1 -p 6489 -a mypassword"

	options, err := ParseRedisCLI(rediscommand)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(options)
}
