package proxy

import (
	"os"
	"strings"
)

func fastModeDisabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("CODEX_DISABLE_FAST_MODE"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
