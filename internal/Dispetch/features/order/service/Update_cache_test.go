package dispetch_service

import (
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	dispetch_orders_mocks "couriers/internal/Dispetch/features/order/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestUpdateCache(t *testing.T) {
	tests := []struct {
		name            string
		wantErr         bool
		waitErr         error
		prepareDatabase func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase)
		prepareCache    func(m *dispetch_orders_mocks.MockOrdersDispetchCache)
	}{
		{
			name:    "success test",
			wantErr: false,
			waitErr: nil,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().GetFreeCouriers(gomock.Any()).Return([]dispetch_domains.CourierInfo{dispetch_domains.NewCourierInfo(1, 1, "Moscow"), dispetch_domains.NewCourierInfo(2, 1, "Moscow"), dispetch_domains.NewCourierInfo(1, 1, "London"), dispetch_domains.NewCourierInfo(2, 1, "London")}, nil).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().UpdateCache(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail from database",
			wantErr: true,
			waitErr: errFromDatabase,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().GetFreeCouriers(gomock.Any()).Return(nil, errFromDatabase).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().UpdateCache(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "success test",
			wantErr: true,
			waitErr: errFromCache,
			prepareDatabase: func(m *dispetch_orders_mocks.MockOrdersDispetchDatabase) {
				m.EXPECT().GetFreeCouriers(gomock.Any()).Return([]dispetch_domains.CourierInfo{dispetch_domains.NewCourierInfo(1, 1, "Moscow"), dispetch_domains.NewCourierInfo(2, 1, "Moscow"), dispetch_domains.NewCourierInfo(1, 1, "London"), dispetch_domains.NewCourierInfo(2, 1, "London")}, nil).Times(1)
			},
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().UpdateCache(gomock.Any(), gomock.Any()).Return(errFromCache).Times(1)
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

			service := NewOrdersDispetchService(mockDatabase, mochCache, nil)

			err := service.UpdateCache(t.Context())

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, but returned err=%v", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("waitErr=%v, but returned err=%v", tt.waitErr, err)
			}
		})
	}
}
