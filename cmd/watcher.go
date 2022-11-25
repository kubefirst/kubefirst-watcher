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
var ownerFile string

// watcherCmd represents the watcher command
var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Observe k8s events and wait conditions to be satisfied",
	Long:  `TBD`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("watcher called")
		logutils.InitializeLogger()
		logutils.Logger.Info("Watcher ready to start")
		return informer.StartWatcher(configFile, ownerFile, logutils.Logger)
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)

	watcherCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "Provide a yaml with watcher settings")
	watcherCmd.Flags().StringVarP(&ownerFile, "owner-file", "o", "", "Provide a yaml with CRD owner refernece")
}
