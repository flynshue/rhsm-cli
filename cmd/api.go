package cmd

import (
	"github.com/flynshue/rhsm-cli/pkg/rhsm"
	"github.com/spf13/viper"
)

var api *rhsm.API

func rhsmAPI() *rhsm.API {
	if api == nil {
		api = rhsm.NewAPI("https://api.access.redhat.com/management/v1", viper.GetString("token"))
		api.AddResource("systemsFilter", systemsFilterResource())
		api.AddResource("systemsList", systemsListResource())
		api.AddResource("subscriptionList", subscriptionListResource())
		api.AddResource("systemShow", systemsShowResource())
	}
	return api
}
