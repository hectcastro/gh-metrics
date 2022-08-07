package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nbio/st"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	actual := execute(t, "")
	expected := "Error: required flag(s) \"owner\", \"repo\" not set"

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_Version(t *testing.T) {
	actual := execute(t, "--version")
	expected := fmt.Sprintf("gh-metrics version %s", Version)

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_OnlyOwner(t *testing.T) {
	actual := execute(t, "--owner=cli")
	expected := "Error: required flag(s) \"repo\" not set"

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_OnlyRepo(t *testing.T) {
	actual := execute(t, "--repo=cli")
	expected := "Error: required flag(s) \"owner\" not set"

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
