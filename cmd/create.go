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
	"github.com/prometheus/common/log"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

type InputFlag struct{
	DateDelete string
	LimitCPU float32
	LimitMemory int
	RequestCPU float32
	RequestMemory int
	Users []string
	Groups []string
}

type RequestToCreateDNamespace struct {
	DateDelete time.Time `json:"date_delete"`
	Resources struct{
		Limits struct{
			CPU string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"limits"`
		Requests struct{
			CPU string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"requests"`
	} `json:"resources"`
	Access struct{
		Users []string `json:"users"`
		Groups []string `json:"groups"`
	} `json:"access"`
}

const timeParse  = "2006-01-02T15:04:05Z07:00"
var (
	In InputFlag
)
// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create dynamic namespace",
	Long: "",
	Run: In.CreateNamespace,
}

func stringToTime(string string) time.Time {
	date, err := time.Parse(timeParse, string)
	if err != nil {
		log.Fatalf("Cannot parse time::: %s", err)
	}
	return date
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&In.DateDelete, "time-delete", "d", "","Timespamp delete")
	createCmd.Flags().StringArrayVarP(&In.Users, "users", "u", []string{},"Access users")
	createCmd.Flags().StringArrayVarP(&In.Groups, "groups", "g", []string{},"Access groups")
	createCmd.Flags().Float32VarP(&In.LimitCPU, "limit.cpu", "", 0, "Limnit CPU")
	createCmd.Flags().Float32VarP(&In.RequestCPU, "request.cpu", "", 0, "Request CPU")
	createCmd.Flags().IntVarP(&In.LimitMemory, "limit.memory", "", 0, "Limnit memory")
	createCmd.Flags().IntVarP(&In.RequestMemory, "request.memory", "", 0, "Request memory")
}

func (c *InputFlag) CreateNamespace (cmd *cobra.Command, args []string) {
	URL_DYNS_CREATE := fmt.Sprintf("%s://%s.%s.%s/api/ns/", Schema, Name, C.Cluster, Domain)
	var dateToDelete time.Time
	client := C.Login()
	if In.DateDelete == ""{
		dateToDelete = time.Now().AddDate(0,0, 1).UTC()
	} else {
		dateToDelete = stringToTime(In.DateDelete).UTC()
	}

	if In.RequestCPU <= 0 {
		In.RequestCPU = In.LimitCPU
	}
	if In.RequestMemory <= 0 {
		In.RequestMemory = In.LimitMemory
	}

	Output := RequestToCreateDNamespace{
		DateDelete: dateToDelete ,
		Resources: struct {
			Limits struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			}`json:"limits"`
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
		}{
		Limits: struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		}{
			CPU: fmt.Sprint(In.LimitCPU),
			Memory: fmt.Sprint(In.LimitMemory),
		},
		Requests: struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		}{
			CPU: fmt.Sprint(In.RequestCPU),
			Memory: fmt.Sprint(In.RequestMemory),
		},
	},
	Access: struct {
		Users  []string `json:"users"`
		Groups []string `json:"groups"`
	}{
		Users: In.Users,
		Groups: In.Groups,
		},
	}
	outputJson, _ := json.Marshal(Output)
	req, _ := http.NewRequest("POST", URL_DYNS_CREATE, bytes.NewBuffer(outputJson))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println( string(body))
}