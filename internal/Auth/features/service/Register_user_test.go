package auth_service

import (
	auth_domains "couriers/internal/Auth/core/domains"
	auth_mocks "couriers/internal/Auth/features/service/mocks"
	pkg_errors "couriers/pkg/errors"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	errDatabase          = errors.New("Any err from database")
	Moscow               = "Moscow"
	bigLoginPasswordUser = "012345678910111212220324393029395449848749384948398383390000009848494895498494839389899849458989549849949854949849493092309404995090940955094095409554090540909509509554095409"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name        string
		user        auth_domains.User
		wantErr     bool
		waitErr     error
		prepareMock func(m *auth_mocks.MockAuthDatabaseRepository)
	}{
		{
			name:    "Success Test",
			user:    auth_domains.NewRegUser("123456789", "123456789", Moscow),
			wantErr: false,
			waitErr: nil,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, nil).Times(1)
			},
		},

		{
			name:    "Fail validate user login because it is too short",
			user:    auth_domains.NewRegUser("1", "123456789", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user login because it is too long",
			user:    auth_domains.NewRegUser(bigLoginPasswordUser, "123456789", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user password because it is too short",
			user:    auth_domains.NewRegUser("123456789", "1", Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user password because it is too long",
			user:    auth_domains.NewRegUser("123456789", bigLoginPasswordUser, Moscow),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user city because it is too short",
			user:    auth_domains.NewRegUser("123456789", "123456789", ""),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user city because it is too long",
			user:    auth_domains.NewRegUser("123456789", "123456789", bigLoginPasswordUser),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail validate user city because it starts little char",
			user:    auth_domains.NewRegUser("123456789", "123456789", "moscow"),
			wantErr: true,
			waitErr: pkg_errors.ErrInvalidArgument,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, pkg_errors.ErrInvalidArgument).Times(0)
			},
		},

		{
			name:    "Fail answer from database",
			user:    auth_domains.NewRegUser("123456789", "123456789", Moscow),
			wantErr: true,
			waitErr: errDatabase,
			prepareMock: func(m *auth_mocks.MockAuthDatabaseRepository) {
				m.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return(auth_domains.User{}, errDatabase).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mockRepo := auth_mocks.NewMockAuthDatabaseRepository(ctr)

			tt.prepareMock(mockRepo)

			authService := NewAuthService(mockRepo)

			_, err := authService.RegisterUser(t.Context(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser returned error=%v wait=%v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("RegisterUser returned error=%v wait=%v", err, tt.waitErr)
			}
		})
	}
}
