/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"msn/cmd/data"

	"github.com/spf13/cobra"
)

// populateCmd represents the populate command
var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate DB with random data",
	Long:  `Populate DB with random data`,
	Run: func(cmd *cobra.Command, args []string) {
		database := "sqlite"
		if len(args) == 1 && args[0] == "postgres" {
			database = "postgres"
		}
		data.Populate(database)
	},
}

func init() {
	rootCmd.AddCommand(populateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// populateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// populateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
