# Laoliu Codex2API Patches

## 2026-07-11

- Upstream repository: `james-6-23/codex2api`
- Upstream tag: `v2.5.0`
- Upstream commit: `0fae6cc4fb433c1995b8f739f1c3aa55f8caa8bf`
- Patch history: originally added on `v2.4.8`, then rebased and verified on `v2.5.0`

### Disable fast mode

When `CODEX_DISABLE_FAST_MODE` is enabled (`1`, `true`, `yes`, or `on`):

- OpenAI-compatible `service_tier: fast/priority` is not forwarded upstream.
- Anthropic `speed: fast` is not mapped to Codex `service_tier: priority`.
- Requested priority is not used as the billing tier fallback.

Requests without the environment switch keep upstream behavior unchanged.

### Images API local 2K/4K upscale

For non-streaming `/v1/images/generations` and `/v1/images/edits` responses:

- `gpt-image-2-2k` output is locally upscaled to a 2560-pixel long side.
- `gpt-image-2-4k` output is locally upscaled to a 3840-pixel long side.
- Base `gpt-image-2` requests with a 2K/4K `size` receive the same treatment.
- Already-large images are returned without an additional upscale.
- `2160x3840` and `3840x2160` requests are explicitly covered by regression tests,
  including compressed JPEG responses from the upstream image tool.

The implementation reuses `internal/imageproc` and its bounded global cache.

### Verification

```sh
go test ./proxy ./admin ./internal/imageproc
go test ./...
```
