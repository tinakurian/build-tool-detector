package types

import (
	"context"
)

// RepositoryService holds information about
// the repository
type RepositoryService interface {
	Owner() string
	Repository() string
	Branch() string
	DetectBuildTool(ctx context.Context) (*string, error)
}
