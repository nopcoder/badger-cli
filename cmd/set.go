package cmd

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a key and its value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		dir, _ := flags.GetString("dir")
		ttl, _ := flags.GetDuration("ttl")
		opts := badger.DefaultOptions(dir).WithLogger(nil)
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		key := args[0]
		value := args[1]
		err = db.Update(func(txn *badger.Txn) error {
			// with ttl
			if ttl.Seconds() > 0 {
				e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(ttl)
				return txn.SetEntry(e)
			}
			// just set the key/value
			return txn.Set([]byte(key), []byte(value))
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.PersistentFlags().Duration("ttl", 0, "Set ttl for the new key")
}
