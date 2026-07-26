package proxy

import (
	"context"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type fastModePolicyContextKey struct{}

type fastModePolicy struct {
	disabled bool
	forced   bool
}

func normalizeFastModePolicy(disabled, forced bool) fastModePolicy {
	return fastModePolicy{disabled: disabled, forced: forced && !disabled}
}

func attachAPIKeyFastModePolicy(c *gin.Context, disabled, forced bool) {
	if c == nil || c.Request == nil {
		return
	}
	ctx := context.WithValue(c.Request.Context(), fastModePolicyContextKey{}, normalizeFastModePolicy(disabled, forced))
	c.Request = c.Request.WithContext(ctx)
}

func fastModePolicyFromContext(ctx context.Context) fastModePolicy {
	if ctx != nil {
		if policy, ok := ctx.Value(fastModePolicyContextKey{}).(fastModePolicy); ok {
			return policy
		}
	}
	return normalizeFastModePolicy(fastModeDisabled(), false)
}

func fastModeDisabledFromContext(ctx context.Context) bool {
	return fastModePolicyFromContext(ctx).disabled
}

func fastModePolicyForAPIKey(c *gin.Context) fastModePolicy {
	if row := apiKeyRowFromContext(c); row != nil {
		return normalizeFastModePolicy(row.Limits.DisableFastMode, row.Limits.ForceFastMode)
	}
	if c != nil && c.Request != nil {
		return fastModePolicyFromContext(c.Request.Context())
	}
	return normalizeFastModePolicy(fastModeDisabled(), false)
}

func stripFastModeRequestFields(body []byte) []byte {
	for _, path := range []string{"service_tier", "serviceTier"} {
		tier := strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, path).String()))
		if tier == "fast" || tier == "priority" {
			body, _ = sjson.DeleteBytes(body, path)
		}
	}
	if speed := strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "speed").String())); speed == "fast" {
		body, _ = sjson.DeleteBytes(body, "speed")
	}
	return body
}

func forceFastModeRequestFields(body []byte) []byte {
	body, _ = sjson.DeleteBytes(body, "serviceTier")
	body, _ = sjson.SetBytes(body, "service_tier", "priority")
	return body
}

func applyAPIKeyFastModePolicy(c *gin.Context, body []byte) []byte {
	policy := fastModePolicyForAPIKey(c)
	if policy.disabled {
		return stripFastModeRequestFields(body)
	}
	if policy.forced {
		return forceFastModeRequestFields(body)
	}
	return body
}

func applyAnthropicAPIKeyFastModePolicy(c *gin.Context, body []byte) []byte {
	policy := fastModePolicyForAPIKey(c)
	if policy.disabled {
		return stripFastModeRequestFields(body)
	}
	if !policy.forced {
		return body
	}
	body = forceFastModeRequestFields(body)
	body, _ = sjson.SetBytes(body, "speed", "fast")
	return body
}

func applyContextFastModePolicy(ctx context.Context, body []byte) []byte {
	policy := fastModePolicyFromContext(ctx)
	if policy.disabled {
		return stripFastModeRequestFields(body)
	}
	if policy.forced {
		return forceFastModeRequestFields(body)
	}
	return body
}

func fastModeDisabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("CODEX_DISABLE_FAST_MODE"))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
