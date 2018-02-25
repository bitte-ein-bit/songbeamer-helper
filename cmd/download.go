package cmd

import (
	"log"

	"github.com/bitte-ein-bit/songbeamer-helper/aws"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
	"github.com/spf13/cobra"
)

func init() {
	//RootCmd.AddCommand(cmdDownload)
}

var cmdDownload = &cobra.Command{
	Use:   "download",
	Short: "Download from AWS",
	Long:  `Download changes from AWS to local.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Starting download of files")
		err := aws.GetS3Files()
		util.CheckForError(err)
	},
}
