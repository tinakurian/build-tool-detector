package github

import "github.com/davecgh/go-spew/spew"
import "fmt"

// Attributes something
type Attributes struct {
	Owner      string
	Repository string
	Branch     string
}

const sMASTER = "master"
const sTREE = "tree"

// GetAttributes something
func GetAttributes(segments []string, ctxBranch string) Attributes {

	// Default branch that will be used is master
	branch := sMASTER

	// If the ctxBranch is not empty then
	// set the branch name to the branch which
	// the user has specified.
	//
	// If the user has not specified the branch
	// check whether it is passed in through
	// the URL.
	//
	// If the the user has not specified the
	// branch through the URL or through the
	// API parameter 'branch' then default
	// to using master
	if ctxBranch != "" {
		branch = ctxBranch
	} else if len(segments) > 3 {
		if segments[3] == sTREE {
			branch = segments[4]
		}
	}

	spew.Dump(segments)
	fmt.Printf("\n\n\n TINAAAAAAAAAAAAAA \n\n\n %s", "tina")
	return Attributes{Owner: segments[1], Repository: segments[2], Branch: branch}
}
