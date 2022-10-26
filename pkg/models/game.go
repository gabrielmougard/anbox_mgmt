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
)

type Game struct {
	ID          uint      `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	AgeRating   uint      `json:"ageRating"`
	Publisher   string    `json:"publisher"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type GameFilter struct {
	ID          *uint
	Title       *string
	Description *string
	URL         *string
	AgeRating   *uint
	Publisher   *string

	Limit  int
	Offset int
}

type GamePatch struct {
	Title       *string
	Description *string
	URL         *string
	AgeRating   *uint
	Publisher   *string
}

type GameService interface {
	CreateGame(context.Context, *Game) error
	Games(context.Context, GameFilter) ([]*Game, error)
	UpdateGame(context.Context, *Game, GamePatch) error
	DeleteGame(context.Context, uint) error
}
