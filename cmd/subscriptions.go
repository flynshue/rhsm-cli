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
	Short: "List subscriptions under account",
	Long: `List subscriptions under account
Example:

# List all subscriptions
rhsm-cli list subscriptions

# List systems that are consuming subscription
rhsm-cli list subscriptions --subscription <SUBSCRIPTION NUMBER>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if subscription != "" {
			return subSystemsList(subscription)
		}
		return subscriptionList()
	},
}

var subscription string

func subscriptionList() error {
	return rhsmAPI().Call("subscriptionList", nil, nil)
}

func subSystemsList(subscription string) error {
	params := map[string]string{"subscription": subscription}
	return rhsmAPI().Call("subSystemsList", params, nil)
}

type SubscriptionListSuccess struct {
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

type SubSystemResponse struct {
	Body []SubSystemBody `json:"body"`
}

type SubSystemBody struct {
	Name                 string `json:"systemName"`
	UUID                 string `json:"uuid"`
	EntitlementsConsumed int    `json:"totalEntitlementQuantity"`
}

func subscriptionListSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	subscriptions := &SubscriptionListSuccess{}
	if err := json.Unmarshal(b, subscriptions); err != nil {
		return err
	}
	fmt.Println("Name, Subscription Number, SKU, Status, Pool ID, Quantity, Consumed")
	for _, sub := range subscriptions.Body {
		var (
			poolID       string
			poolQuantity int
			poolConsumed int
		)
		if len(sub.Pools) != 0 {
			poolID, poolQuantity, poolConsumed = sub.Pools[0].ID, sub.Pools[0].Quantity, sub.Pools[0].Consumed
		}
		fmt.Printf("%s, %s, %s, %s, %s, %d, %d\n", sub.Name, sub.SubscriptionNumber,
			sub.Sku, sub.Status, poolID, poolQuantity, poolConsumed)
	}
	return nil
}

func subSystemsSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	systems := &SubSystemResponse{}
	if err := json.Unmarshal(b, systems); err != nil {
		return err
	}
	fmt.Println("System Name, UUID, Entitlements Consumed")
	for _, system := range systems.Body {
		fmt.Printf("%s, %s, %d\n", system.Name, system.UUID, system.EntitlementsConsumed)
	}
	return nil
}

func subscriptionListResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, subscriptionListSuccess)
	return rhsm.NewRestResource("GET", "/subscriptions", router)
}

func subSystemResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, subSystemsSuccess)
	return rhsm.NewRestResource("GET", "/subscriptions/{{ .subscription }}/systems", router)
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
	subscriptionsCmd.Flags().StringVar(&subscription, "subscription", "", "subscription number")
}
