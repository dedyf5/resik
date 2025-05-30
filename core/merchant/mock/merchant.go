// Code generated by MockGen. DO NOT EDIT.
// Source: merchant.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	dto "github.com/dedyf5/resik/core/merchant/dto"
	ctx "github.com/dedyf5/resik/ctx"
	merchant "github.com/dedyf5/resik/entities/merchant"
	param "github.com/dedyf5/resik/entities/merchant/param"
	response "github.com/dedyf5/resik/pkg/response"
	gomock "github.com/golang/mock/gomock"
)

// MockIService is a mock of IService interface.
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService.
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance.
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// MerchantDelete mocks base method.
func (m *MockIService) MerchantDelete(ctx *ctx.Ctx, param *merchant.Merchant) (bool, *response.Status) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MerchantDelete", ctx, param)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*response.Status)
	return ret0, ret1
}

// MerchantDelete indicates an expected call of MerchantDelete.
func (mr *MockIServiceMockRecorder) MerchantDelete(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MerchantDelete", reflect.TypeOf((*MockIService)(nil).MerchantDelete), ctx, param)
}

// MerchantInsert mocks base method.
func (m *MockIService) MerchantInsert(ctx *ctx.Ctx, merchant *merchant.Merchant) (bool, *response.Status) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MerchantInsert", ctx, merchant)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*response.Status)
	return ret0, ret1
}

// MerchantInsert indicates an expected call of MerchantInsert.
func (mr *MockIServiceMockRecorder) MerchantInsert(ctx, merchant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MerchantInsert", reflect.TypeOf((*MockIService)(nil).MerchantInsert), ctx, merchant)
}

// MerchantListGet mocks base method.
func (m *MockIService) MerchantListGet(param *param.MerchantListGet) (*dto.MerchantList, *response.Status) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MerchantListGet", param)
	ret0, _ := ret[0].(*dto.MerchantList)
	ret1, _ := ret[1].(*response.Status)
	return ret0, ret1
}

// MerchantListGet indicates an expected call of MerchantListGet.
func (mr *MockIServiceMockRecorder) MerchantListGet(param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MerchantListGet", reflect.TypeOf((*MockIService)(nil).MerchantListGet), param)
}

// MerchantUpdate mocks base method.
func (m *MockIService) MerchantUpdate(ctx *ctx.Ctx, merchant *merchant.Merchant) (bool, *response.Status) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MerchantUpdate", ctx, merchant)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(*response.Status)
	return ret0, ret1
}

// MerchantUpdate indicates an expected call of MerchantUpdate.
func (mr *MockIServiceMockRecorder) MerchantUpdate(ctx, merchant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MerchantUpdate", reflect.TypeOf((*MockIService)(nil).MerchantUpdate), ctx, merchant)
}
