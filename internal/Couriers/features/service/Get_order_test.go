package couriers_service

import (
	courier_domains "couriers/internal/Couriers/core/domains"
	couriers_mocks "couriers/internal/Couriers/features/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetOrder(t *testing.T) {
	tests := []struct {
		name         string
		wantErr      bool
		waitErr      error
		prepareCache func(m *couriers_mocks.MockCouriersCache)
	}{
		{
			name:    "Successful test",
			wantErr: false,
			waitErr: nil,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, nil).Times(1)
			},
		},

		{
			name:    "fail test",
			wantErr: true,
			waitErr: errFromCache,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, errFromCache).Times(1)
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

			_, err := service.GetOrder(t.Context(), 1)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToCache returned err=%v, but wantErr=%v:", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("SaveToCache returned err=%v, but waitErr=%v:", err, tt.waitErr)
			}
		})
	}
}
