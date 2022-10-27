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

var _ models.GameService = (*GameService)(nil)

type GameService struct {
	db *DB
}

func NewGameService(db *DB) *GameService {
	return &GameService{db}
}

func (gs *GameService) CreateGame(ctx context.Context, game *models.Game) error {
	tx, err := gs.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createGame(ctx, tx, game); err != nil {
		return err
	}

	return tx.Commit()
}

func (gs *GameService) Games(ctx context.Context, filter models.GameFilter) ([]*models.Game, error) {
	tx, err := gs.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	games, err := findGames(ctx, tx, filter)

	if err != nil {
		return nil, err
	}

	return games, tx.Commit()
}

func (gs *GameService) UpdateGame(ctx context.Context, game *models.Game, patch models.GamePatch) error {
	tx, err := gs.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = updateGame(ctx, tx, game, patch)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (gs *GameService) DeleteGame(ctx context.Context, id uint) error {
	tx, err := gs.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = deleteGame(ctx, tx, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func createGame(ctx context.Context, tx *sqlx.Tx, game *models.Game) error {
	query := `
	INSERT INTO games (title, description, url, age_rating, publisher) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	args := []interface{}{
		game.Title,
		game.Description,
		game.URL,
		game.AgeRating,
		game.Publisher,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&game.ID, &game.CreatedAt, &game.UpdatedAt)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func findGames(ctx context.Context, tx *sqlx.Tx, filter models.GameFilter) ([]*models.Game, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0 // used to set correct postgres argument enums i.e $1, $2

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Title; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("title = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Description; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("description = $%d", argPosition)), append(args, *v)
	}

	if v := filter.URL; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("url = $%d", argPosition)), append(args, *v)
	}

	if v := filter.AgeRating; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("age_rating = $%d", argPosition)), append(args, *v)
	}

	if v := filter.Publisher; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("publisher = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from games" + formatWhereClause(where) + " ORDER BY created_at DESC"
	games, err := queryGames(ctx, tx, query, args...)

	if err != nil {
		return games, err
	}

	return games, nil
}

func deleteGame(ctx context.Context, tx *sqlx.Tx, id uint) error {
	query := "DELETE FROM games WHERE id = $1"
	return execQuery(ctx, tx, query, id)
}

func queryGames(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) ([]*models.Game, error) {
	games := make([]*models.Game, 0)
	if err := findMany(ctx, tx, &games, query, args...); err != nil {
		return games, err
	}

	return games, nil
}

func updateGame(ctx context.Context, tx *sqlx.Tx, game *models.Game, patch models.GamePatch) error {
	if v := patch.Title; v != nil {
		game.Title = *v
	}

	if v := patch.Description; v != nil {
		game.Description = *v
	}

	if v := patch.URL; v != nil {
		game.URL = *v
	}

	if v := patch.AgeRating; v != nil {
		game.AgeRating = *v
	}

	if v := patch.Publisher; v != nil {
		game.Publisher = *v
	}

	args := []interface{}{
		game.Title,
		game.Description,
		game.URL,
		game.AgeRating,
		game.Publisher,
		game.ID,
	}

	query := `
	UPDATE games 
	SET title = $1, description = $2, url = $3, age_rating = $4, publisher = $5, updated_at = NOW() WHERE id = $6
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&game.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return models.ErrInternal
	}

	return nil
}
