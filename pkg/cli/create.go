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
var createCmd = &cobra.Command{
	Use:   "create [ENTITY]",
	Short: "Create entities",
	Long:  `Create entities`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			entity := args[0]
			if entity == "game" {
				createGame := CreateGame{}
				title, _ := cmd.Flags().GetString("title")
				if len(title) > 0 {
					createGame.Title = title
				} else {
					fmt.Println("--title is a mandatory flag")
					return
				}
				if desc, _ := cmd.Flags().GetString("desc"); len(desc) > 0 {
					createGame.Description = desc
				}
				if url, _ := cmd.Flags().GetString("url"); len(url) > 0 {
					createGame.URL = url
				}
				if ageRating, _ := cmd.Flags().GetInt("age_rating"); ageRating > 0 {
					createGame.AgeRating = ageRating
				}
				if publisher, _ := cmd.Flags().GetString("publisher"); len(publisher) > 0 {
					createGame.Publisher = publisher
				}

				payload := struct {
					Game CreateGame `json:"game"`
				}{
					createGame,
				}
				apiCallPayload("POST", "games", payload)
			} else if entity == "user" {
				newUser := NewUser{}
				email, _ := cmd.Flags().GetString("email")
				if len(email) > 0 {
					newUser.Email = email
				} else {
					fmt.Println("--email is a mandatory flag")
					return
				}
				age, _ := cmd.Flags().GetInt("age")
				if age > 0 {
					newUser.Age = age
				} else {
					fmt.Println("--age is a mandatory flag")
					return
				}
				username, _ := cmd.Flags().GetString("username")
				if len(username) > 0 {
					newUser.Username = username
				} else {
					fmt.Println("--username is a mandatory flag")
					return
				}
				password, _ := cmd.Flags().GetString("password")
				if len(password) > 0 {
					newUser.Password = password
				} else {
					fmt.Println("--password is a mandatory flag")
					return
				}

				payload := struct {
					User NewUser `json:"user"`
				}{
					newUser,
				}
				apiCallPayload("POST", "users", payload)
			} else {
				fmt.Println("Entity not recognized")
			}
		} else {
			fmt.Println("You must provide an entity to create: 'game' or 'user' ?")
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("title", "t", "", "Title of a game")
	createCmd.Flags().StringP("desc", "d", "", "Description of a game")
	createCmd.Flags().String("url", "", "URL of a game")
	createCmd.Flags().Int("age_rating", 0, "Age rating of a game")
	createCmd.Flags().String("publisher", "", "Title of a game")

	createCmd.Flags().StringP("email", "e", "", "Email of a user")
	createCmd.Flags().String("username", "", "Username of a user")
	createCmd.Flags().Int("age", 0, "Age of a user")
	createCmd.Flags().String("password", "", "Password of a user")
}
