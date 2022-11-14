/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/6za/k1-watcher/pkg/k1/informer"
	"github.com/spf13/cobra"
)

// watcherCmd represents the watcher command
var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Observe k8s events and wait conditions to be satisfied",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("watcher called")
		informer.StartWatcher()
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watcherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watcherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
