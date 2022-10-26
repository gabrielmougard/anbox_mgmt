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
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"anbox_mgmt/pkg/models"
	"anbox_mgmt/pkg/postgresql"

	"github.com/gorilla/mux"
)

type Server struct {
	server          *http.Server
	router          *mux.Router
	userService     models.UserService
	gameService     models.GameService
	metadataService models.MetadataService
}

func NewServer(db *postgresql.DB) *Server {
	s := Server{
		server: &http.Server{
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
		router: mux.NewRouter().StrictSlash(true),
	}

	s.routes()

	s.userService = postgresql.NewUserService(db)
	s.gameService = postgresql.NewGameService(db)
	s.metadataService = postgresql.NewMetadataService(db)
	s.server.Handler = s.router

	return &s
}

func (s *Server) Run(port string, gameTrafficFrequency time.Duration, gameTrafficLimitPlayTimePerFreq int32) error {
	sigHandler := make(chan struct{})
	terminate := func() {
		sigHandler <- struct{}{}
	}

	defer terminate()

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	s.server.Addr = port
	log.Printf("server starting on %s", port)

	// generate artificial game traffic
	go generateGameTraffic(gameTrafficFrequency, gameTrafficLimitPlayTimePerFreq, s.metadataService, sigHandler)

	return s.server.ListenAndServe()
}

func healthCheck() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := M{
			"status":  "available",
			"message": "healthy",
			"data":    M{"hello": "beautiful"},
		}
		writeJSON(rw, http.StatusOK, resp)
	})
}

func generateGameTraffic(gameTrafficFrequency time.Duration, gameTrafficLimitPlayTimePerFreq int32, metadataService models.MetadataService, sigHandler chan struct{}) {
	ticker := time.NewTicker(gameTrafficFrequency)
	ctx := context.Background()
	mdFilter := models.MetadataFilter{} // This simulator assumes that there has been traffic for everyone.
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			allMetadata, err := metadataService.Metadata(ctx, mdFilter)
			if err != nil {
				ticker.Stop()
				log.Printf("error fetching gaming metadata: %s", err)
				return
			}
			for _, md := range allMetadata {
				randomPlayTime := md.PlayTime + uint(rand.Int31n(gameTrafficLimitPlayTimePerFreq))
				randomPatch := models.MetadataPatch{
					PlayTime: &randomPlayTime,
				}
				// Every `gameTrafficFrequency` minute, all the game metadata (play time
				// for a player/game pair) is incremented by a random limited number of minutes
				// thus simulating traffic.
				metadataService.UpdateMetadata(ctx, md, randomPatch)
			}
			elapsed := time.Since(start)
			log.Printf("Gaming metadata updated in %s", elapsed)
		case <-sigHandler:
			ticker.Stop()
			log.Print("game traffic generator stopped.")
			return
		}
	}
}
