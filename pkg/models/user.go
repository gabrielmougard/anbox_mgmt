// Copyright 2022 gab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"-"`
	Email        string    `json:"email,omitempty"`
	Age          uint      `json:"age,omitempty"`
	Username     string    `json:"username,omitempty"`
	Token        string    `json:"token,omitempty"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	UpdatedAt    time.Time `json:"-" db:"updated_at"`
}

var AnonymousUser User

type UserFilter struct {
	ID       *uint
	Email    *string
	Username *string
	Age      *uint

	Limit  int
	Offset int
}

type UserPatch struct {
	Email        *string `json:"email"`
	Username     *string `json:"username"`
	Age          *uint   `json:"age"`
	PasswordHash *string `json:"-" db:"password_hash"`
}

func (u *User) SetPassword(password string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		// return better error message
		return err
	}

	u.PasswordHash = string(hashBytes)

	return nil
}

func (u User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) IsAnonymous() bool {
	return u == &AnonymousUser
}

type UserService interface {
	Authenticate(ctx context.Context, email, password string) (*User, error)

	CreateUser(context.Context, *User) error

	UserByID(context.Context, uint) (*User, error)

	UserByEmail(context.Context, string) (*User, error)

	UserByUsername(context.Context, string) (*User, error)

	Users(context.Context, UserFilter) ([]*User, error)

	UpdateUser(context.Context, *User, UserPatch) error

	DeleteUser(context.Context, uint) error
}
