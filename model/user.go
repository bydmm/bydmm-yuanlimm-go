package model

import "time"

// User 用户
type User struct {
	ID             uint `gorm:"primary_key"`
	Name           string
	NickName       string
	Demo           bool
	Address        string
	PasswordDigest string
	Token          string `sql:"unique_index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
