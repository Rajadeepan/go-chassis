package egress

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
)

// Init initialize router config
func Init() error {
	return nil
}

// ValidateRule validate the route rules of each service
func ValidateRule(rules map[string][]*model.RouteRule) bool {
	for name, rule := range rules {

		for _, route := range rule {
			allWeight := 0
			for _, routeTag := range route.Routes {
				allWeight += routeTag.Weight
			}

			if allWeight > 100 {
				lager.Logger.Warnf("route rule for [%s] is not valid: ruleTag weight is over 100%", name)
				return false
			}
		}

	}
	return true
}
