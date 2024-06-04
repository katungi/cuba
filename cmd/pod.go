package cmd

import (
	 "github.com/spf13/cobra"
)

var podCmd = &cobra.Command{
	Use: "pod",
	Short: "The command line tool to run commands on pods",
}

var createCmd = &cobra.Command{
	Use: "create",
	Short: "Create a pod",
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all pods",
	// Run:
}


func init() {
	// rootCmd := &cobra.Command{} // Define rootCmd
	rootCmd.AddCommand(podCmd)
	podCmd.AddCommand(createCmd)
	podCmd.AddCommand(listCmd)
}

