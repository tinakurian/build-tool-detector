/*

Package github implements a way to extract
and construct a request to github in order
to retrieve a pom file. If the pom file is
not present, we assume the project is not
build using maven.

*/
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

// GetAttributes will use the path segments and
// query params to populate the Attributes
// struct. The attributes struct will be used
// to make a request to github to determine
// the build tool type.
func GetAttributes(segments []string, ctxBranch *string) Attributes {

	// Default branch that will be used if a branch
	// is not passed in though the optional 'branch'
	// query parameter and is not part of the url.
	branch := sMASTER

	// If the query parameter field 'branch' is not
	// empty then set the branch name to the query
	// parameter value.
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
