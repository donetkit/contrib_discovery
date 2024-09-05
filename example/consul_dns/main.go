package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// 通过Consul DNS解析服务
	// 假设Consul运行在本地或配置了DNS解析
	// "my-service.service.consul" 是在Consul中注册的服务名
	service := "my-service.service.consul"

	// 查找服务的IP地址
	ips, err := net.LookupHost(service)
	if err != nil {
		log.Fatalf("无法通过DNS发现服务: %v", err)
	}

	// 输出发现的IP地址
	for _, ip := range ips {
		fmt.Printf("发现服务实例的IP地址: %s\n", ip)
	}
}
