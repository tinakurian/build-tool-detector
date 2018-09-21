package service

import (
	"github.com/tinakurian/build-tool-detector/controllers/git"
)

// System something
type System struct{}

// ISystem something
type ISystem interface {
	GetGitService()
}

// GetGitService something
func (s System) GetGitService() git.Service {
	return git.Service{}
}
