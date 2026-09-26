# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Personal Go playground for implementing algorithms and solving problems, mostly from Donald Knuth's
*The Art of Computer Programming* (TAOCP), primarily Volume 4 (Combinatorial Algorithms:
Boolean basics, permutations, dancing links / exact cover, SAT, CSP). Module path:
`github.com/wallberg/sandbox-go`, Go 1.26.

## Commands

Uses [Task](https://taskfile.dev) (`Taskfile.yaml`) as the primary build tool. It only covers
`./cmd/taocp ./graph ./math ./sgb ./slice ./sortx ./taocp` (i.e. not `golang/`, `cmd/sat_rand`,
`cmd/actor_rpc`, `cmd/test`, `samples/`).

```sh
task build       # go build -v over the core packages
task install      # go install -v over the core packages
task test         # go test over the core packages
task test:long    # long-running tests: --tags=longtests '-run=^TestLong$' -timeout 24h
task doc          # regenerate doc/taocp-7-2-2-2-AlgorithmL.pdf from the .dot source (requires graphviz `dot`)
```

Plain `go` commands work too, and are needed for anything Task doesn't cover (e.g. `golang/`):

```sh
go build ./...
go test ./...
go test ./taocp/...                        # single package
go test ./taocp/... -run TestExercise_7223_3 -v   # single test
```

### Long tests

Tests gated behind the `longtests` build tag (see `taocp/longtests_test.go`) are excluded from
normal `go test` runs and can take hours (some are annotated with actual past run times, e.g.
"Passed in 3.3 hours"). Only run them with `--tags=longtests` when specifically asked to validate
long-running search results — don't run them as part of routine verification.

### Assets

Word lists and shape definitions (`taocp/assets/*.txt`, `*.yaml`) are embedded via
`github.com/gobuffalo/packr`. `sgb.LoadWords()` and friends read these via `packr.NewBox(...)`
with a path relative to the calling package (e.g. `sgb/sgb.go` uses `"../taocp/assets"`,
`taocp/trie.go` and `taocp/polyomino.go` use `"./assets"`). Keep these relative paths in sync if
files move.

## Architecture

- **`taocp/`** — the core library. One file per algorithm/topic, generally named after the TAOCP
  section it implements (each file's top comment cites the Volume/Fascicle and section number,
  e.g. "§7.2.2.1 Dancing Links"). Notable areas:
  - `dancing_links.go` / `dancing_links_xcc.go` (Exact Cover w/ Colors) /
    `dancing_links_mcc.go` (Exact Cover w/ Multiplicities and Colors) — Algorithm X/XCC/MCC
    implementations. Options are the `XCCOptions` struct (Minimax, MinimaxSingle, Exercise83,
    sharp preference heuristic, etc.), stats/progress reporting via `ExactCoverStats`.
  - `sat*.go` — SAT solvers: `sat_algorithm_a.go` (with an `_all` variant enumerating all
    solutions), `sat_algorithm_b.go`, `sat_algorithm_d.go`, `sat_algorithm_l.go`. Shared
    `SatStats`/`SatOptions` types live in `sat.go`.
  - `permutations.go`, `compositions.go`, `partitions.go`, `boolean.go` — combinatorial
    generation from Vol 4A.
  - `polyomino.go` / `polyomino_shapes.go` — polyomino packing (uses `graph/` and YAML shape
    definitions in `taocp/assets/`).
  - `words.go`, `trie.go` — word-based puzzles (double word squares, word stairs, word
    crossings) built on Stanford GraphBase word lists via `sgb/`.
  - `csp_test.go` — Constraint Satisfaction Problems (§7.2.2.3), expressed as SAT and solved via
    the SAT solvers above (no dedicated `csp.go`; CSPs are encoded as SAT instances in tests).
  - Many generator-style functions return Go 1.23+ `iter.Seq`/`iter.Seq2` iterators rather than
    slices or channels — follow this convention for new generators.
- **`sgb/`** — loads Stanford GraphBase data (5-letter word lists, OSPD4 word lists by length)
  from `taocp/assets/`.
- **`graph/`** — thin extensions on top of `github.com/yourbasic/graph` (path/cycle generators,
  helpers), used by polyomino packing.
- **`math/`**, **`slice/`**, **`sortx/`** — small standalone utility packages (numeric helpers,
  slice helpers like cycle-detection/find, sorting extras) used across `taocp/`.
- **`golang/`** — scratch tests exploring Go language features themselves (generics, iterators,
  matrices), not part of the TAOCP work; not included in the Task build/test package list.
- **`cmd/taocp/`** — CLI built with `github.com/jessevdk/go-flags`, exposing the library as
  subcommands (`xcc`, `mcc`, `po` for polyominoes, `wordcross`, `wordstair`, `wr` for word
  rectangles). Each subcommand lives in its own file and self-registers via an `init()` that
  calls `parser.AddCommand(...)`. Commands typically read a YAML problem description from stdin
  (or `-i`) and write YAML/compact results to stdout (or `-o`); mirrors the options structs in
  the underlying `taocp` package (e.g. `xccCommand` mirrors `taocp.XCCOptions`).
  Sample input: `samples/7.2.2.2-sat-as-covering.yaml`.
- **`cmd/sat_rand/`**, **`cmd/actor_rpc/`**, **`cmd/test/`** — standalone experiment mains, not
  wired into the `taocp` CLI or the Taskfile package list.

## Git commits

Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `test:`,
`docs:`, `refactor:`, etc. — matches existing history, `git log --oneline`). Include Claude
attribution in the commit message body and add the model:

```
Co-Authored-By: Claude <noreply@anthropic.com>
```

## Conventions

- Code comments frequently reference the specific TAOCP volume/fascicle, section, exercise
  number, and sometimes page number the implementation follows (e.g. "Exercise 29.", "Algorithm
  L, p. 319"). Preserve/add these citations when touching algorithm code — they're the primary
  way to trace an implementation back to its source.
- A few non-TAOCP helper functions are explicitly marked as such in their doc comment (e.g.
  `StrictPartitions` in `taocp/partitions.go` says "Not a TAOCP implementation"). Keep that
  distinction visible when adding helpers that aren't from the book.
