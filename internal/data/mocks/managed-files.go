package mocks

import (
	"path"

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

func StartMockingManagedFiles(sysConf *data.SystemConfig) *managedFilesMockOpts {
	return &managedFilesMockOpts{
		sysConf,
	}
}

type managedFilesMockOpts struct {
	sysConf *data.SystemConfig
}
type managedFilesMockOptFunc func(*MockManagedFiles)

func (self *managedFilesMockOpts) SetMocks(opts ...managedFilesMockOptFunc) *MockManagedFiles {
	mock := &MockManagedFiles{}
	for _, opt := range opts {
		opt(mock)
	}
	return mock
}

func (self *managedFilesMockOpts) GitAppPresent(
	appName string,
	present bool,
) managedFilesMockOptFunc {
	return func(mock *MockManagedFiles) {
		mock.On("IsGitAppPresent", path.Join(self.sysConf.SourcesPath(), appName)).Return(present)
	}
}

func (self *managedFilesMockOpts) ExecutableExists(
	appName string,
	exists bool,
) managedFilesMockOptFunc {
	return func(mock *MockManagedFiles) {
		mock.On("ExecutableExists", path.Join(self.sysConf.ArtifactsPath(), appName)).Return(exists)
	}
}

func (self *managedFilesMockOpts) LinkExists(
	appName string,
	exists bool,
) managedFilesMockOptFunc {
	return func(mock *MockManagedFiles) {
		mock.On("LinkExists", path.Join(*self.sysConf.BinaryDir, appName)).Return(exists)
	}
}

func (self *managedFilesMockOpts) MetaData(
	appName string,
	metadata *data.Meta,
) managedFilesMockOptFunc {
	return func(mock *MockManagedFiles) {
		mock.On(
			"GetMetaData",
			path.Join(self.sysConf.MetaPath(), data.MetaFileNameForApp(appName)),
		).Return(*metadata)
	}
}
