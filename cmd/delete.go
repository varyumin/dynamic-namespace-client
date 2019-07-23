/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
)


// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [NAMESPACE] ",
	Example: "dynamictl delete dynamic-namespace-dkwcc",
	Short: "Delete dynamic namespace",
	Long: "",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, namespace := range args {
			URL_DYNS_DELETE := fmt.Sprintf("%s://%s.%s.%s/api/ns/%s/", Schema, Name, C.Cluster, Domain, namespace)
			client := C.Login()
			req, _ := http.NewRequest("DELETE", URL_DYNS_DELETE, bytes.NewBuffer([]byte{}))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
