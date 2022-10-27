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

package cli

import (
	//"fmt"
	//"os"
	"anbox_mgmt/pkg/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var cfg config.Config

func init() {
	cfg = config.EnvConfig()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anbox-cli",
	Short: "Managing the Anbox application",
	Long:  `Managing the Anbox application. We can do CRUD operations on "users" and "games" and also create links between entities`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func queryBuild(query string, k string, v string) string {
	if query == "" {
		query += "?"
	}
	query += fmt.Sprintf("%s=%s", k, v)
	return query
}

func apiCall(verb string, path string, query string) {
	req, err := http.NewRequest(verb, fmt.Sprintf("http://0.0.0.0:%s/api/v1/%s%s", cfg.Port, path, query), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}

func apiCallPayload(verb string, path string, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(verb, fmt.Sprintf("http://0.0.0.0:%s/api/v1/%s", cfg.Port, path), body)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "<API Key>") // TODO

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}
