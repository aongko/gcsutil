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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"

	"github.com/spf13/cobra"
)

var bucket string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download bucket object",
	Short: "Download file(s) in the bucket",
	Long: `Download file(s) from the specified bucket into current directories.
For filename(s) with "/", directory(-ies) will be created as well.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("need bucket and object name")
		}
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
			return err
		}
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println("Current working directory: " + wd)

		bucket = args[0]
		for _, object := range args[1:] {
			fmt.Printf("Downloading %v:%v\n", bucket, object)

			data, err := read(client, bucket, object)
			if err != nil {
				log.Fatalf("Cannot read object: %v", err)
				return err
			}

			d := strings.Split(object, "/")
			var fullPath string
			if len(d) > 1 {
				dirs := strings.Join(d[:len(d)-1], "/")
				dir := wd + "/" + dirs
				filename := d[len(d)-1]
				fullPath = dir + "/" + filename
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			} else {
				filename := object
				fullPath = wd + "/" + filename
			}

			ioutil.WriteFile(fullPath, data, 0644)
		}
		return nil
	},
}

func read(client *storage.Client, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
