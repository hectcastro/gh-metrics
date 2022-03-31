package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/nbio/st"
)

func execute(args string) string {
	actual := new(bytes.Buffer)

	RootCmd.SetOut(actual)
	RootCmd.SetErr(actual)
	RootCmd.SetArgs(strings.Split(args, " "))
	RootCmd.Execute()

	return actual.String()
}

func Test_RootCmd_NoArgs(t *testing.T) {
	actual := execute("")
	expected := "Error: required flag(s) \"owner\", \"repo\" not set"

	st.Assert(t, strings.Contains(actual, expected), true)
}

func Test_RootCmd_OnlyOwner(t *testing.T) {
	actual := execute("--owner=cli")
	expected := "Error: required flag(s) \"repo\" not set"

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
