/*

Package system implements a way
to retrieve different services
for this system.

*/
package system

import (
	"github.com/tinakurian/build-tool-detector/controllers/git"
)

// System struct.
type System struct{}

// ISystem interface.
type ISystem interface {
	GetGitService()
}

// GetGitService gets the git service.
func (s System) GetGitService() git.Service {
	return git.Service{}
}
