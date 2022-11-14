/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/6za/k1-watcher/cmd/internal/logutils"
	"github.com/6za/k1-watcher/pkg/k1/informer"
	"github.com/spf13/cobra"
)

var configFile string

// watcherCmd represents the watcher command
var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Observe k8s events and wait conditions to be satisfied",
	Long:  `TBD`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("watcher called")
		logutils.InitializeLogger()
		logutils.Logger.Info("Hello World")
		informer.StartWatcher(logutils.Logger)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)

	watcherCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "Provide a yaml witth watcher settings")
}
