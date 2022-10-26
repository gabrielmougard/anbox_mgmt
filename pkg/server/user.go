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

package server

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"anbox_mgmt/pkg/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type mergedUserWithMetadata struct {
	*models.User
	Metadata []*models.Metadata `json:"metadata"`
}

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

// userResponse is a helper function used to return the User response in the format specified
// by the API spec.
func userResponse(user *models.User, _token ...string) M {
	if user == nil {
		return nil
	}
	var token string
	if len(_token) > 0 {
		token = _token[0]
	}
	return M{
		"email":    user.Email,
		"token":    token,
		"username": user.Username,
	}
}

func (s *Server) createUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    string `json:"email" validate:"required,email"`
			Username string `json:"username" validate:"required,min=2"`
			Age      uint   `json:"age" validate:"required,min=1"`
			Password string `json:"password" validate:"required,min=8,max=72"`
		} `json:"user" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		user := models.User{
			Email:    input.User.Email,
			Username: input.User.Username,
			Age:      input.User.Age,
		}

		user.SetPassword(input.User.Password)

		if err := s.userService.CreateUser(r.Context(), &user); err != nil {
			switch {
			case errors.Is(err, models.ErrDuplicateEmail):
				err = ErrorM{"email": []string{"this email is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			case errors.Is(err, models.ErrDuplicateUsername):
				err = ErrorM{"username": []string{"this username is already in use"}}
				errorResponse(w, http.StatusConflict, err)
			default:
				serverError(w, err)
			}
			return
		}

		userWithMD, err := mergeUserWithGamingMetadata(r.Context(), &user, s.metadataService)
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, M{"user": userWithMD})
	}
}

// Currently, for the sake of simplicity, if a user is logged in,
// it has root privilege and thus has no API restrictions.
// An improvement could be to add an ACL feature to
// fine-tune the API access according to the user group.
func (s *Server) loginUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			errorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.userService.Authenticate(r.Context(), input.User.Email, input.User.Password)

		if err != nil || user == nil {
			invalidUserCredentialsError(w)
			return
		}

		token, err := generateUserToken(user)

		if err != nil {
			serverError(w, err)
			return
		}

		user.Token = token

		userWithMD, err := mergeUserWithGamingMetadata(r.Context(), user, s.metadataService)
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"user": userWithMD})

	}
}

func (s *Server) getCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := userFromContext(ctx)
		user.Token = userTokenFromContext(ctx)

		// TODO : merge metadata here
		writeJSON(w, http.StatusOK, M{"user": user})
	}
}

func (s *Server) listUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := models.UserFilter{}

		if v := query.Get("id"); v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uid := uint(id)
			filter.ID = &uid
		}

		if v := query.Get("email"); v != "" {
			filter.Email = &v
		}

		if v := query.Get("username"); v != "" {
			filter.Username = &v
		}

		if v := query.Get("age"); v != "" {
			age, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uage := uint(age)
			filter.Age = &uage
		}

		users, err := s.userService.Users(r.Context(), filter)

		if err != nil {
			serverError(w, err)
			return
		}

		usersWithMd := []*mergedUserWithMetadata{}
		for _, user := range users {
			userWithMD, err := mergeUserWithGamingMetadata(r.Context(), user, s.metadataService)
			if err != nil {
				serverError(w, err)
				return
			}
			usersWithMd = append(usersWithMd, &userWithMD)
		}
		writeJSON(w, http.StatusOK, M{"users": usersWithMd})
	}
}

func (s *Server) updateUser() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    *string `json:"email,omitempty"`
			Username *string `json:"username,omitempty"`
			Age      *uint   `json:"age,omitempty"`
			Password *string `json:"password,omitempty"`
		} `json:"user,omitempty" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := &Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}

		ctx := r.Context()
		user := userFromContext(ctx)
		patch := models.UserPatch{
			Username: input.User.Username,
			Email:    input.User.Email,
			Age:      input.User.Age,
		}

		if v := input.User.Password; v != nil {
			user.SetPassword(*v)
		}

		err := s.userService.UpdateUser(ctx, user, patch)
		if err != nil {
			serverError(w, err)
			return
		}

		user.Token = userTokenFromContext(ctx)

		userWithMD, err := mergeUserWithGamingMetadata(r.Context(), user, s.metadataService)
		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"user": userWithMD})
	}
}

func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := models.UserFilter{}

		if v := query.Get("id"); v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uid := uint(id)
			filter.ID = &uid
		}

		if v := query.Get("email"); v != "" {
			filter.Email = &v
		}

		if v := query.Get("username"); v != "" {
			filter.Username = &v
		}

		if v := query.Get("age"); v != "" {
			age, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uage := uint(age)
			filter.Age = &uage
		}

		users, err := s.userService.Users(r.Context(), filter)

		if err != nil {
			switch {
			case errors.Is(err, models.ErrNotFound):
				err := ErrorM{"user": []string{"requested user not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}

		user := userFromContext(r.Context())
		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		for _, user := range users {
			if err := s.userService.DeleteUser(r.Context(), user.ID); err != nil {
				serverError(w, err)
				return
			}
		}

		writeJSON(w, http.StatusNoContent, nil)
	}
}

func mergeUserWithGamingMetadata(ctx context.Context, user *models.User, metadataService models.MetadataService) (mergedUserWithMetadata, error) {
	filterMd := models.MetadataFilter{
		PlayerID: &user.ID,
	}
	md, err := metadataService.Metadata(ctx, filterMd)
	if err != nil {
		return mergedUserWithMetadata{}, err
	}
	return mergedUserWithMetadata{user, md}, nil
}
