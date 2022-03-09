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
	"strconv"

	"github.com/flynshue/rhsm-cli/pkg/rhsm"
	"github.com/spf13/cobra"
)

// flag vars
var (
	filter   string
	systemID string
)

// systemsCmd represents the systems command
var systemsCmd = &cobra.Command{
	Use:   "systems",
	Short: "list systems associated with redhat account",
	Long: `
# List all systems associated with account
rhsm-cli list systems

# List systems matching filter
rhsm-cli list systems --filter ocp`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case len(filter) > 0:
			return systemsFilter(filter)
		case len(systemID) > 0:
			return systemShow(systemID)
		default:
			return systemsList()
		}
	},
}

var offset, count int

type SystemsListResponse struct {
	Body       []SystemsListBody `json:"body"`
	Pagination `json:"pagination"`
}

type Pagination struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
}

type SystemsListBody struct {
	EntitlementCount int    `json:"entitlementCount"`
	Hostname         string `json:"name"`
	LastCheckin      string `json:"lastCheckin"`
	UUID             string `json:"uuid"`
	Type             string `json:"type"`
}

type SystemShowResponse struct {
	Body SystemShowBody `json:"body"`
}

type SystemShowBody struct {
	ID           string               `json:"uuid"`
	Name         string               `json:"name"`
	Entitlements EntitlementsAttached `json:"entitlementsAttached"`
}

type EntitlementsAttached struct {
	Value []EntitlementValues `json:"value"`
}

type EntitlementValues struct {
	EntitlementID    string `json:"id"`
	SubscriptionName string `json:"subscriptionName"`
	Sku              string `json:"sku"`
}

func systemsFilter(keyword string) error {
	fmt.Println("Hostname, Type, Entitlement Status, Entitlement Count, UUID, Last Checkin")
	params := map[string]string{"filter": keyword, "offset": strconv.Itoa(offset)}
	if err := rhsmAPI().Call("systemsFilter", params, nil); err != nil {
		return err
	}
	return paginationHelper("systemsFilter", params, nil)
}

func systemsList() error {
	params := map[string]string{"offset": strconv.Itoa(offset)}
	if err := rhsmAPI().Call("systemsList", params, nil); err != nil {
		return err
	}
	return paginationHelper("systemsList", params, nil)
}

func systemShow(uuid string) error {
	params := map[string]string{"uuid": uuid}
	return rhsmAPI().Call("systemShow", params, nil)
}

func systemsSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	systems := &SystemsListResponse{}
	if err := json.Unmarshal(b, systems); err != nil {
		return err
	}
	count = systems.Count
	for _, system := range systems.Body {
		fmt.Printf("%s, %s, %d, %s, %s\n",
			system.Hostname, system.Type,
			system.EntitlementCount, system.UUID, system.LastCheckin)
	}
	return nil
}

func systemShowSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	system := &SystemShowResponse{}
	if err := json.Unmarshal(b, system); err != nil {
		return err
	}
	fmt.Println("Hostname, UUID, Subscription Name, Sku, Entitlement ID")
	var subscription, sku, entitlement string
	if len(system.Body.Entitlements.Value) != 0 {
		subscription = system.Body.Entitlements.Value[0].SubscriptionName
		sku = system.Body.Entitlements.Value[0].Sku
		entitlement = system.Body.Entitlements.Value[0].EntitlementID
	}
	fmt.Printf("%s, %s, %s, %s, %s\n", system.Body.Name, system.Body.ID, subscription, sku, entitlement)
	return nil
}

func systemsFilterResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, systemsSuccess)
	return rhsm.NewRestResource("GET", "/systems?offset={{ .offset }}&filter={{ .filter }}", router)
}

func systemsListResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, systemsSuccess)
	return rhsm.NewRestResource("GET", "/systems?offset={{ .offset }}", router)
}

func systemsShowResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, systemShowSuccess)
	return rhsm.NewRestResource("GET", "/systems/{{ .uuid }}?include=entitlements", router)
}

func paginationHelper(resource string, params map[string]string, body interface{}) error {
	for count == 100 {
		offset += 100
		params["offset"] = strconv.Itoa(offset)
		if err := rhsmAPI().Call("systemsFilter", params, body); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	listCmd.AddCommand(systemsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// systemsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// systemsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	systemsCmd.Flags().StringVar(&filter, "filter", "", "filter systems by system name")
	systemsCmd.Flags().StringVar(&systemID, "systemID", "", "get system by system uuid and its entitlements")
}
