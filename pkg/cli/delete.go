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
	"fmt"

	"github.com/spf13/cobra"
)

// the create command
var deleteCmd = &cobra.Command{
	Use:   "delete [ENTITY]",
	Short: "Delete entities",
	Long:  `Delete entities`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			entity := args[0]
			query := ""
			if entity == "game" {
				if title, _ := cmd.Flags().GetString("title"); len(title) > 0 {
					query = queryBuild(query, "title", title)
				}
				if desc, _ := cmd.Flags().GetString("desc"); len(desc) > 0 {
					query = queryBuild(query, "desc", desc)
				}
				if url, _ := cmd.Flags().GetString("url"); len(url) > 0 {
					query = queryBuild(query, "url", url)
				}
				if ageRating, _ := cmd.Flags().GetInt("age_rating"); ageRating > 0 {
					query = queryBuild(query, "age_rating", fmt.Sprint(ageRating))
				}
				if publisher, _ := cmd.Flags().GetString("publisher"); len(publisher) > 0 {
					query = queryBuild(query, "publisher", publisher)
				}
				apiCall("DELETE", "games", query)
			} else if entity == "user" {
				if email, _ := cmd.Flags().GetString("email"); len(email) > 0 {
					query = queryBuild(query, "email", email)
				}
				if username, _ := cmd.Flags().GetString("username"); len(username) > 0 {
					query = queryBuild(query, "username", username)
				}
				if age, _ := cmd.Flags().GetInt("age"); age > 0 {
					query = queryBuild(query, "age", fmt.Sprint(age))
				}
				apiCall("DELETE", "users", query)
			} else {
				fmt.Println("Entity not recognized")
			}
		} else {
			fmt.Println("You must provide an entity to delete: 'game' or 'user' ?")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("title", "t", "", "Title of a game")
	deleteCmd.Flags().StringP("desc", "d", "", "Description of a game")
	deleteCmd.Flags().String("url", "", "URL of a game")
	deleteCmd.Flags().Int("age_rating", 0, "Age rating of a game")
	deleteCmd.Flags().StringP("publisher", "p", "", "Title of a game")

	deleteCmd.Flags().StringP("email", "e", "", "Email of a user")
	deleteCmd.Flags().String("username", "", "Username of a user")
	deleteCmd.Flags().Int("age", 0, "Age of a user")
}
