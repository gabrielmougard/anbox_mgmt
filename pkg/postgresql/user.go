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

package postgresql

import (
	"anbox_mgmt/pkg/models"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func (us *UserService) UserByID(ctx context.Context, id uint) (*models.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	user, err := findOneUser(ctx, tx, models.UserFilter{ID: &id})

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	user, err := findOneUser(ctx, tx, models.UserFilter{Email: &email})

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) UserByUsername(ctx context.Context, uname string) (*models.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	user, err := findOneUser(ctx, tx, models.UserFilter{Username: &uname})

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Users(ctx context.Context, uf models.UserFilter) ([]*models.User, error) {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	users, err := findUsers(ctx, tx, uf)

	if err != nil {
		return nil, err
	}

	return users, tx.Commit()
}

func (us *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := us.UserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, models.ErrUnAuthorized
	}

	return user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *models.User, patch models.UserPatch) error {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	defer tx.Rollback()

	if err := updateUser(ctx, tx, user, patch); err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	return nil
}

func (us *UserService) DeleteUser(ctx context.Context, id uint) error {
	tx, err := us.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = deleteUser(ctx, tx, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func deleteUser(ctx context.Context, tx *sqlx.Tx, id uint) error {
	query := "DELETE FROM users WHERE id = $1"
	return execQuery(ctx, tx, query, id)
}

func createUser(ctx context.Context, tx *sqlx.Tx, user *models.User) error {
	query := `
	INSERT INTO users (email, username, age, password_hash)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	args := []interface{}{user.Email, user.Username, user.Age, user.PasswordHash}
	err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return models.ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return models.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func findUserByID(ctx context.Context, tx *sqlx.Tx, id uint) (*models.User, error) {
	return findOneUser(ctx, tx, models.UserFilter{ID: &id})
}

func findOneUser(ctx context.Context, tx *sqlx.Tx, filter models.UserFilter) (*models.User, error) {
	us, err := findUsers(ctx, tx, filter)

	if err != nil {
		return nil, err
	} else if len(us) == 0 {
		return nil, models.ErrNotFound
	}

	return us[0], nil
}

func findUsers(ctx context.Context, tx *sqlx.Tx, filter models.UserFilter) ([]*models.User, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Email; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("email = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Username; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("username = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Age; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("age = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from users" + formatWhereClause(where) +
		" ORDER BY id ASC" + formatLimitOffset(filter.Limit, filter.Offset)

	users, err := queryUsers(ctx, tx, query, args...)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func updateUser(ctx context.Context, tx *sqlx.Tx, user *models.User, patch models.UserPatch) error {
	if v := patch.Email; v != nil {
		user.Email = *v
	}

	if v := patch.PasswordHash; v != nil {
		user.PasswordHash = *v
	}

	if v := patch.Username; v != nil {
		user.Username = *v
	}

	if v := patch.Age; v != nil {
		user.Age = *v
	}

	args := []interface{}{
		user.Username,
		user.Email,
		user.Age,
		user.PasswordHash,
		user.ID,
	}

	query := `
	UPDATE users 
	SET username = $1, email = $2, age = $3, password_hash = $4, updated_at = NOW()
	WHERE id = $5
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&user.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return models.ErrInternal
	}

	return nil
}

func queryUsers(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) ([]*models.User, error) {
	users := make([]*models.User, 0)

	if err := findMany(ctx, tx, &users, query, args...); err != nil {
		return users, err
	}

	return users, nil
}
