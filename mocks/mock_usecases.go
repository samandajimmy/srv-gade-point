package mocks

import gomock "github.com/golang/mock/gomock"

type MockUsecases struct {
	MockQUs   *MockQUseCase
	MockTUs   *MockTUseCase
	MockVUs   *MockVUsecase
	MockRefUs *MockRefUseCase
}

func NewMockUsecases(mockCtrl *gomock.Controller) MockUsecases {
	return MockUsecases{
		MockQUs:   NewMockQUseCase(mockCtrl),
		MockTUs:   NewMockTUseCase(mockCtrl),
		MockVUs:   NewMockVUsecase(mockCtrl),
		MockRefUs: NewMockRefUseCase(mockCtrl),
	}
}
