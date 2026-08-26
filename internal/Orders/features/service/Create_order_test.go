package orders_service

import (
	orders_domains "couriers/internal/Orders/core/domains"
	orders_mocks "couriers/internal/Orders/features/service/mocks"
	pkg_errors "couriers/pkg/errors"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errDatabase          = errors.New("error from database")
	errBroker            = errors.New("error from broker")
	Moscow               = "Moscow"
	bigLoginPasswordUser = "012345678910111212220324393029395449848749384948398383390000009848494895498494839389899849458989549849949854949849493092309404995090940955094095409554090540909509509554095409"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name                string
		order               orders_domains.Order
		wantErr             bool
		waitErr             error
		prepareDatabaseMock func(m *orders_mocks.MockDatabasePostgres)
		prepareBrokerMock   func(m *orders_mocks.MockBrokerKafka)
	}{
		{
			name:    "success test",
			order:   orders_domains.NewUnitializedOrder("cake", 1000, Moscow, 1),
			wantErr: false,
			waitErr: nil,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(1)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail to validate because order name less 3 chars",
			order:   orders_domains.NewUnitializedOrder("ca", 1000, Moscow, 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail to validate because order name more 20 chars",
			order:   orders_domains.NewUnitializedOrder("ca4dxfhvxzcnhvxsrtrtffgcvkdiaoajdjsghnfgcvg", 1000, Moscow, 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail to validate because order price is negative",
			order:   orders_domains.NewUnitializedOrder("cake", -8, Moscow, 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail to validate because order city haves less 1 chars",
			order:   orders_domains.NewUnitializedOrder("cake", 1000, "", 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail to validate because order city haves more 100 chars",
			order:   orders_domains.NewUnitializedOrder("cake", 1000, bigLoginPasswordUser, 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail to validate because order city begins little char",
			order:   orders_domains.NewUnitializedOrder("ca", 1000, "moscow", 1),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(0)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "err from database",
			order:   orders_domains.NewUnitializedOrder("cake", 1000, Moscow, 1),
			wantErr: true,
			waitErr: errDatabase,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, errDatabase).Times(1)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "err from broker",
			order:   orders_domains.NewUnitializedOrder("cake", 1000, Moscow, 1),
			wantErr: true,
			waitErr: errBroker,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(orders_domains.Order{}, nil).Times(1)
			},
			prepareBrokerMock: func(m *orders_mocks.MockBrokerKafka) {
				m.EXPECT().SendOrder(gomock.Any(), gomock.Any()).Return(errBroker).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := orders_mocks.NewMockDatabasePostgres(ctrl)
			tt.prepareDatabaseMock(mockDB)

			mockBroker := orders_mocks.NewMockBrokerKafka(ctrl)
			tt.prepareBrokerMock(mockBroker)

			service := NewOrdersService(mockDB, mockBroker)

			_, err := service.CreateOrder(ctx, tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err=%v, but returned=%v:", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("wait err=%v, but returned=%v:", tt.waitErr, err)
			}
		})
	}
}
