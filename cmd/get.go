package cmd

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [key]...",
	Short: "Get content of a specific key",
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

		err = db.View(func(txn *badger.Txn) error {
			for _, key := range args {
				item, err := txn.Get([]byte(key))
				if err != nil {
					if err == badger.ErrKeyNotFound {
						return fmt.Errorf("key %s not found", key)
					}
					return err
				}

				value, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				fmt.Println(string(value))
			}

			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
