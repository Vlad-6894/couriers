package dispetch_confirm_service

import (
	dispetch_domains "couriers/internal/Dispetch/core/domains"
	dispetch_confirm_mocks "couriers/internal/Dispetch/features/confirm/service/mocks"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errDB = errors.New("error from database")
)

func TestConfirmOrder(t *testing.T) {
	tests := []struct {
		name            string
		confirm         dispetch_domains.Confirm
		wantErr         bool
		waitErr         error
		prepareDatabase func(m *dispetch_confirm_mocks.MockDispetchPostgresDatabase)
	}{
		{
			name:    "success test",
			confirm: dispetch_domains.NewConfirm(1, 1),
			wantErr: false,
			waitErr: nil,
			prepareDatabase: func(m *dispetch_confirm_mocks.MockDispetchPostgresDatabase) {
				m.EXPECT().ConfirmOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:    "fail test",
			confirm: dispetch_domains.NewConfirm(1, 1),
			wantErr: true,
			waitErr: errDB,
			prepareDatabase: func(m *dispetch_confirm_mocks.MockDispetchPostgresDatabase) {
				m.EXPECT().ConfirmOrder(gomock.Any(), gomock.Any()).Return(errDB).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mochDatabase := dispetch_confirm_mocks.NewMockDispetchPostgresDatabase(ctrl)
			tt.prepareDatabase(mochDatabase)

			service := NewDispetchConfirmService(mochDatabase)

			err := service.ConfirmOrder(t.Context(), tt.confirm)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, but returned err=%v", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("waitErr=%v, but returned err=%v", tt.waitErr, err)
			}
		})
	}
}
