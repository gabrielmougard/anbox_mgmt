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
			authToken := r.Header.Get("Authorization")

			if authToken == "" {
				if mustAuth {
					invalidAuthTokenError(w)
				} else {
					r = setContextUser(r, &models.AnonymousUser)
					h.ServeHTTP(w, r)
				}

				return
			}

			claims, err := parseUserToken(authToken)

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
			r = setContextUserToken(r, authToken)
			h.ServeHTTP(w, r)
		})
	}
}
