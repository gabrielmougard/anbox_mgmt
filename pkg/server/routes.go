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
	"os"

	"github.com/rs/cors"
)

func (s *Server) routes() {
	s.router.Use(cors.AllowAll().Handler)
	s.router.Use(Logger(os.Stdout))
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	noAuth := apiRouter.PathPrefix("").Subrouter()
	{
		noAuth.Handle("/health", healthCheck())
		noAuth.Handle("/users", s.createUser()).Methods("POST")
		noAuth.Handle("/users/login", s.loginUser()).Methods("POST")
	}

	authApiRoutes := apiRouter.PathPrefix("").Subrouter()
	authApiRoutes.Use(s.authenticate(true))
	{
		authApiRoutes.Handle("/users", s.listUsers()).Methods("GET")
		authApiRoutes.Handle("/users", s.deleteUser()).Methods("DELETE")
		authApiRoutes.Handle("/users", s.updateUser()).Methods("PUT", "PATCH")

		authApiRoutes.Handle("/games", s.createGames()).Methods("POST")
		authApiRoutes.Handle("/games", s.listGames()).Methods("GET")
		authApiRoutes.Handle("/games", s.deleteGames()).Methods("DELETE")
		authApiRoutes.Handle("/games", s.updateGames()).Methods("PUT", "PATCH")

		authApiRoutes.Handle("/games/link", s.linkGames()).Methods("POST")
	}
}
