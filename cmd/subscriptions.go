/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flynshue/rhsm-cli/pkg/rhsm"
	"github.com/spf13/cobra"
)

// subscriptionsCmd represents the subscriptions command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("subscriptions called")
		return subscriptionList()
	},
}

func subscriptionList() error {
	return rhsmAPI().Call("subscriptionList", nil, nil)
}

type SubscriptionListSucess struct {
	Body []SubscriptionBody
}

type SubscriptionBody struct {
	SubscriptionNumber string  `json:"subscriptionNumber"`
	Name               string  `json:"subscriptionName"`
	Status             string  `json:"status"`
	Sku                string  `json:"sku"`
	Pools              []Pools `json:"pools"`
}

type Pools struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Consumed int    `json:"consumed"`
}

func subscriptionListSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	subscriptions := &SubscriptionListSucess{}
	if err := json.Unmarshal(b, subscriptions); err != nil {
		return err
	}
	fmt.Println("Name, Subscription Number, SKU, Status, Pool ID, Quantity, Consumed")
	for _, sub := range subscriptions.Body {
		subscriptionMeta := fmt.Sprintf("%s, %s, %s, %s", sub.Name, sub.SubscriptionNumber, sub.Sku, sub.Status)
		if len(sub.Pools) != 0 {
			poolMeta := fmt.Sprintf(" %s, %d, %d", sub.Pools[0].ID, sub.Pools[0].Quantity, sub.Pools[0].Consumed)
			subscriptionMeta = subscriptionMeta + "," + poolMeta
		}
		fmt.Printf("%s\n", subscriptionMeta)
	}
	return nil
}

func subscriptionListResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, subscriptionListSuccess)
	return rhsm.NewRestResource("GET", "/subscriptions", router)
}

func init() {
	listCmd.AddCommand(subscriptionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscriptionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscriptionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
