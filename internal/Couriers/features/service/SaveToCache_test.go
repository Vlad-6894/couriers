package couriers_service

import (
	courier_domains "couriers/internal/Couriers/core/domains"
	couriers_mocks "couriers/internal/Couriers/features/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errFromCache  = errors.New("error from cache")
	errFromBroker = errors.New("error from broker")
)

func TestSaveToCache(t *testing.T) {
	tests := []struct {
		name         string
		order        courier_domains.DispetchedOrder
		wantErr      bool
		waitErr      error
		prepareCache func(m *couriers_mocks.MockCouriersCache)
	}{
		{
			name:    "Successful test",
			order:   courier_domains.NewDispetchedOrder(1, 1, 1),
			wantErr: false,
			waitErr: nil,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().SaveToCache(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail test",
			order:   courier_domains.NewDispetchedOrder(1, 1, 1),
			wantErr: true,
			waitErr: errFromCache,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().SaveToCache(gomock.Any(), gomock.Any()).Return(errFromCache).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := couriers_mocks.NewMockCouriersCache(ctrl)
			tt.prepareCache(mockCache)

			service := NewCouriersService(mockCache, nil)

			err := service.SaveToCache(t.Context(), tt.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToCache returned err=%v, but wantErr=%v:", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("SaveToCache returned err=%v, but waitErr=%v:", err, tt.waitErr)
			}
		})
	}
}
