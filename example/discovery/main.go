package main

import (
	"github.com/donetkit/contrib_discovery/consul"
	"github.com/donetkit/contrib_discovery/discovery"
)

const GatewayConsulTagsKey = "app/gateway/consul/tags"
const GatewayConsulTagsHeaderKey = "app/gateway/consul/tags/header"
const GatewayConsulTagsTCPKey = "service/gateway/consul/tcp"

func main() {
	clientR, err := consul.New(discovery.WithRegisterAddr("192.168.5.110"), discovery.WithRegisterPort(18500))
	clientW, err := consul.New(discovery.WithRegisterAddr("192.168.5.111"), discovery.WithRegisterPort(8500))
	if err != nil {
		return
	}

	buff, err := clientR.Get(GatewayConsulTagsKey)
	if err == nil {
		if len(buff) > 0 {
			clientW.Set(GatewayConsulTagsKey, string(buff))
		}
	}

	buff, err = clientR.Get(GatewayConsulTagsHeaderKey)
	if err == nil {
		if len(buff) > 0 {
			clientW.Set(GatewayConsulTagsHeaderKey, string(buff))
		}
	}

	buff, err = clientR.Get(GatewayConsulTagsTCPKey)
	if err == nil {
		if len(buff) > 0 {
			clientW.Set(GatewayConsulTagsTCPKey, string(buff))
		}
	}

}
