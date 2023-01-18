package cmd

import (
	"errors"
	"strings"

	"github.com/cli/go-gh/pkg/auth"
)

// GHRepo represents a GitHub repository.
type GHRepo struct {
	Owner string
	Name  string
	Host  string
}

// newGHRepo returns a GHRepo configured via the '/'-delimited string it's passed.
func newGHRepo(name string) (*GHRepo, error) {
	defaultHost, _ := auth.DefaultHost()
	nameParts := strings.Split(name, "/")

	switch len(nameParts) {
	case 2:
		return &GHRepo{
			Owner: nameParts[0],
			Name:  nameParts[1],
			Host:  defaultHost,
		}, nil
	case 3:
		return &GHRepo{
			Owner: nameParts[1],
			Name:  nameParts[2],
			Host:  nameParts[0],
		}, nil
	default:
		return nil, errors.New("invalid repository name")
	}
}
