package cmd

import (
	"testing"

	"github.com/nbio/st"
)

func TestNewGHRepo(t *testing.T) {
	tests := []struct {
		name         string
		wantName     string
		wantOwner    string
		wantHost     string
		wantFullName string
		wantErr      bool
		errMsg       string
	}{
		{
			name:      "foo/bar",
			wantOwner: "foo",
			wantName:  "bar",
			wantHost:  "github.com",
			wantErr:   false,
		}, {
			name:      "other-github.com/foo/bar",
			wantOwner: "foo",
			wantName:  "bar",
			wantHost:  "other-github.com",
			wantErr:   false,
		}, {
			name:    "bar",
			wantErr: true,
			errMsg:  "invalid repository name",
		}, {
			name:    "",
			wantErr: true,
			errMsg:  "invalid repository name",
		}}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ghr, err := newGHRepo(tt.name)

			if tt.wantErr {
				st.Assert(t, err.Error(), tt.errMsg)

			} else {
				st.Assert(t, err, nil)
				st.Assert(t, ghr.Owner, tt.wantOwner)
				st.Assert(t, ghr.Name, tt.wantName)
				st.Assert(t, ghr.Host, tt.wantHost)
			}
		})
	}
}
