package auth_service

import (
	auth_domains "couriers/internal/Auth/core/domains"
	auth_mocks "couriers/internal/Auth/features/service/mocks"
	pkg_errors "couriers/pkg/errors"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestRegisterCourier(t *testing.T) {
	tests := []struct {
		name        string
		courier     auth_domains.Courier
		wantErr     bool
		waitErr     error
		prepareMock func(m *auth_mocks.MockAuthDatabaseRepository)
	}{
		{
			name:    "Success Test",
			courier: auth_domains.NewRegCourier("123456789", "123456789", Moscow),
			wantErr: false,
			waitErr: nil,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, nil).Times(1)
			},
		},

		{
			name:    "Fail validate courier login because it is too short",
			courier: auth_domains.NewRegCourier("1", "123456789", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate courier login because it is too long",
			courier: auth_domains.NewRegCourier(bigLoginPasswordUser, "123456789", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate courier password because it is too short",
			courier: auth_domains.NewRegCourier("123456789", "1", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user password because it is too long",
			courier: auth_domains.NewRegCourier("123456789", bigLoginPasswordUser, Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate courier city because it is too short",
			courier: auth_domains.NewRegCourier("123456789", "123456789", ""),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate courier city because it is too long",
			courier: auth_domains.NewRegCourier("123456789", "123456789", bigLoginPasswordUser),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate courier city because it starts little char",
			courier: auth_domains.NewRegCourier("123456789", "123456789", "moscow"),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail answer from database",
			courier: auth_domains.NewRegCourier("123456789", "123456789", Moscow),
			wantErr: true,
			waitErr: errDatabase,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterCourier(gomock.Any(), gomock.Any()).Return(auth_domains.Courier{}, errDatabase).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mockRepo := auth_mocks.NewMockAuthDatabaseRepository(ctr)
			tt.prepareMock(mockRepo)

			servise := NewAuthService(mockRepo)

			_, err := servise.RegisterCourier(t.Context(), tt.courier)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterCourier returned error=%v wait=%v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("RegisterCourier returned error=%v wait=%v", err, tt.waitErr)
			}
		})
	}
}
