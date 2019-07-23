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
  "fmt"
  "os"
  "github.com/spf13/cobra"

  "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"

  "net/http"
  "crypto/tls"
  "net/http/cookiejar"
  "golang.org/x/net/publicsuffix"
  "net/url"
)

var C CredentialServer
var cfgFile string
var URL string

const Schema = "https"
const Name  = "dynamic-ns"
const Domain  = "k8s.tcsbank.ru"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "kubectl-dynamicns",
  Short: "",
  Long: "",
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}
type CredentialServer struct {
  User string
  Password string
  Cluster string
}
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.
  rootCmd.PersistentFlags().StringVarP(&C.User, "user", "U", os.Getenv("LDAP_USER"), "User from LDAP. Environment LDAP_USER")
  rootCmd.PersistentFlags().StringVarP(&C.Password, "password", "P", os.Getenv("LDAP_PASSWORD"), "Password from LDAP. Environment LDAP_PASSWORD")
  rootCmd.PersistentFlags().StringVarP(&C.Cluster, "cluster", "c", os.Getenv("CLUSTER"), "Cluster ENV. Environment CLUSTER")
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dynamictl.yaml)")


  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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

    // Search config in home directory with name ".dynamictl" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".dynamictl")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

func SetHttpClient() (http.Client)  {
  http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
  options := cookiejar.Options{
    PublicSuffixList: publicsuffix.List,
  }
  jar, err := cookiejar.New(&options)
  if err != nil {
    panic(err)
  }
  return http.Client{Jar: jar}
}
func (c *CredentialServer) Login() (client http.Client) {
  URL = fmt.Sprintf("%s://%s.%s.%s", Schema, Name, C.Cluster, Domain)
  client = SetHttpClient()
  resp, err := client.Get(URL)
  if err != nil {
    //panic(err)
  }
  auth_url := resp.Request.URL.String()
  resp, err = client.PostForm(auth_url, url.Values{"login": {C.User}, "password": {C.Password}})
  if err != nil {
     panic("Can not login! Bad login or password")
     os.Exit(1)
  }
  return client
}