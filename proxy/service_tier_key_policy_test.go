package proxy

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/codex2api/database"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func TestApplyAPIKeyFastModePolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tc := range []struct {
		name      string
		disabled  bool
		forced    bool
		wantTier  string
		wantSpeed string
	}{
		{name: "allowed", disabled: false, wantTier: "fast", wantSpeed: "fast"},
		{name: "disabled", disabled: true},
		{name: "forced", forced: true, wantTier: "priority", wantSpeed: "fast"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("POST", "/v1/responses", nil)
			c.Set(contextAPIKeyRow, &database.APIKeyRow{
				ID: 1,
				Limits: database.APIKeyLimits{
					DisableFastMode: tc.disabled,
					ForceFastMode:   tc.forced,
				},
			})

			got := applyAPIKeyFastModePolicy(c, []byte(`{"service_tier":"fast","speed":"fast"}`))
			if tier := gjson.GetBytes(got, "service_tier").String(); tier != tc.wantTier {
				t.Fatalf("service_tier = %q, want %q; body=%s", tier, tc.wantTier, got)
			}
			if speed := gjson.GetBytes(got, "speed").String(); speed != tc.wantSpeed {
				t.Fatalf("speed = %q, want %q; body=%s", speed, tc.wantSpeed, got)
			}
		})
	}
}

func TestApplyContextFastModePolicy(t *testing.T) {
	ctx := context.WithValue(context.Background(), fastModePolicyContextKey{}, fastModePolicy{disabled: true})
	got := applyContextFastModePolicy(ctx, []byte(`{"service_tier":"priority"}`))
	if gjson.GetBytes(got, "service_tier").Exists() {
		t.Fatalf("priority tier was not removed: %s", got)
	}
}

func TestApplyAnthropicAPIKeyFastModePolicyForced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/v1/messages", nil)
	c.Set(contextAPIKeyRow, &database.APIKeyRow{
		ID:     1,
		Limits: database.APIKeyLimits{ForceFastMode: true},
	})

	got := applyAnthropicAPIKeyFastModePolicy(c, []byte(`{"model":"claude-test"}`))
	if tier := gjson.GetBytes(got, "service_tier").String(); tier != "priority" {
		t.Fatalf("service_tier = %q, want priority; body=%s", tier, got)
	}
	if speed := gjson.GetBytes(got, "speed").String(); speed != "fast" {
		t.Fatalf("speed = %q, want fast; body=%s", speed, got)
	}
}
