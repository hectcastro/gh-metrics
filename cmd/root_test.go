package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/nbio/st"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/h2non/gock.v1"
)

func ResetSubCommandFlagValues(t *testing.T, root *cobra.Command) {
	t.Helper()

	root.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			f.Value.Set(f.DefValue)
			f.Changed = false
		}
	})
}

func execute(t *testing.T, args string) string {
	t.Helper()

	actual := new(bytes.Buffer)

	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	ResetSubCommandFlagValues(t, RootCmd)
	RootCmd.SetArgs(strings.Split(args, " "))
	RootCmd.Execute()

	return actual.String()
}

func Test_RootCmd_NoArgs(t *testing.T) {
	defer gock.Off()

	currentRepo, err := gh.CurrentRepository()
	if err != nil {
		t.Error(err)
	}

	gock.New("https://api.github.com/graphql").
		Post("/").
		MatchType("json").
		AddMatcher(gqlSearchQueryMatcher(currentRepo.Owner(), currentRepo.Name(), defaultStart, defaultEnd)).
		Reply(200).
		BodyString(ResponseJSON)

	actual := execute(t, "")
	expected := `┌──────┬─────────┬───────────┬───────────┬───────────────┬──────────────────────┬──────────┬──────────────┬───────────────────┬──────────────────────┬─────────────────────────┐
│   PR │ COMMITS │ ADDITIONS │ DELETIONS │ CHANGED FILES │ TIME TO FIRST REVIEW │ COMMENTS │ PARTICIPANTS │ FEATURE LEAD TIME │ FIRST TO LAST REVIEW │ FIRST APPROVAL TO MERGE │
├──────┼─────────┼───────────┼───────────┼───────────────┼──────────────────────┼──────────┼──────────────┼───────────────────┼──────────────────────┼─────────────────────────┤
│ 5339 │       1 │         6 │         3 │             1 │ 155h26m              │        0 │            3 │ 1h12m             │ 24h0m                │ 22h51m                  │
│ 5340 │       1 │        12 │         6 │             2 │ 155h26m              │        0 │            3 │ 1h12m             │ 24h0m                │ 22h51m                  │
└──────┴─────────┴───────────┴───────────┴───────────────┴──────────────────────┴──────────┴──────────────┴───────────────────┴──────────────────────┴─────────────────────────┘`

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_Version(t *testing.T) {
	actual := execute(t, "--version")
	expected := fmt.Sprintf("gh-metrics version %s", Version)

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_OnlyRepo(t *testing.T) {
	actual := execute(t, "--repo=cli")
	expected := "invalid repository name"

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_WorkdayOnlyWeekdays(t *testing.T) {
	friday := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	saturday := time.Date(2022, 4, 2, 0, 0, 0, 0, time.UTC)
	sunday := time.Date(2022, 4, 3, 0, 0, 0, 0, time.UTC)

	st.Assert(t, WorkdayOnlyWeekdays(friday), true)
	st.Assert(t, WorkdayOnlyWeekdays(saturday), false)
	st.Assert(t, WorkdayOnlyWeekdays(sunday), false)
}

func Test_WorkdayAllDays(t *testing.T) {
	friday := time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC)
	saturday := time.Date(2022, 4, 2, 0, 0, 0, 0, time.UTC)
	sunday := time.Date(2022, 4, 3, 0, 0, 0, 0, time.UTC)

	st.Assert(t, WorkdayOnlyWeekdays(friday), true)
	st.Assert(t, WorkdayAllDays(saturday), true)
	st.Assert(t, WorkdayAllDays(sunday), true)
}
