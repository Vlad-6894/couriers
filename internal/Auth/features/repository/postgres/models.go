package auth_repository_postgres

import auth_domains "couriers/internal/Auth/core/domains"

type UserModel struct {
	ID       int
	Version  int
	Login    string
	Password string
	City     string
}

type CourierModel struct {
	ID             int
	Version        int
	Login          string
	Password       string
	City           string
	OrdersComplete int
	IsFree         bool
}

func NewUserDomainFromModel(userModel UserModel) auth_domains.User {
	return auth_domains.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.Login,
		userModel.Password,
		userModel.City,
	)
}

func NewCourierDomainFromModel(courierModel CourierModel) auth_domains.Courier {
	return auth_domains.NewCourier(
		courierModel.ID,
		courierModel.Version,
		courierModel.Login,
		courierModel.Password,
		courierModel.City,
		courierModel.OrdersComplete,
		courierModel.IsFree,
	)
}
