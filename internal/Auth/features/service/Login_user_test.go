package auth_service

import (
	auth_domains "couriers/internal/Auth/core/domains"
	auth_mocks "couriers/internal/Auth/features/service/mocks"
	pkg_errors "couriers/pkg/errors"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name        string
		login       string
		password    string
		wantErr     bool
		waitErr     error
		prepareMock func(m *auth_mocks.MockAuthDatabaseRepository)
	}{
		{
			name:     "Success test",
			login:    "123456789",
			password: "123456789",
			wantErr:  false,
			waitErr:  nil,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{
					ID:       1,
					Version:  1,
					Login:    "123456789",
					Password: "123456789",
					City:     Moscow,
				}, nil).Times(1)
			},
		},

		{
			name:     "Fail Login user because password is wrong",
			login:    "123456789",
			password: "987654321",
			wantErr:  true,
			waitErr:  pkg_errors.ErrInvalidPassword,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{
					Password: "123456789",
				}, nil).Times(1)
			},
		},

		{

			name:     "Fail answer from database",
			login:    "123456789",
			password: "123456789",
			wantErr:  true,
			waitErr:  errDatabase,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{Password: "12345678"}, errDatabase).Times(1)
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

			key, err := servise.LoginUser(t.Context(), tt.login, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser returned error=%v wait=%v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("RegisterCourier returned error=%v wait=%v", err, tt.waitErr)
			}

			if err == nil && key == "" {
				t.Errorf("empty jwt key error")
			}
		})
	}
}
