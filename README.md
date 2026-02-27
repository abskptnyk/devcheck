# devcheck

Check if your local dev environment is ready to run a project.

```
$ devcheck

Detected: Go, Docker

✅  go installed
✅  Docker daemon running
❌  .env missing keys: DB_PASSWORD, API_KEY

────────────────────────────────────────
2 passed  0 warnings  1 failed
Run devcheck --fix for suggestions
```

Run it in any project directory. devcheck looks at your project files, figures out what stack you're using, and checks that everything is actually running and configured before you waste time debugging.

---

## Install

**One-liner:**
```bash
curl -fsSL https://raw.githubusercontent.com/vidya381/devcheck/main/scripts/install.sh | bash
```

**Go install:**
```bash
go install github.com/vidya381/devcheck/cmd/devcheck@latest
```

Or grab a binary from [releases](https://github.com/vidya381/devcheck/releases) — linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64.

---

## Usage

```bash
devcheck              # run all checks
devcheck --fix        # show a suggested fix for each failure
devcheck --json       # output results as JSON
devcheck --ci         # exit code 1 on any failure (for CI pipelines)
devcheck --verbose    # show skipped checks too
```

---

## What it checks

devcheck detects your stack from project files and runs the relevant checks automatically — no config needed.

| Detected from | Checks |
|---------------|--------|
| `go.mod` | go binary, go version |
| `package.json` | node, npm, node version |
| `requirements.txt` / `pyproject.toml` | python3, pip |
| `pom.xml` / `build.gradle` | java, mvn or gradle |
| `Dockerfile` / `docker-compose.yml` | docker binary, daemon running |
| `DATABASE_URL` in env | PostgreSQL reachable |
| `REDIS_URL` in env | Redis reachable |
| `.env.example` present | all keys exist in .env |

---

## CI usage

```yaml
- name: Check dev environment
  run: devcheck --ci
```

Exits with code 1 if any check fails, so your pipeline stops early.

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Adding a new check is straightforward — each check is a small self-contained struct. Issues labeled `good first issue` are a good place to start.
