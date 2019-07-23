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
	"time"
	"github.com/spf13/cobra"
	"text/template"
	"os"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"text/tabwriter"
)


type InfoDnamespace  struct {
	Name string `json:"name"`
	DateCreate time.Time `json:"date_create"`
	DateDelete time.Time `json:"date_delete"`
	Owner string `json:"owner"`
}
var listText = `
NAMESPACE	OWNER	CREATE	DELETE	
{{ range . }}
{{ .Name }}	{{ .Owner }}	{{ .DateCreate }}	{{ .DateDelete }}	
{{ end }}
	`
var list []InfoDnamespace
// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your dynamic namespaces",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		URL_DYNS_LIST := fmt.Sprintf("%s://%s.%s.%s/api/ns/", Schema, Name, C.Cluster, Domain)
		client := C.Login()
		req, _ := http.NewRequest("GET", URL_DYNS_LIST, bytes.NewBuffer([]byte{}))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&list)
		if err != nil {
			fmt.Println(err)
		}

		for i, namespace := range list{
			if namespace.Owner != C.User{
				list = append(list[:i], list[i+1:]...)
			}
		}

		tmpl, _ := template.New("list").Parse(listText)

		w := tabwriter.NewWriter(os.Stdout,  32, 0, 0, ' ', 0)
		if err := tmpl.Execute(w, list); err != nil {
			fmt.Println(err)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
