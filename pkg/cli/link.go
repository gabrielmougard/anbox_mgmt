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
var linkCmd = &cobra.Command{
	Use:   "link [ENTITY]",
	Short: "Link entities",
	Long:  `Link entities`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			entity := args[0]
			if entity == "game" {
				game := CreateGame{}
				user := UpdateUser{}
				title, _ := cmd.Flags().GetString("title")
				if len(title) > 0 {
					game.Title = title
				} else {
					fmt.Println("--title is a mandatory flag")
					return
				}
				if desc, _ := cmd.Flags().GetString("desc"); len(desc) > 0 {
					game.Description = desc
				}
				if url, _ := cmd.Flags().GetString("url"); len(url) > 0 {
					game.URL = url
				}
				if ageRating, _ := cmd.Flags().GetInt("age_rating"); ageRating > 0 {
					game.AgeRating = ageRating
				}
				if publisher, _ := cmd.Flags().GetString("publisher"); len(publisher) > 0 {
					game.Publisher = publisher
				}
				username, _ := cmd.Flags().GetString("username")
				if len(username) > 0 {
					user.Username = username
				} else {
					fmt.Println("--username is a mandatory flag")
					return
				}
				email, _ := cmd.Flags().GetString("email")
				if len(email) > 0 {
					user.Email = email
				}
				age, _ := cmd.Flags().GetInt("age")
				if age > 0 {
					user.Age = age
				}
				password, _ := cmd.Flags().GetString("password")
				if len(password) > 0 {
					user.Password = password
				}

				payload := struct {
					Game CreateGame `json:"game"`
					User UpdateUser `json:"user"`
				}{
					game,
					user,
				}

				apiCallPayload("POST", "games/link", payload)
			} else {
				fmt.Println("The link command is only available of 'game' entity.")
			}
		} else {
			fmt.Println("You must provide an entity to link: only 'game' is available.")
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	linkCmd.Flags().StringP("title", "t", "", "Title of a game")
	linkCmd.Flags().StringP("desc", "d", "", "Description of a game")
	linkCmd.Flags().String("url", "", "URL of a game")
	linkCmd.Flags().Int("age_rating", 0, "Age rating of a game")
	linkCmd.Flags().String("publisher", "", "Title of a game")

	linkCmd.Flags().StringP("email", "e", "", "Email of a user")
	linkCmd.Flags().String("username", "", "Username of a user")
	linkCmd.Flags().Int("age", 0, "Age of a user")
	linkCmd.Flags().String("password", "", "Password of a user")
}
