/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/kubefirst/kubefirst-watcher/cmd/internal/logutils"
	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	"github.com/kubefirst/kubefirst-watcher/pkg/k1/informer"
	"github.com/kubefirst/kubefirst-watcher/pkg/k1/k8s"
	"github.com/spf13/cobra"
)

var configFile string
var ownerFile string
var crdAPIVersion string
var crdNamespace string
var crdName string
var crdResource string

// watcherCmd represents the watcher command
var watcherCmd = &cobra.Command{
	Use:   "watcher",
	Short: "Observe k8s events and wait conditions to be satisfied",
	Long:  `TBD`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("watcher called")
		logutils.InitializeLogger()
		logutils.Logger.Info("Watcher ready to start")
		clientSet := k8s.GetK8SConfig()
		myClient := &crd.CRDClient{
			Logger: logutils.Logger,
			CRD: &crd.CRDConfig{
				APIVersion:   crdAPIVersion,
				Namespace:    crdNamespace,
				InstanceName: crdName,
				Resource:     crdResource,
			},
		}
		return informer.StartCRDWatcher(clientSet, myClient, logutils.Logger)
		//return informer.StartWatcher(configFile, ownerFile, logutils.Logger)
	},
}

func init() {
	rootCmd.AddCommand(watcherCmd)

	watcherCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "Provide a yaml with watcher settings")
	watcherCmd.Flags().StringVarP(&ownerFile, "owner-file", "o", "", "Provide a yaml with CRD owner refernece")
	watcherCmd.Flags().StringVar(&crdAPIVersion, "crd-api-version", "k1.kubefirst.io/v1beta1", `CRD API Version.`)
	watcherCmd.Flags().StringVar(&crdNamespace, "crd-namespace", "default", `CRD Namespace.`)
	watcherCmd.Flags().StringVar(&crdName, "crd-instance", "", `CRD instance name. Mandatory in CRD mode.`)
	watcherCmd.Flags().StringVar(&crdResource, "crd-resource", "watchers", `CRD Resource name`)
}
