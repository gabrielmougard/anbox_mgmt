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

package main

import (
	"log"

	"anbox_mgmt/pkg/config"
	"anbox_mgmt/pkg/postgresql"
	"anbox_mgmt/pkg/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.EnvConfig()

	db, err := postgresql.Open(cfg.DbURI)
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Run(cfg.Port, cfg.GameTrafficFreq, cfg.GameTrafficLimitPlayTimePerFreq))
}
