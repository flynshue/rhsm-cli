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
	"strings"

	"github.com/flynshue/rhsm-cli/pkg/rhsm"
	"github.com/spf13/cobra"
)

// subscriptionAddCmd represents the subscription command
var subscriptionAddCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Attach subscription pool to system",
	Long: `Attach subscription pool to system.

Example:
rhsm-cli add subscription --systemID <SYSTEM UUID> --poolID <POOL ID>

Note: This endpoint isn't working as described in the docs https://access.redhat.com/management/api/rhsm#/system/attachEntitlement.

Please see https://access.redhat.com/management/systems to attach subscription until the feature is fixed.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(systemID) == 0 || len(pool) == 0 {
			fmt.Println(cmd.Long)
			return fmt.Errorf("must supply system uuid and pool id")
		}
		if err := attachSub(systemID, pool); err != nil {
			fmt.Println(cmd.Long)
			return err
		}
		return nil
	},
}

// subscriptionAddCmd represents the subscription command
var subscriptionDelCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Remove subscription/entitlement ID from system",
	Long: `Remove subscription/entitlement ID from system
Example:

rhsm-cli remove subscription --systemID <SYSTEM UUID> --entitlementID <ENTITLEMENT ID>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(systemID) == 0 || len(entitlementID) == 0 {
			return fmt.Errorf("must supply system UUID and entitlement ID")
		}
		return removeSub(systemID, entitlementID)
	},
}

var (
	entitlementID string
	pool          string
)

type AttachSubResponse struct {
	Body AttachSubBody `json:"body"`
}

type AttachSubBody struct {
	Hostname     string `json:"name"`
	UUID         string `json:"uuid"`
	Entitlements int    `json:"entitlementsAttachedCount"`
	LastCheckin  string `json:"lastCheckin"`
}

func removeSub(system, entitlement string) error {
	params := map[string]string{"uuid": systemID, "entitlementID": entitlementID}
	return rhsmAPI().Call("removeSub", params, nil)
}

func attachSub(system, pool string) error {
	params := map[string]string{"uuid": system, "pool": pool}
	return rhsmAPI().Call("attachSub", params, nil)
}

func attachSubSuccess(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	system := &AttachSubResponse{}
	if err := json.Unmarshal(b, system); err != nil {
		return err
	}
	fmt.Printf("%s, %s, %d, %s\n", system.Body.Hostname, system.Body.UUID, system.Body.Entitlements, system.Body.LastCheckin)
	return nil
}

func removeSubSuccess(resp *http.Response) error {
	requestPath := strings.Split(resp.Request.URL.Path, "/")
	systemID := requestPath[4]
	entitlementID := requestPath[5]
	fmt.Printf("Successfully removed entitlement %s from %s\n", entitlementID, systemID)
	return nil
}

func attachSubResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(200, attachSubSuccess)
	router.DefaultRouter = func(resp *http.Response) error {
		return fmt.Errorf("%d from %s %s", resp.StatusCode, resp.Request.Method, resp.Request.URL.String())
	}
	return rhsm.NewRestResource("POST", "/systems/{{ .uuid }}/entitlements?pool={{ .pool }}&quantity=1", router)
}

func removeSubResource() *rhsm.RestResource {
	router := rhsm.NewRouter()
	router.AddFunc(204, removeSubSuccess)
	return rhsm.NewRestResource("DELETE", "/systems/{{ .uuid }}/{{ .entitlementID }}", router)
}

func init() {
	addCmd.AddCommand(subscriptionAddCmd)
	removeCmd.AddCommand(subscriptionDelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscriptionAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscriptionAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	subscriptionAddCmd.Flags().StringVar(&systemID, "systemID", "", "system UUID")
	subscriptionAddCmd.Flags().StringVar(&pool, "pool", "", "pool ID to attach")
	subscriptionDelCmd.Flags().StringVar(&entitlementID, "entitlementID", "", "entitlement ID to remove")
	subscriptionDelCmd.Flags().StringVar(&systemID, "systemID", "", "system UUID")
}
