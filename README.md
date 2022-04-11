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
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬─────────────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST REVIEW TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼─────────────────────────────┼─────────────────────────┤
│ 5437 │       2 │       240 │         2 │             9 │ 28m                  │        0 │            2 │ 5h51m             │ --                          │ 5h14m                   │
│ 5436 │       1 │        39 │        18 │             6 │ 9m                   │        0 │            2 │ 1h51m             │ --                          │ 1h40m                   │
│ 5423 │       1 │        41 │         8 │             2 │ 1h7m                 │        0 │            2 │ 67h55m            │ --                          │ 66h44m                  │
│ 5410 │       4 │        10 │        74 │             4 │ 12h45m               │        3 │            4 │ 44h39m            │ 26h47m                      │ 30h41m                  │
│ 5409 │       2 │       113 │         1 │             3 │ 11h58m               │        0 │            3 │ 24h56m            │ 7h29m                       │ 9m                      │
│ 5408 │       1 │         3 │         5 │             1 │ 13h25m               │        0 │            2 │ 18h15m            │ --                          │ 4h52m                   │
│ 5406 │       1 │         1 │         0 │             1 │ 18h5m                │        1 │            4 │ 18h19m            │ --                          │ 14m                     │
│ 5405 │       1 │         2 │         8 │             1 │ 2h17m                │        0 │            5 │ 19h0m             │ 15h47m                      │ 16h38m                  │
│ 5396 │       1 │        37 │         2 │             2 │ 207h33m              │        0 │            3 │ 207h40m           │ --                          │ 6m                      │
│ 5390 │       7 │       183 │         0 │             3 │ 65h41m               │        1 │            4 │ 329h11m           │ 23h56m                      │ 4h29m                   │
│ 5388 │       4 │         8 │         3 │             5 │ 89h0m                │        0 │            4 │ 238h47m           │ 149h47m                     │ 1h39m                   │
│ 5380 │       5 │        54 │         7 │             3 │ --                   │        1 │            6 │ 282h6m            │ 277h6m                      │ --                      │
│ 5378 │       2 │        58 │        10 │             9 │ 184h0m               │        0 │            3 │ 285h40m           │ --                          │ 101h36m                 │
│ 5366 │       6 │       101 │       376 │            12 │ 23h58m               │        1 │            2 │ 234h38m           │ 195h57m                     │ 27m                     │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴─────────────────────────────┴─────────────────────────┘
```

Or, within a more precise window of time:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬─────────────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST REVIEW TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼─────────────────────────────┼─────────────────────────┤
│ 5339 │       4 │         6 │         3 │             1 │ 2m                   │        0 │            3 │ 1h12m             │ 59m                         │ 1h9m                    │
│ 5336 │       1 │         2 │         2 │             2 │ 7m                   │        0 │            1 │ 2h30m             │ --                          │ 2h24m                   │
│ 5327 │       1 │         1 │         1 │             1 │ 41h57m               │        1 │            4 │ 65h44m            │ 23h21m                      │ 23h36m                  │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴─────────────────────────────┴─────────────────────────┘
```

Alternatively, instead of the default table output, output can be generated in CSV format:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22 --csv
PR,Commits,Additions,Deletions,Changed Files,Time to First Review,Comments,Participants,Feature Lead Time,First Review to Last Review,First Approval to Merge
5339,4,6,3,1,2m,0,3,1h12m,59m,1h9m
5336,1,2,2,2,7m,0,1,2h30m,--,2h24m
5327,1,1,1,1,41h57m,1,4,65h44m,23h21m,23h36m
```

## Metric definitions

- **Time to first review**: The duration from when the pull request was created to when the first review against it was completed.
- **Feature lead time**: The duration from when the first commit contained in the pull request was created to when the pull request was merged.
- **First review to last review**: The duration between the first non-author review and the last approving non-author review ([Background](https://github.com/hectcastro/gh-metrics/issues/13)) 
- **First approval to merge**: The duration from when the first approval review is given to when the pull request is merged.

## Influences

Development of this extension was heavily inspired by [jmartin82/mkpis](https://github.com/jmartin82/mkpis).
