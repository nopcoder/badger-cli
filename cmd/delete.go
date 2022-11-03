package cmd

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [key]...",
	Short: "Delete a key and its contents",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		dir, _ := flags.GetString("dir")
		opts := badger.DefaultOptions(dir).WithLogger(nil)
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		err = db.Update(func(txn *badger.Txn) error {
			for _, key := range args {
				if err := txn.Delete([]byte(key)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
