package dispetch_service

import (
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	dispetch_orders_mocks "couriers/internal/Dispetch/features/order/service/mocks"
	pkg_repository_redis "couriers/pkg/repository/redis"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestDispetchOrder(t *testing.T) {
	tests := []struct {
		name            string
		order           dispetch_domains.Order
		wantErr         bool
		waitErr         error
		prepareDatabase func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase)
		prepareCache    func(m *dispetch_orders_mocks.MockOrdersDispetchCache)
		prepareBroker   func(m *dispetch_orders_mocks.MockOrdersDispetchBroker)
	}{
		{
			name:    "success test search from cache",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: false,
			waitErr: nil,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(0)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "success test search from database, but cache is empty",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: false,
			waitErr: nil,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, pkg_repository_redis.ErrEmpty).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail test because cache returned error",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: true,
			waitErr: errFromCache,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(0)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, errFromCache).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail test because database SearchCourier returned error",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: true,
			waitErr: errFromDatabase,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, errFromDatabase).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, pkg_repository_redis.ErrEmpty).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail test because broker after database SendOrder returned error",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: true,
			waitErr: errFromBroker,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, pkg_repository_redis.ErrEmpty).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errFromBroker).Times(1)
			},
		},

		{
			name:    "fail test because database DoBusy returned error",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: true,
			waitErr: errFromDatabase,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(errFromDatabase).Times(1)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(0)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail test because broker after database DoBusy returned error",
			order:   dispetch_domains.NewOrder(1, 1, "pizza", 1000, false, "Moscow", 1),
			wantErr: true,
			waitErr: errFromBroker,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().DoBusy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(0)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().SearchCourier(gomock.Any(), gomock.Any()).Return(1, 1, nil).Times(1)
			},
			prepareBroker: func(m *dispetch_orders_mocks.MockOrdersDispetchBroker) {
				m.EXPECT().SendDispetchedOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errFromBroker).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatabase := dispetch_orders_mocks.NewMockOrdersDispetchDatabase(ctrl)
			tt.prepareDatabase(mockDatabase)

			mochCache := dispetch_orders_mocks.NewMockOrdersDispetchCache(ctrl)
			tt.prepareCache(mochCache)

			mockBroker := dispetch_orders_mocks.NewMockOrdersDispetchBroker(ctrl)
			tt.prepareBroker(mockBroker)

			service := NewOrdersDispetchService(mockDatabase, mochCache, mockBroker)

			err := service.DispetchOrder(t.Context(), tt.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, but returned err=%v", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("waitErr=%v, but returned err=%v", tt.waitErr, err)
			}
		})
	}
}
