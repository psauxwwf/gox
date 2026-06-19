# Project Notes

- CLI is built with `cobra` and launched with `fang`.
- Use subcommands like `gox setup`, `gox remove`, `gox save` instead of legacy top-level flags.
- Use `log/slog` for all project logging.
- Read and write YAML config with `gopkg.in/yaml.v3` marshal/unmarshal directly; do not reintroduce `cleanenv`.
- Run `garble` via `go tool garble`; do not depend on a separately installed `garble` binary.
- Keep current behavior where missing config falls back to defaults.
- Keep current `setup` behavior where all installation steps are attempted and collected even if one of them fails.
