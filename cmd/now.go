// Copyright © 2019 Prathik Raj <prathik011@gmail.com>
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

package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prathik/spacedrepetition/service"
	"github.com/spf13/cobra"
	"strings"
)

// todayCmd represents the today command
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, dbErr := sql.Open("sqlite3", "./sr.db")
		if dbErr != nil {
			panic(dbErr)
		}

		spacedRepetition := service.SpacedRepetition{
			SqlDataBase: db,
		}

		spacedRepetition.Init()

		loadTopic:
			topic := spacedRepetition.GetTopicNow()

		if topic == nil {
			fmt.Println("Nothing to recall now. Add more topics.")
			return
		}

		fmt.Println(topic.Title)

		fmt.Print("Did you recall this item? [Y/n]: ")
		var input string
		_, _ = fmt.Scanln(&input)

		if strings.ToLower(input) == "y" || input == "" {
			fmt.Println("Great! Let's recall this later!")
			spacedRepetition.RescheduleTopic(topic)
		} else {
			// TODO: Add ask remove question.
			fmt.Print("Do you want to review and recall after an hour " +
				"[No deletes the topic]? [Y/n]: ")

			// Clear previous input.
			input = ""

			// Get user input for the current question.
			_, _ = fmt.Scanln(&input)

			if strings.ToLower(input) == "y" || input == "" {
				fmt.Println("Do re-review! Let's recall this later!")
				spacedRepetition.RescheduleTopicOneHour(topic)
			}
		}

		goto loadTopic
	},
}

func init() {
	rootCmd.AddCommand(nowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// todayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
