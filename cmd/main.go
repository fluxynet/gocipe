package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Versioning info
var (
	appCommit = "n/a"
	appBuilt  = "n/a"
)

func main() {
	cmdRoot := &cobra.Command{
		Use:   "gocipe",
		Short: "Launch gocipe",
	}

	// Versioning
	cmdRoot.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Software Version information",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Printf("\nCommit : %v\nBuilt: %v\n", appCommit, appBuilt)
			},
		},
	)

	cmdRoot.Execute()
}
