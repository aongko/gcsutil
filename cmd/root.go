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
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var googleServiceAccountFile string

// RootCmd is the root command of Cobra
var RootCmd = &cobra.Command{
	Use:   "gcsutil",
	Short: "Some helper for Google Cloud Storage",
	Long:  `Provide some helper to list and downloads object(s) from Google Cloud Storage`,
}

// Execute to run the Cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gcsutil.yaml)")
	RootCmd.PersistentFlags().StringVar(&googleServiceAccountFile, "service_account_file", "", "path to google service account file")
	viper.BindPFlag("service_account_file", RootCmd.PersistentFlags().Lookup("service_account_file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gcsutil" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gcsutil")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if viper.GetString("service_account_file") == "" {
			fmt.Println("Could not read config file:", err)
			fmt.Println("Please provide the file, or set --service_account_file")
			os.Exit(1)
		}
	}
}
