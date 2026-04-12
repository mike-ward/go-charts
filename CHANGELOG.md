# Changelog

## v0.5.3 - 2026-04-12

- Simplify codebase with modern Go 1.26 idioms: `cmp.Or` for defaults,
  `slices.SortFunc`/`cmp.Compare`, builtin `min`/`max`, `slices.Clone`,
  `wg.Go`; extract helpers to deduplicate legend, validation, and
  ring-buffer logic; flatten guards and remove dead code (-285 net lines)

## v0.5.2 - 2026-04-10

- Extract `InteractionCfg` from `BaseCfg`; zoom/pan/range-select/animate-transitions
  fields now live only on XY chart configs
- Expand axis/scale test coverage: table-driven tests for `axis.Linear`,
  `axis.Category`, `scale.Linear`, `scale.Log`

## v0.5.1 - 2026-04-08

- Bump go-gui v0.9.0 → v0.9.1
- Bump go-glyph v1.6.3 → v1.6.4
- Bump golang.org/x/sys v0.42.0 → v0.43.0
