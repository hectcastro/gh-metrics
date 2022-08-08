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
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┼─────────────────────────┤
│ 6029 │       1 │         3 │         2 │             1 │ 26m                  │        1 │            4 │ 1h9m              │ 40m                  │ 40m                     │
│ 6019 │       2 │         8 │         0 │             1 │ 19h13m               │        1 │            4 │ 23h15m            │ --                   │ 3h58m                   │
│ 6008 │       1 │         1 │        12 │             2 │ 12h19m               │        1 │            4 │ 185h5m            │ 167h54m              │ 4h51m                   │
│ 6004 │       1 │        18 │         0 │             1 │ 149h59m              │        3 │            5 │ 208h47m           │ 6h7m                 │ 58h48m                  │
│ 5974 │       1 │         1 │         1 │             1 │ 130h54m              │        1 │            5 │ 262h58m           │ 6h55m                │ 178h34m                 │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┴─────────────────────────┘
```

Or, within a more precise window of time:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┼─────────────────────────┤
│ 5339 │       4 │         6 │         3 │             1 │ 2m                   │        0 │            3 │ 1h12m             │ 59m                  │ 1h9m                    │
│ 5336 │       1 │         2 │         2 │             2 │ 7m                   │        0 │            1 │ 2h30m             │ --                   │ 2h24m                   │
│ 5327 │       1 │         1 │         1 │             1 │ 41h57m               │        1 │            4 │ 65h44m            │ 23h21m               │ 23h36m                  │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┴─────────────────────────┘
```

Or, with an additional query filter:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22 --query "author:josebalius"
┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┼─────────────────────────┤
│ 5339 │       4 │         6 │         3 │             1 │ 2m                   │        0 │            3 │ 1h12m             │ 59m                  │ 1h9m                    │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┴─────────────────────────┘
```

Alternatively, instead of the default table output, output can be generated in CSV format:

```console
$ gh metrics --owner cli --repo cli --start 2022-03-21 --end 2022-03-22 --csv
PR,Commits,Additions,Deletions,Changed Files,Time to First Review,Comments,Participants,Feature Lead Time,First to Last Review,First Approval to Merge
5339,4,6,3,1,00:02,0,3,01:12,00:59,01:09
5336,1,2,2,2,00:07,0,1,02:30,00:00,02:24
5327,1,1,1,1,41:57,1,4,65:44,23:21,23:36
```

## Metric definitions

- **Time to first review**: The duration from when the pull request was created or marked *Ready for review* to when the first review against it was completed.
- **Feature lead time**: The duration from when the first commit contained in the pull request was created to when the pull request was merged.
- **First review to last review**: The duration between the first non-author review and the last approving non-author review ([Background](https://github.com/hectcastro/gh-metrics/issues/13)) 
- **First approval to merge**: The duration from when the first approval review is given to when the pull request is merged.

## Influences

Development of this extension was heavily inspired by [jmartin82/mkpis](https://github.com/jmartin82/mkpis).
