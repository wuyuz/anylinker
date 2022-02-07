package helper

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// EqualFold is strings.EqualFold, ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func EqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

func BasicAuth(c *fiber.Ctx) (username, passwords string, ok bool) {
	const prefix = "Basic "
	auth := c.Get("Authorization")
	if auth == ""{
		return
	}

	if len(auth) < len(prefix) || !EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	cx, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(cx)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}