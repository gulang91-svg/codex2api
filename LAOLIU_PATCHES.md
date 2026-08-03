# Laoliu Codex2API Patches

## 2026-08-11

- Official baseline upgraded through tag `v2.7.4` and main commit `2a57d51`.
- Custom patches were reapplied on top of the official baseline; keep this file and the
   custom commits when applying future upstream updates.
 
- Upstream repository: `james-6-23/codex2api`
- Current upstream tag: `v2.7.4` (`638829d`)
- Upstream commit: `2a57d51`
- Patch history: originally added on `v2.4.8`, then rebased through `v2.5.0`,
  `v2.6.7`, `v2.7.0`, and `v2.7.4`.

### Disable fast mode

In addition to the environment fallback, each API key can now use one of three policies
in the admin UI: disabled, passthrough, or forced. The policy is applied at request
ingress for OpenAI Responses/Chat Completions, Anthropic Messages, and Responses WebSocket
turns. The setting is stored in `APIKeyLimits` as `disable_fast_mode` / `force_fast_mode`.

When `CODEX_DISABLE_FAST_MODE` is enabled (`1`, `true`, `yes`, or `on`) and no per-key
override is present:

- OpenAI-compatible `service_tier: fast/priority` is not forwarded upstream.
- Anthropic `speed: fast` is not mapped to Codex `service_tier: priority`.
- Requested priority is not used as the billing tier fallback.

Requests without the environment switch keep upstream behavior unchanged.

### Images API local 2K/4K upscale

For non-streaming `/v1/images/generations` responses:

- `gpt-image-2-2k` output is locally upscaled to a 2560-pixel long side.
- `gpt-image-2-4k` output is locally upscaled to a 3840-pixel long side.
- Base `gpt-image-2` requests with a 2K/4K `size` receive the same treatment.
- Already-large images are returned without an additional upscale.

The implementation reuses `internal/imageproc` and its bounded global cache.

### Per-account concurrency ceiling

Per-account and account-group base concurrency settings accept values from 1 through
200. Backend validation, frontend input validation, and pool runway estimates use the
same ceiling. The global setting remains governed by the official runtime behavior.

### Verification

```sh
go test ./proxy ./admin ./internal/imageproc
go test ./...
```
