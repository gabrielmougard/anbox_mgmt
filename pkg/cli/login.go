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
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a user account",
	Long:  `Login to a user account to perform operations in Anbox cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		loginUser := LoginUser{}
		email, _ := cmd.Flags().GetString("email")
		if len(email) > 0 {
			loginUser.Email = email
		} else {
			fmt.Println("--email is a mandatory flag")
			return
		}
		password, _ := cmd.Flags().GetString("password")
		if len(password) > 0 {
			loginUser.Password = password
		} else {
			fmt.Println("--password is a mandatory flag")
			return
		}

		payload := struct {
			User LoginUser `json:"user"`
		}{
			loginUser,
		}
		apiCallPayload("POST", "users/login", payload, SAVE_TOKEN)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "Email of a user")
	loginCmd.Flags().StringP("password", "p", "", "Password of a user")
}
