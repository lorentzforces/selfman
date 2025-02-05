package mocks

import "github.com/stretchr/testify/mock"

type MockManagedFiles struct {
	mock.Mock
}

func (self MockManagedFiles) isGitAppPresent(appPath string) bool {
	args := self.Called(appPath)
	return args.Bool(0)
}
