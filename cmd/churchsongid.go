package cmd

import (
	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cmdChurchSongID)
}

var cmdChurchSongID = &cobra.Command{
	Use:   "churchsongid",
	Short: "Set ChurchSongID in all songs",
	Long:  `Set ChurchSongID in all songs.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("No longer supported")
	},
}
