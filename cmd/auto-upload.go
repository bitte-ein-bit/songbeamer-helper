package cmd

import (
	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(autoUpload)
}

var autoUpload = &cobra.Command{
	Use:   "auto-upload",
	Short: "Push songs to ChurchTools",
	Long:  `Sync changes to ChurchTools`,
	Run: func(cmd *cobra.Command, args []string) {
		uploadToChurchTools(true)
		log.Finalize()
	},
}
