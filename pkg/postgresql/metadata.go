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

var _ models.MetadataService = (*MetadataService)(nil)

type MetadataService struct {
	db *DB
}

func NewMetadataService(db *DB) *MetadataService {
	return &MetadataService{db}
}

func (ms *MetadataService) CreateMetadata(ctx context.Context, md *models.Metadata) error {
	tx, err := ms.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := createMetadata(ctx, tx, md); err != nil {
		return err
	}

	return tx.Commit()
}

func (ms *MetadataService) Metadata(ctx context.Context, filter models.MetadataFilter) ([]*models.Metadata, error) {
	tx, err := ms.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	md, err := findMetadata(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return md, tx.Commit()
}

func (ms *MetadataService) UpdateMetadata(ctx context.Context, md *models.Metadata, patch models.MetadataPatch) error {
	tx, err := ms.db.BeginTxx(ctx, nil)

	if err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	defer tx.Rollback()

	if err := updateMetadata(ctx, tx, md, patch); err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return models.ErrInternal
	}

	return nil
}

func (ms *MetadataService) DeleteMetadata(ctx context.Context, id uint) error {
	tx, err := ms.db.BeginTxx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = deleteMetadata(ctx, tx, id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func createMetadata(ctx context.Context, tx *sqlx.Tx, md *models.Metadata) error {
	query := `
	INSERT INTO metadata (player_id, played_game_id, play_time) 
	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`

	args := []interface{}{
		md.Player.ID,
		md.PlayedGame.ID,
		md.PlayTime,
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&md.ID, &md.CreatedAt, &md.UpdatedAt)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func deleteMetadata(ctx context.Context, tx *sqlx.Tx, id uint) error {
	query := "DELETE FROM metadata WHERE id = $1"
	return execQuery(ctx, tx, query, id)
}

func updateMetadata(ctx context.Context, tx *sqlx.Tx, md *models.Metadata, patch models.MetadataPatch) error {
	if v := patch.PlayTime; v != nil {
		md.PlayTime = *v
	}

	args := []interface{}{
		md.PlayTime,
		md.ID,
	}

	query := `
	UPDATE metadata
	SET playTime = $1, updated_at = NOW() WHERE id = $2
	RETURNING updated_at`

	if err := tx.QueryRowxContext(ctx, query, args...).Scan(&md.UpdatedAt); err != nil {
		log.Printf("error updating record: %v", err)
		return models.ErrInternal
	}

	return nil
}

func findMetadata(ctx context.Context, tx *sqlx.Tx, filter models.MetadataFilter) ([]*models.Metadata, error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0 // used to set correct postgres argument enums i.e $1, $2

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.PlayerID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("player_id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.PlayedGameID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("played_game_id = $%d", argPosition)), append(args, *v)
	}

	if v := filter.PlayTime; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf("play_time = $%d", argPosition)), append(args, *v)
	}

	query := "SELECT * from metadata" + formatWhereClause(where) + " ORDER BY created_at DESC"
	articles, err := queryMetadata(ctx, tx, query, args...)

	if err != nil {
		return articles, err
	}

	return articles, nil
}

func queryMetadata(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) ([]*models.Metadata, error) {
	md := make([]*models.Metadata, 0)

	if err := findMany(ctx, tx, &md, query, args...); err != nil {
		return md, err
	}

	return md, nil
}
