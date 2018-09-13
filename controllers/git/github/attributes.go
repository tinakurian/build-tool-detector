package github

// Attributes something
type Attributes struct {
	Owner      string
	Repository string
	Branch     string
}

const (
	sMASTER = "master"
	sTREE   = "tree"
)

// GetAttributes something
func GetAttributes(segments []string, ctxBranch *string) Attributes {

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url
	branch := sMASTER

	// If the query parameter field 'branch' is not
	// empty then set the branch name to the query
	// parameter value
	if ctxBranch != nil {
		branch = *ctxBranch
	} else if len(segments) > 3 {
		// If the user has not specified the branch
		// check whether it is passed in through
		// the URL.
		if segments[3] == sTREE {
			branch = segments[4]
		}
	}

	return Attributes{Owner: segments[1], Repository: segments[2], Branch: branch}
}
