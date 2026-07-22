# Laoliu Codex2API Patches

## 2026-07-10

- Upstream repository: `james-6-23/codex2api`
- Upstream tag: `v2.4.8`
- Upstream commit: `a4a9268bf8c839046072a24c728617c281367c63`

### Disable fast mode

When `CODEX_DISABLE_FAST_MODE` is enabled (`1`, `true`, `yes`, or `on`):

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

The global, per-account, batch-account, and account-group base concurrency
settings accept values from 1 through 200. Backend validation, frontend input
validation, and pool runway estimates use the same ceiling.

### Verification

```sh
go test ./proxy ./admin ./internal/imageproc
go test ./...
```
