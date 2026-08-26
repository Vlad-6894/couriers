package orders_service

import (
	orders_domains "couriers/internal/Orders/core/domains"
	orders_mocks "couriers/internal/Orders/features/service/mocks"
	pkg_errors "couriers/pkg/errors"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

var (
	rightLimit  = 12
	rightOffset = 12
	failLimit   = -4
	failOffset  = -4
)

func TestGetOrders(t *testing.T) {
	tests := []struct {
		name                string
		personID            int
		limit               *int
		offset              *int
		wantErr             bool
		waitErr             error
		prepareDatabaseMock func(m *orders_mocks.MockDatabasePostgres)
		prepareBrokerMock   func(m *orders_mocks.MockBrokerKafka)
	}{
		{
			name:     "success test",
			personID: 1,
			limit:    &rightLimit,
			offset:   &rightOffset,
			wantErr:  false,
			waitErr:  nil,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, nil).Times(1)
			},
			prepareBrokerMock: nil,
		},

		{
			name:     "fail limit",
			personID: 1,
			limit:    &failLimit,
			offset:   &rightOffset,
			wantErr:  true,
			waitErr:  pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, nil).Times(0)
			},
			prepareBrokerMock: nil,
		},

		{
			name:     "fail offset",
			personID: 1,
			limit:    &rightLimit,
			offset:   &failOffset,
			wantErr:  true,
			waitErr:  pkg_errors.ErrInvalidArgument,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, nil).Times(0)
			},
			prepareBrokerMock: nil,
		},

		{
			name:     "nil limit",
			personID: 1,
			limit:    nil,
			offset:   &rightOffset,
			wantErr:  false,
			waitErr:  nil,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, nil).Times(1)
			},
			prepareBrokerMock: nil,
		},

		{
			name:     "nil offset",
			personID: 1,
			limit:    &rightLimit,
			offset:   nil,
			wantErr:  false,
			waitErr:  nil,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, nil).Times(1)
			},
			prepareBrokerMock: nil,
		},

		{
			name:     "error from database",
			personID: 1,
			limit:    &rightLimit,
			offset:   &rightOffset,
			wantErr:  true,
			waitErr:  errDatabase,
			prepareDatabaseMock: func(m *orders_mocks.MockDatabasePostgres) {
				m.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]orders_domains.GetOrder{}, errDatabase).Times(1)
			},
			prepareBrokerMock: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			ctr := gomock.NewController(t)
			defer ctr.Finish()

			mockDB := orders_mocks.NewMockDatabasePostgres(ctr)
			tt.prepareDatabaseMock(mockDB)

			service := NewOrdersService(mockDB, nil)

			_, err := service.GetOrders(ctx, tt.personID, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err=%v, but returned=%v:", tt.wantErr, err)
			}

			if !errors.Is(err, tt.waitErr) {
				t.Errorf("wait err=%v, but returned=%v:", tt.waitErr, err)
			}
		})
	}
}
