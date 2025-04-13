package mocks

import "github.com/stretchr/testify/mock"

type MockManagedFiles struct {
	mock.Mock
}

func (self *MockManagedFiles) IsGitAppPresent(repoPath string) bool {
	args := self.Called(repoPath)
	return args.Bool(0)
}

func (self *MockManagedFiles) DirExistsNotEmpty(path string) bool {
	args := self.Called(path)
	return args.Bool(0)
}

func (self *MockManagedFiles) ExecutableExists(path string) bool {
	args := self.Called(path)
	return args.Bool(0)
}

func (self *MockManagedFiles) LinkExists(path string) bool {
	args := self.Called(path)
	return args.Bool(0)
}
