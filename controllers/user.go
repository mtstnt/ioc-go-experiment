package controllers

import (
	"database/sql"

	"github.com/mtstnt/ioc-go-experiment/services"
)

type User struct {
	UserService services.IUser `inject:""`
	MyDB        *sql.DB        `inject:""`
}

func (u *User) GetA() int {
	return u.UserService.GetValue()
}
