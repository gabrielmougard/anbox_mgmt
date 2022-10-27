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
var updateCmd = &cobra.Command{
	Use:   "update [ENTITY]",
	Short: "Update entities",
	Long:  `Update entities`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			entity := args[0]
			if entity == "game" {
				updateGame := UpdateGame{}
				title, _ := cmd.Flags().GetString("title")
				if len(title) > 0 {
					updateGame.Title = title
				}
				if desc, _ := cmd.Flags().GetString("desc"); len(desc) > 0 {
					updateGame.Description = desc
				}
				if url, _ := cmd.Flags().GetString("url"); len(url) > 0 {
					updateGame.URL = url
				}
				if ageRating, _ := cmd.Flags().GetInt("age_rating"); ageRating > 0 {
					updateGame.AgeRating = ageRating
				}
				if publisher, _ := cmd.Flags().GetString("publisher"); len(publisher) > 0 {
					updateGame.Publisher = publisher
				}

				payload := struct {
					Game UpdateGame `json:"game"`
				}{
					updateGame,
				}
				apiCallPayload("PUT", "games", payload)
			} else if entity == "user" {
				updateUser := UpdateUser{}
				username, _ := cmd.Flags().GetString("username")
				if len(username) > 0 {
					updateUser.Username = username
				} else {
					fmt.Println("--username is a mandatory flag")
					return
				}
				email, _ := cmd.Flags().GetString("email")
				if len(email) > 0 {
					updateUser.Email = email
				}
				age, _ := cmd.Flags().GetInt("age")
				if age > 0 {
					updateUser.Age = age
				}

				password, _ := cmd.Flags().GetString("password")
				if len(password) > 0 {
					updateUser.Password = password
				}

				payload := struct {
					User UpdateUser `json:"user"`
				}{
					updateUser,
				}
				apiCallPayload("PUT", "users", payload)
			} else {
				fmt.Println("Entity not recognized")
			}
		} else {
			fmt.Println("You must provide an entity to update: 'game' or 'user' ?")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("title", "t", "", "Title of a game")
	updateCmd.Flags().StringP("desc", "d", "", "Description of a game")
	updateCmd.Flags().String("url", "", "URL of a game")
	updateCmd.Flags().Int("age_rating", 0, "Age rating of a game")
	updateCmd.Flags().String("publisher", "", "Title of a game")

	updateCmd.Flags().StringP("email", "e", "", "Email of a user")
	updateCmd.Flags().String("username", "", "Username of a user")
	updateCmd.Flags().Int("age", 0, "Age of a user")
	updateCmd.Flags().String("password", "", "Password of a user")
}
