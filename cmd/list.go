// Copyright © 2017 NAME HERE andrew.ongko@gmail.com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	"github.com/aongko/gcsutil/gcsutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var prefix string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list bucket [bucket...]",
	Short: "list objects in one or multiple buckets",
	Long:  `List object in one or multiple buckets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no bucket specified")
		}
		ctx := context.Background()
		client, err := gcsutil.NewStorageClient(ctx, viper.GetString("service_account_file"))
		if err != nil {
			log.Fatal(err)
			return err
		}
		for _, bucket := range args {
			if prefix == "" {
				list(os.Stdout, client, bucket)
			} else {
				listByPrefix(os.Stdout, client, bucket, prefix, "")
			}
		}
		return nil
	},
}

func list(w io.Writer, client *storage.Client, bucket string) error {
	ctx := context.Background()
	// [START storage_list_files]
	it := client.Bucket(bucket).Objects(ctx, nil)
	fmt.Fprintln(w, "Bucket: "+bucket)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files]
	return nil
}

func listByPrefix(w io.Writer, client *storage.Client, bucket, prefix, delim string) error {
	ctx := context.Background()
	// [START storage_list_files_with_prefix]
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimeter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	fmt.Fprintln(w, "Bucket: "+bucket)
	fmt.Fprintln(w, "Using Prefix: "+prefix)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files_with_prefix]
	return nil
}

func init() {
	listCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "This will list files in a bucket with the specified prefix")
	RootCmd.AddCommand(listCmd)
}
