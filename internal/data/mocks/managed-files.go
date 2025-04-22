package mocks

import (
	"github.com/lorentzforces/selfman/internal/data"
	"github.com/stretchr/testify/mock"
)

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

func (self *MockManagedFiles) GetMetaData(path string) data.Meta {
	args := self.Called(path)
	return args.Get(0).(data.Meta)
}

func (self *MockManagedFiles) WriteMetaData(path string, meta data.Meta) error {
	args := self.Called(path, meta)
	return args.Error(0)
}

func (self *MockManagedFiles) MockAllFilesPresent() {
	self.On("IsGitAppPresent", mock.Anything).Return(true)
	self.On("DirExistsNotEmpty", mock.Anything).Return(true)
	self.On("ExecutableExists", mock.Anything).Return(true)
	self.On("LinkExists", mock.Anything).Return(true)
}

// TODO(jdtls): overhaul mocking a bit to make it less onerous to mock all the status stuff
