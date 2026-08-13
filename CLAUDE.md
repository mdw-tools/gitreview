# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

- `make test` — format, vet, and run tests with the race detector
- `make fmt` — run `go mod tidy` and `go fmt ./...`
- `make install` — install the binary with version from `git describe`
- `make docs` — regenerate README.md from the built-in help text in `config.go`
- `go test -run TestName ./...` — run a single test

## Overview

gitreview is a Go CLI tool (zero external dependencies) that scans a directory
tree for git repositories, concurrently analyzes their state (dirty, ahead,
behind, fetch results), and presents the repositories needing attention for
interactive review via an external git GUI (default: smerge).

## Architecture

**Startup flow:** `main.go` → `ReadConfig()` parses flags/args →
`NewGitReviewer(config)` → `GitAnalyzeAll()` → `ReviewAll()`.

**Concurrency model:** `analyzer.go` manages a pool of 16 worker goroutines
(`worker.go`). Workers receive repo paths from an input channel, run all git
commands (`git.go`), and send `GitReport` results to an output channel.
`merge()` collects results via fan-in.

**Key types:**
- `Config` (config.go) — CLI flags and the resolved root directory
- `GitReport` (git.go) — per-repo analysis results; methods run git commands
  and build the progress indicator `[!MABFS]`
- `GitReviewer` (review.go) — orchestrates analysis, categorizes repos into maps
  (erred/messy/ahead/behind/fetched/skipped), and drives interactive review
- `Analyzer` / `Worker` (analyzer.go, worker.go) — concurrent worker pool

**I/O helpers** (io.go): `collectGitRepositories()` walks the directory tree
(skipping `.git` and not descending into found repos; symlinks pointing at git
repos are followed but never descended into, deduplicated by resolved path),
`execute()` shells out to git, `prompt()` reads stdin for interactive review.

## Root Directory Resolution

When no path argument is supplied, the scan root is resolved in this order:

1. the first non-flag command-line argument, else
2. `$CODEPATH/src` (when `CODEPATH` is set), else
3. the current working directory.

## Per-repo Git Config

- `review.skip true` — skip the repo entirely
- `review.branch <name>` — override default-branch detection (otherwise
  main/master is auto-detected)
