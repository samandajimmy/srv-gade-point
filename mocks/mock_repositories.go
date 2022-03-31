package mocks

import gomock "github.com/golang/mock/gomock"

type MockRepositories struct {
	MockCRp       *MockCRepository
	MockRefTRp    *MockRefTRepository
	MockRRp       *MockRRepository
	MockRtRp      *MockRtRepository
	MockVcRp      *MockVcRepository
	MockRefRp     *MockRefRepository
	MockRestRefRp *MockRestRefRepository
}

func NewMockRepository(mockCtrl *gomock.Controller) MockRepositories {
	return MockRepositories{
		MockCRp:       NewMockCRepository(mockCtrl),
		MockRefTRp:    NewMockRefTRepository(mockCtrl),
		MockRRp:       NewMockRRepository(mockCtrl),
		MockRtRp:      NewMockRtRepository(mockCtrl),
		MockVcRp:      NewMockVcRepository(mockCtrl),
		MockRefRp:     NewMockRefRepository(mockCtrl),
		MockRestRefRp: NewMockRestRefRepository(mockCtrl),
	}
}
