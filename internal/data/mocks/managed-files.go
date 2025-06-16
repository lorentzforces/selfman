package mocks

import (
	"github.com/lorentzforces/selfman/internal/data"
	"github.com/stretchr/testify/mock"
)

type MockManagedFiles struct {
	mock.Mock
}

func (self *MockManagedFiles) AppStatus(appName string) data.AppStatus {
	args := self.Called(appName)
	return args.Get(0).(data.AppStatus)
}

func (self *MockManagedFiles) MockAllPresent(appName string) {
	self.On("AppStatus", appName).Return(data.AppStatus{
		IsConfigured: true,
		SourcePresent: true,
		TargetPresent: true,
		LinkPresent: true,
		VersionPresent: true,
	})
}
