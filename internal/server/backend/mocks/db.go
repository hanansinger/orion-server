// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	config "github.ibm.com/blockchaindb/server/config"

	types "github.ibm.com/blockchaindb/server/pkg/types"

	x509 "crypto/x509"
)

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

// BootstrapDB provides a mock function with given fields: conf
func (_m *DB) BootstrapDB(conf *config.Configurations) error {
	ret := _m.Called(conf)

	var r0 error
	if rf, ok := ret.Get(0).(func(*config.Configurations) error); ok {
		r0 = rf(conf)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *DB) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DoesUserExist provides a mock function with given fields: userID
func (_m *DB) DoesUserExist(userID string) (bool, error) {
	ret := _m.Called(userID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeader provides a mock function with given fields: userId, blockNum
func (_m *DB) GetBlockHeader(userId string, blockNum uint64) (*types.GetBlockResponseEnvelope, error) {
	ret := _m.Called(userId, blockNum)

	var r0 *types.GetBlockResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64) *types.GetBlockResponseEnvelope); ok {
		r0 = rf(userId, blockNum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetBlockResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64) error); ok {
		r1 = rf(userId, blockNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCertificate provides a mock function with given fields: userID
func (_m *DB) GetCertificate(userID string) (*x509.Certificate, error) {
	ret := _m.Called(userID)

	var r0 *x509.Certificate
	if rf, ok := ret.Get(0).(func(string) *x509.Certificate); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*x509.Certificate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields:
func (_m *DB) GetConfig() (*types.GetConfigResponseEnvelope, error) {
	ret := _m.Called()

	var r0 *types.GetConfigResponseEnvelope
	if rf, ok := ret.Get(0).(func() *types.GetConfigResponseEnvelope); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetConfigResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDBStatus provides a mock function with given fields: dbName
func (_m *DB) GetDBStatus(dbName string) (*types.GetDBStatusResponseEnvelope, error) {
	ret := _m.Called(dbName)

	var r0 *types.GetDBStatusResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetDBStatusResponseEnvelope); ok {
		r0 = rf(dbName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDBStatusResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dbName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetData provides a mock function with given fields: dbName, querierUserID, key
func (_m *DB) GetData(dbName string, querierUserID string, key string) (*types.GetDataResponseEnvelope, error) {
	ret := _m.Called(dbName, querierUserID, key)

	var r0 *types.GetDataResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string, string) *types.GetDataResponseEnvelope); ok {
		r0 = rf(dbName, querierUserID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetDataResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(dbName, querierUserID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLedgerPath provides a mock function with given fields: userId, start, end
func (_m *DB) GetLedgerPath(userId string, start uint64, end uint64) (*types.GetLedgerPathResponseEnvelope, error) {
	ret := _m.Called(userId, start, end)

	var r0 *types.GetLedgerPathResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, uint64, uint64) *types.GetLedgerPathResponseEnvelope); ok {
		r0 = rf(userId, start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetLedgerPathResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint64, uint64) error); ok {
		r1 = rf(userId, start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNodeConfig provides a mock function with given fields: nodeID
func (_m *DB) GetNodeConfig(nodeID string) (*types.GetNodeConfigResponseEnvelope, error) {
	ret := _m.Called(nodeID)

	var r0 *types.GetNodeConfigResponseEnvelope
	if rf, ok := ret.Get(0).(func(string) *types.GetNodeConfigResponseEnvelope); ok {
		r0 = rf(nodeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetNodeConfigResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(nodeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: querierUserID, targetUserID
func (_m *DB) GetUser(querierUserID string, targetUserID string) (*types.GetUserResponseEnvelope, error) {
	ret := _m.Called(querierUserID, targetUserID)

	var r0 *types.GetUserResponseEnvelope
	if rf, ok := ret.Get(0).(func(string, string) *types.GetUserResponseEnvelope); ok {
		r0 = rf(querierUserID, targetUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.GetUserResponseEnvelope)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(querierUserID, targetUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Height provides a mock function with given fields:
func (_m *DB) Height() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDBExists provides a mock function with given fields: name
func (_m *DB) IsDBExists(name string) bool {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsReady provides a mock function with given fields:
func (_m *DB) IsReady() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LedgerHeight provides a mock function with given fields:
func (_m *DB) LedgerHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitTransaction provides a mock function with given fields: tx
func (_m *DB) SubmitTransaction(tx interface{}) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
