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
	"anbox_mgmt/pkg/models"
	"io"
	"net/http"
	"strings"

	//"anbox_mgmt/pkg/models"
	"github.com/gorilla/handlers"
)

func Logger(w io.Writer) func(h http.Handler) http.Handler {
	return (func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(w, h)
	})
}

func (s *Server) authenticate(mustAuth bool) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				if mustAuth {
					invalidAuthTokenError(w)
				} else {
					r = setContextUser(r, &models.AnonymousUser)
					h.ServeHTTP(w, r)
				}

				return
			}

			ss := strings.Split(authHeader, " ")

			if len(ss) < 2 {
				invalidAuthTokenError(w)
				return
			}

			token := ss[1]

			claims, err := parseUserToken(token)

			if err != nil {
				invalidAuthTokenError(w)
				return
			}

			email := claims["email"].(string)

			user, err := s.userService.UserByEmail(r.Context(), email)

			if err != nil {
				serverError(w, err)
				return
			}

			r = setContextUser(r, user)
			r = setContextUserToken(r, token)
			h.ServeHTTP(w, r)
		})
	}
}
