package couriers_service

import (
	courier_domains "couriers/internal/Couriers/core/domains"
	couriers_mocks "couriers/internal/Couriers/features/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestConfirm(t *testing.T) {
	tests := []struct {
		name          string
		wantErr       bool
		waitErr       error
		prepareCache  func(m *couriers_mocks.MockCouriersCache)
		prepareBroker func(m *couriers_mocks.MockCouriersBroker)
	}{
		{
			name:    "Successful test",
			wantErr: false,
			waitErr: nil,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, nil).Times(1)
				m.EXPECT().DeleteOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			prepareBroker: func(m *couriers_mocks.MockCouriersBroker) {
				m.EXPECT().SendConfirm(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail test because SendConfirm",
			wantErr: true,
			waitErr: errFromBroker,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, nil).Times(1)
				m.EXPECT().DeleteOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			prepareBroker: func(m *couriers_mocks.MockCouriersBroker) {
				m.EXPECT().SendConfirm(gomock.Any(), gomock.Any()).Return(errFromBroker).Times(1)
			},
		},

		{
			name:    "fail test because DeleteOrder",
			wantErr: true,
			waitErr: errFromCache,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, nil).Times(1)
				m.EXPECT().DeleteOrder(gomock.Any(), gomock.Any()).Return(errFromCache).Times(1)
			},
			prepareBroker: func(m *couriers_mocks.MockCouriersBroker) {
				m.EXPECT().SendConfirm(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},

		{
			name:    "fail test because GetOrder",
			wantErr: true,
			waitErr: errFromCache,
			prepareCache: func(m *couriers_mocks.MockCouriersCache) {
				m.EXPECT().GetOrder(gomock.Any(), gomock.Any()).Return(courier_domains.DispetchedOrder{}, errFromCache).Times(1)
				m.EXPECT().DeleteOrder(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			prepareBroker: func(m *couriers_mocks.MockCouriersBroker) {
				m.EXPECT().SendConfirm(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := couriers_mocks.NewMockCouriersCache(ctrl)
			tt.prepareCache(mockCache)

			mockBroker := couriers_mocks.NewMockCouriersBroker(ctrl)
			tt.prepareBroker(mockBroker)

			service := NewCouriersService(mockCache, mockBroker)

			err := service.Confirm(t.Context(), 1)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToCache returned err=%v, but wantErr=%v:", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("SaveToCache returned err=%v, but waitErr=%v:", err, tt.waitErr)
			}
		})
	}
}
