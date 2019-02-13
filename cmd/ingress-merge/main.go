package main

import (
	"context"
	goflag "flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"github.com/jakubkulhan/ingress-merge"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Merge Ingress Controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err        error
				controller = ingress_merge.NewController()
			)

			if controller.MasterURL, err = cmd.Flags().GetString("apiserver"); err != nil {
				return err
			}

			if controller.KubeconfigPath, err = cmd.Flags().GetString("kubeconfig"); err != nil {
				return err
			}

			if controller.IngressClass, err = cmd.Flags().GetString("ingress-class"); err != nil {
				return err
			}

			if controller.IngressSelector, err = cmd.Flags().GetString("ingress-selector"); err != nil {
				return err
			}

			if controller.ConfigMapSelector, err = cmd.Flags().GetString("configmap-selector"); err != nil {
				return err
			}

			if controller.IngressWatchIgnore, err = cmd.Flags().GetStringArray("ingress-watch-ignore"); err != nil {
				return err
			}

			if controller.ConfigMapWatchIgnore, err = cmd.Flags().GetStringArray("configmap-watch-ignore"); err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			interrupts := make(chan os.Signal, 1)
			go func() {
				select {
				case <-interrupts:
					cancel()
				}
			}()
			signal.Notify(interrupts, syscall.SIGINT, syscall.SIGTERM)

			glog.Infoln("Starting controller")

			return controller.Run(ctx)
		},
	}

	rootCmd.PersistentFlags().String(
		"apiserver",
		"",
		"Address of Kubernetes API server.",
	)

	rootCmd.PersistentFlags().String(
		"kubeconfig",
		"",
		"Path to kubeconfig file.",
	)

	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
	goflag.CommandLine.Parse([]string{}) // prevents glog errors

	rootCmd.Flags().String(
		"ingress-class",
		"merge",
		"Process ingress resources with this `kubernetes.io/ingress.class` annotation.",
	)

	rootCmd.Flags().String(
		"ingress-selector",
		"",
		"Process ingress resources with labels matching this selector string.",
	)

	rootCmd.Flags().String(
		"configmap-selector",
		"",
		"Process configmap resources with labels matching this selector string.",
	)

	rootCmd.Flags().StringArray(
		"ingress-watch-ignore",
		[]string{},
		"Ignore ingress resources with matching annotations (can be specified multiple times).",
	)

	rootCmd.Flags().StringArray(
		"configmap-watch-ignore",
		[]string{},
		"Ignore configmap resources with matching annotations (can be specified multiple times).",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
