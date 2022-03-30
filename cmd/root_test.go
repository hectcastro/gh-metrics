package cmd

import (
	"bytes"
	"strings"
	"testing"

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
