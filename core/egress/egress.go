package egress

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"sync"
	"github.com/ServiceComb/go-chassis/core/config"
	"fmt"
	"regexp"
)

var lock sync.RWMutex

//FetchRouteRule return all rules
func  FetchEgressRule() map[string][]*model.EgressRule {
	return GetEgressRule()
}

// GetRouteRule get route rule
func GetEgressRule() map[string][]*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return config.EgressDefinition.Destinations;
}


func hostallslices(){
	EgressRules := GetEgressRule()
	for key, EgressRule := range EgressRules {
		fmt.Println("Raj Values",key,EgressRule)
	}
}

//Check Egress rule matches
func Match(hostname string) (bool, *model.EgressRule){
	EgressRules := GetEgressRule()
	for key, egressRules := range EgressRules {
		fmt.Println("Raj Values",key,egressRules)
		for _, egress := range  egressRules {
			fmt.Println("Raj egress host ", egress.Host)
			fmt.Println("Raj egress ports ", egress.Ports)
			for _, host := range egress.Host {
				if len(host) > 1 {
					if string(host[0]) != "*" {
						fmt.Println("The first character is not *", host)
						if host == hostname {
							fmt.Println("Raj : the value of host from configuration and from the request", host, hostname)
							fmt.Println("1. MATCHED")
							return true, egress
						}
					} else if string(host[0]) == "*" {
						fmt.Println("THe first character is *")
						substring := host[1:]
						fmt.Println("llllll", substring)
						fmt.Println("iiii", substring+"$")
						match, err := regexp.MatchString(substring+"$", hostname);
						fmt.Println("oooooooo", match, err)
						if match == true {
							fmt.Println("2. MATCHED")
							return true, egress
						}
					}
				}
			}
		}
	}

	return false, nil
}

//Check Egress rule matches
func Match1() (map[string]*model.EgressRule, map[string]*model.EgressRule){

	searchstructure := make(map[string]*model.EgressRule)
	regexstructure := make(map[string]*model.EgressRule)

	EgressRules := GetEgressRule()
	for key, egressRules := range EgressRules {
		fmt.Println("Raj Values",key,egressRules)
		for _, egress := range  egressRules{
			fmt.Println("Raj egress host ", egress.Host)
			fmt.Println("Raj egress ports ", egress.Ports)
			for _, host := range egress.Host{
				if len(host) > 1 && string(host[0]) != "*" {
					fmt.Println("Raj: The value of string 0", string(host[0]))
					searchstructure[host] = egress
				}else if string(host[0]) == "*"{
					substring := host[1:]
					regexstructure[substring] = egress
				}
			}
		}
	}
    fmt.Println("Raj: ", searchstructure)
	return searchstructure, regexstructure
}

func SearchinMap(hostname string)(bool, *model.EgressRule){
	searchstruct, regexstructure := Match1()

	if val, ok := searchstruct[hostname]; ok{
		fmt.Println("3. MATCHED")
		return true, val
	}

	for key, value := range regexstructure {
			match, _ := regexp.MatchString(key+"$", hostname);
			if match == true {
				fmt.Println("3. MATCHED")
				return true, value

		}
	}
	return false, nil
}