// slb_controller project main.go
package main

import (
	"fmt"

	slb "github.com/slb-controller/utils"
)

const (
	basepath = "http://slb-agent:9000/polycube/v1"
	slb_name = "my-slb"
)

var (
	slbAPI *slb.SlbApiService
)

func main() {
	cfg := slb.Configuration{BasePath: basepath}
	slbAPI := slb.NewAPIClient(&cfg).SlbApi

	_, err := slbAPI.CreateSlbByID(nil, slb_name, slb.Slb{Name: slb_name})
	if err != nil {
		fmt.Printf("failed to create %s\n", slb_name)
		return
	}
	fmt.Printf("created slb %s successfully\n", slb_name)
}
