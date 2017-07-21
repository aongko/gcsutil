// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"path/filepath"

	"github.com/aongko/gcsutil/gcsutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var paths []string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload bucket objects",
	Short: "Upload objects to the bucket",
	Long: `Upload objects to the bucket.
If objects is in a dir, specify the directory name, everything
inside it will be uploaded accordingly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("upload called")
		if len(args) < 1 {
			return errors.New("no bucket specified")
		} else if len(args) < 2 {
			return errors.New("no object name")
		}
		ctx := context.Background()
		client, err := gcsutil.NewStorageClient(ctx, viper.GetString("service_account_file"))
		if err != nil {
			log.Fatal(err)
			return err
		}
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println("Current working directory: " + wd)

		bucket := args[0]
		fmt.Println("Bucket:", bucket)

		object := args[1]
		// if _, err := os.Stat(object); err == nil {
		// 	fmt.Println("file exists")
		// } else {
		// 	fmt.Println("file not exists")
		// }

		err = filepath.Walk(object, visit)
		if err != nil {
			return err
		}
		for _, path := range paths {
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
				return err
			}
			defer f.Close()

			if prefix != "" {
				path = prefix + path
			}
			wc := client.Bucket(bucket).Object(path).NewWriter(ctx)
			if _, err = io.Copy(wc, f); err != nil {
				fmt.Printf("error copying %+v", err)
				return err
			}
			if err := wc.Close(); err != nil {
				fmt.Printf("error closing %+v", err)
				return err
			}
			fmt.Println("done working on:", path)
		}

		return nil
	},
}

func visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !f.Mode().IsDir() {
		paths = append(paths, path)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "The prefix, to simulate directories in the buckets")
}
