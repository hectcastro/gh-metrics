# gh-metrics

A [`gh`](https://cli.github.com/) extension that provides summary pull request metrics.

- [Usage](#usage)
- [Metric definitions](#metric-definitions)
- [Influences](#influences)

## Usage

To install the extension use:

```console
$ gh extension install hectcastro/gh-metrics
```

Once installed, you can summarize all pull requests for the `cli/cli` repository over the last 10 days:

```console
$ gh metrics --owner cli --repo cli
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ LAST REVIEW TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┤
│ 5339 │       4 │         6 │         3 │             1 │ 2m0s                 │        0 │            3 │ 1h12m0s           │ 1h9m0s               │
│ 5336 │       1 │         2 │         2 │             2 │ 7m0s                 │        0 │            1 │ 2h30m0s           │ 2h24m0s              │
│ 5327 │       1 │         1 │         1 │             1 │ 41h57m0s             │        1 │            4 │ 65h44m0s          │ 14m0s                │
│ 5323 │       1 │         3 │         2 │             2 │ 11m0s                │        0 │            1 │ 15m0s             │ 4m0s                 │
│ 5319 │      16 │       318 │        39 │            12 │ 14h23m0s             │        1 │            4 │ 18h9m0s           │ 2h43m0s              │
│ 5298 │       1 │         1 │         1 │             1 │ 16m0s                │        0 │            2 │ 72h29m0s          │ 72h13m0s             │
│ 5297 │       1 │         1 │         0 │             1 │ 10h51m0s             │        0 │            4 │ 83h8m0s           │ 72h11m0s             │
│ 5296 │       6 │        46 │         1 │             2 │ 2h7m0s               │        1 │            5 │ 84h35m0s          │ 82h7m0s              │
│ 5279 │       3 │        31 │         4 │             6 │ 252h35m0s            │        2 │            3 │ 252h37m0s         │ 1m0s                 │
│ 5276 │       1 │         2 │         1 │             1 │ 15h41m0s             │        0 │            4 │ 284h19m0s         │ 29h2m0s              │
│ 5275 │       7 │       179 │        37 │             4 │ 169h8m0s             │        3 │            4 │ 787h15m0s         │ 23h35m0s             │
│ 5270 │      27 │       520 │        72 │            11 │ 10h49m0s             │        1 │            6 │ 719h41m0s         │ 336h17m0s            │
│ 5251 │       4 │       133 │        24 │             3 │ 180h49m0s            │        0 │            4 │ 423h13m0s         │ 2h17m0s              │
│ 5160 │       2 │       215 │        65 │             4 │ 987h48m0s            │        2 │            4 │ 31h33m0s          │ 0s                   │
│ 5134 │       9 │       445 │        66 │             6 │ 854h10m0s            │        0 │            3 │ 1043h38m0s        │ 0s                   │
│ 4895 │      15 │       497 │        32 │             3 │ 236h37m0s            │        7 │            4 │ 2211h51m0s        │ 1h41m0s              │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┘
```

Or, within a more precise window of time:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ LAST REVIEW TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┤
│ 5339 │       4 │         6 │         3 │             1 │ 2m0s                 │        0 │            3 │ 1h12m0s           │ 1h9m0s               │
│ 5336 │       1 │         2 │         2 │             2 │ 7m0s                 │        0 │            1 │ 2h30m0s           │ 2h24m0s              │
│ 5327 │       1 │         1 │         1 │             1 │ 41h57m0s             │        1 │            4 │ 65h44m0s          │ 14m0s                │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┘
```

Alternatively, instead of the default table output, output can be generated in CSV format:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22 --csv
PR,Commits,Additions,Deletions,Changed Files,Time to First Review,Comments,Participants,Feature Lead Time,Last Review to Merge
5339,4,6,3,1,2m0s,0,3,1h12m0s,1h9m0s
5336,1,2,2,2,7m0s,0,1,2h30m0s,2h24m0s
5327,1,1,1,1,41h57m0s,1,4,65h44m0s,14m0s
```

## Metric definitions

- **Time to first review**: The duration from when the pull request was created to when the first review against it was completed.
- **Feature lead time**: The duration from when the first commit contained in the pull request was created to when the pull request was merged.
- **First approval to merge**: The duration from when the first approval review is given to when the pull request is merged.

## Influences

Development of this extension was heavily inspired by [jmartin82/mkpis](https://github.com/jmartin82/mkpis).
