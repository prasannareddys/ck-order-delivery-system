package main

import (
	"log"

	"github.com/Propertyfinder/ck-order-delivery-system/pkg/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ck-order-delivery-system",
		Short: "Entry point command for the application",
	}

	// start order process command
	rootCmd.AddCommand(
		cmd.NewServerCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
