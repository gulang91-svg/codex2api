# Laoliu Codex2API Patches

## 2026-07-29

- Upstream tag: `v2.6.7`
- Upstream commit: `4aa58cb6d3c7fe6fb4eece7d870ebdc44d0def96`
- Integration: merged the official release into the Laoliu branch while
  retaining the local image, fast-mode, concurrency, proxy, logging, pricing,
  and admin customizations.
- Upstream additions retained: weak-network WebSocket mode, lower streaming
  memory retention, provider-qualified model names, structured batch operation
  results, proxy pagination, and the `0.144.1` minimum Codex CLI baseline.
- Database migration: add `codex_ws_weak_network_mode` with a default of
  `false`; existing deployments remain behaviorally unchanged until enabled.

Verification:

```sh
go test ./...
cd frontend
npm run typecheck
npm test
npm run build
```

## 2026-07-28

- Upstream tag: `v2.6.6`
- Upstream commit: `76b2f38d0e1ed379017b25eeeaafa1a0dc407df0`
- Integration: merged the official release into the Laoliu branch without
  replacing the local image, fast-mode, concurrency, proxy, logging, pricing,
  and admin customizations.
- Compatibility fix: preserve both the upstream no-affinity account-group
  routing fields and the local per-key fast-mode policy fields.
- Reliability fix: make the legacy Responses-success recovery path robust when
  the platform clock returns the same timestamp as the preceding rate limit.
- Tooling fix: make the pool-runway TypeScript import work in both the Node test
  runner and the Vite production build.

Verification:

```sh
go test ./...
cd frontend
npm run typecheck
npm test
npm run build
```

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

### Per-account concurrency ceiling

The global, per-account, batch-account, and account-group base concurrency
settings accept values from 1 through 200. Backend validation, frontend input
validation, and pool runway estimates use the same ceiling.

### Verification

```sh
go test ./proxy ./admin ./internal/imageproc
go test ./...
```
