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
		log.Infof("Auto Upload Funktion ist noch deaktivert. Verwende bei Bedarf ct-upload.")
	},
}
