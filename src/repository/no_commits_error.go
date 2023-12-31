package repository

type NoCommits struct{}

func (n NoCommits) Error() string {
	return "No commits found for repository"
}

func IsErrorTypeNoCommitError(otherError error) bool {
	_, isNotCommitError := otherError.(NoCommits)
	return isNotCommitError
}
