package services

type IUser interface {
	GetValue() int
}

type User struct {
	A int
}

func (u User) GetValue() int {
	return u.A
}

func UserServiceFactory() IUser {
	return &User{
		A: 10,
	}
}
