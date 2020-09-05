package cmd

import (
	"database/sql"
	"github.com/prathik/spacedrepetition/service"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
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
		sr := service.SpacedRepetition{
			SqlDataBase: db,
		}
		sr.Init()

		sr.Add(&service.Topic{
			Title: args[0],
		})
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
