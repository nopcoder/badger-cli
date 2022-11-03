package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys in the database",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		dir, _ := flags.GetString("dir")
		prefix, _ := flags.GetString("prefix")
		limit, _ := flags.GetInt("limit")
		offset, _ := flags.GetInt("offset")
		opts := badger.DefaultOptions(dir).WithReadOnly(true).WithLogger(nil)
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		w := tabwriter.NewWriter(os.Stdout, 0, 80, 1, ' ', tabwriter.Debug)
		defer w.Flush()
		_, _ = fmt.Fprintf(w, "Key\tSize\tVersion\tMeta\n")

		err = db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false

			if limit > 0 {
				opts.PrefetchSize = limit
			}
			if prefix != "" {
				opts.Prefix = []byte(prefix)
			}
			it := txn.NewIterator(opts)
			defer it.Close()

			currentOffset := 0
			keys := 0
			for it.Rewind(); it.ValidForPrefix([]byte(prefix)); it.Next() {
				currentOffset++
				if currentOffset < offset {
					continue
				}

				item := it.Item()
				value := string(item.KeyCopy(nil))
				_, _ = fmt.Fprintf(w, "%s\t%d\t%d\t%v\n", value, item.EstimatedSize(), item.Version(), item.UserMeta())
				keys++
				if limit > 0 && keys >= limit {
					break
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
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringP("prefix", "p", "", "Key prefix for the search")
	listCmd.PersistentFlags().IntP("limit", "l", 200, "Number of results to return")
	listCmd.PersistentFlags().IntP("offset", "o", 0, "Offset to start at")
}
