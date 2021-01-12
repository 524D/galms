package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// decoyCmd represents the decoy command
var decoyCmd = &cobra.Command{
	Use:   "decoy",
	Short: "Generate decoy FASTA database",
	Long: `This command generates a decoy database according to
the one of the methods set by the --method flag.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("decoy called")
	},
}

func init() {
	rootCmd.AddCommand(decoyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decoyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decoyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
