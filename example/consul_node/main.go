package main

import (
	"fmt"
	"github.com/donetkit/contrib_discovery/consul"
	"github.com/donetkit/contrib_discovery/discovery"
	"log"
)

const GatewayConsulTagsKey = "app/gateway/consul/tags"

func main() { // 1. 创建Consul客户端
	client, err := consul.New(discovery.WithRegisterAddr("192.168.5.110"), discovery.WithRegisterPort(18500))
	if err != nil {
		log.Fatal(err)
	}

	buff, err := client.Get(GatewayConsulTagsKey)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(buff))

}
