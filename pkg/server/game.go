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
	"errors"
	"net/http"
	"strconv"

	"anbox_mgmt/pkg/models"
)

func (s *Server) createGames() http.HandlerFunc {
	type Input struct {
		Game struct {
			Title       string `json:"title" validate:"required"`
			Description string `json:"description"`
			URL         string `json:"url"`
			AgeRating   uint   `json:"ageRating"`
			Publisher   string `json:"publisher"`
		} `json:"game"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.Game); err != nil {
			validationError(w, err)
			return
		}

		game := models.Game{
			Title:       input.Game.Title,
			Description: input.Game.Description,
			URL:         input.Game.URL,
			AgeRating:   input.Game.AgeRating,
			Publisher:   input.Game.Publisher,
		}

		user := userFromContext(r.Context())

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		if err := s.gameService.CreateGame(r.Context(), &game); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"game": game})
	}
}

func (s *Server) listGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := models.GameFilter{}

		if v := query.Get("title"); v != "" {
			filter.Title = &v
		}

		if v := query.Get("desc"); v != "" {
			filter.Description = &v
		}

		if v := query.Get("url"); v != "" {
			filter.URL = &v
		}

		if v := query.Get("age"); v != "" {
			age, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uage := uint(age)
			filter.AgeRating = &uage
		}

		if v := query.Get("publisher"); v != "" {
			filter.Publisher = &v
		}

		games, err := s.gameService.Games(r.Context(), filter)

		if err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"games": games, "gamesCount": len(games)})
	}
}

func (s *Server) updateGames() http.HandlerFunc {
	type Input struct {
		Game struct {
			Title       *string `json:"title,omitempty"`
			Description *string `json:"description,omitempty"`
			URL         *string `json:"url,omitempty"`
			AgeRating   *uint   `json:"ageRating,omitempty"`
			Publisher   *string `json:"publisher,omitempty"`
		} `json:"game,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := models.GameFilter{}

		if v := query.Get("title"); v != "" {
			filter.Title = &v
		}
		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		games, err := s.gameService.Games(r.Context(), filter)

		if err != nil {
			switch {
			case errors.Is(err, models.ErrNotFound):
				err := ErrorM{"game": []string{"requested game not found"}}
				notFoundError(w, err)
			default:
				serverError(w, err)
			}
			return
		}
		game := games[0]
		user := userFromContext(r.Context())
		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		patch := models.GamePatch{
			Title:       input.Game.Title,
			Description: input.Game.Description,
			URL:         input.Game.URL,
			AgeRating:   input.Game.AgeRating,
			Publisher:   input.Game.Publisher,
		}

		if err := s.gameService.UpdateGame(r.Context(), game, patch); err != nil {
			serverError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, M{"game": game})
	}
}

func (s *Server) deleteGames() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filter := models.GameFilter{}

		if v := query.Get("title"); v != "" {
			filter.Title = &v
		}

		if v := query.Get("desc"); v != "" {
			filter.Description = &v
		}

		if v := query.Get("url"); v != "" {
			filter.URL = &v
		}

		if v := query.Get("age_rating"); v != "" {
			age, err := strconv.Atoi(v)
			if err != nil {
				serverError(w, err)
				return
			}
			uage := uint(age)
			filter.AgeRating = &uage
		}

		if v := query.Get("publisher"); v != "" {
			filter.Publisher = &v
		}

		games, err := s.gameService.Games(r.Context(), filter)

		if err != nil {
			switch {
			case errors.Is(err, models.ErrNotFound):
				err := ErrorM{"game": []string{"requested game not found"}}
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

		for _, game := range games {
			if err := s.gameService.DeleteGame(r.Context(), game.ID); err != nil {
				serverError(w, err)
				return
			}
		}

		writeJSON(w, http.StatusNoContent, nil)
	}
}

func (s *Server) linkGames() http.HandlerFunc {
	type Input struct {
		User struct {
			Email    *string `json:"email,omitempty"`
			Username *string `json:"username,omitempty" validate:"required"`
			Age      *uint   `json:"age,omitempty"`
			Password *string `json:"password,omitempty"`
		} `json:"user,omitempty" validate:"required"`
		Game struct {
			Title       *string `json:"title" validate:"required"`
			Description *string `json:"description"`
			URL         *string `json:"url"`
			AgeRating   *uint   `json:"ageRating"`
			Publisher   *string `json:"publisher"`
		} `json:"game,omitempty" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := userFromContext(r.Context())

		if user.IsAnonymous() {
			invalidAuthTokenError(w)
			return
		}

		input := Input{}

		if err := readJSON(r.Body, &input); err != nil {
			badRequestError(w)
			return
		}

		if err := validate.Struct(input.User); err != nil {
			validationError(w, err)
			return
		}
		if err := validate.Struct(input.Game); err != nil {
			validationError(w, err)
			return
		}

		// 1) Getting userID and gameID
		filterUser := models.UserFilter{
			Email:    input.User.Email,
			Username: input.User.Username,
			Age:      input.User.Age,
		}
		users, err := s.userService.Users(r.Context(), filterUser)
		if err != nil {
			serverError(w, err)
			return
		}
		filterGame := models.GameFilter{
			Title:       input.Game.Title,
			Description: input.Game.Description,
			URL:         input.Game.URL,
			AgeRating:   input.Game.AgeRating,
			Publisher:   input.Game.Publisher,
		}
		games, err := s.gameService.Games(r.Context(), filterGame)
		if err != nil {
			serverError(w, err)
			return
		}
		if len(users) > 0 && len(games) > 0 {
			user := users[0]
			game := games[0]
			if user.Age >= game.AgeRating {
				// 2) Create Metadata
				md := &models.Metadata{
					PlayerID:     user.ID,
					Player:       user,
					PlayedGameID: game.ID,
					PlayedGame:   game,
					PlayTime:     0, // initialize playtime at 0
				}
				if err = s.metadataService.CreateMetadata(r.Context(), md); err != nil {
					serverError(w, err)
					return
				}
				writeJSON(w, http.StatusNoContent, nil)
			} else {
				invalidUserAgeError(w)
				return
			}
		} else {
			badRequestError(w)
			return
		}
	}
}
