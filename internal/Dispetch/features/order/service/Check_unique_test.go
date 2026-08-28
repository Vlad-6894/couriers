package dispetch_service

import (
	dispetch_orders_mocks "couriers/internal/Dispetch/features/order/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errFromCache    = errors.New("error from cache")
	errFromDatabase = errors.New("error from database")
	errFromBroker   = errors.New("error from broker")
)

func TestCheckUnique(t *testing.T) {
	tests := []struct {
		name         string
		orderID      int
		wantErr      bool
		waitErr      error
		prepareCache func(m *dispetch_orders_mocks.MockOrdersDispetchCache)
	}{
		{
			name:    "Success test",
			orderID: 1,
			wantErr: false,
			waitErr: nil,
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().CheckUnique(gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
			},
		},

		{
			name:    "fail test",
			orderID: 1,
			wantErr: true,
			waitErr: errFromCache,
			prepareCache: func(m *dispetch_orders_mocks.MockOrdersDispetchCache) {
				m.EXPECT().CheckUnique(gomock.Any(), gomock.Any()).Return(false, errFromCache).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mochCache := dispetch_orders_mocks.NewMockOrdersDispetchCache(ctrl)
			tt.prepareCache(mochCache)

			service := NewOrdersDispetchService(nil, mochCache, nil)

			_, err := service.CheckUnique(t.Context(), tt.orderID)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, but returned err=%v", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("waitErr=%v, but returned err=%v", tt.waitErr, err)
			}
		})
	}
}
