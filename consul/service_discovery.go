package consul

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (s *Client) Register() error {
	check := &consulApi.AgentServiceCheck{
		Timeout:                        fmt.Sprintf("%ds", s.options.TimeOut),        // 超时时间
		Interval:                       fmt.Sprintf("%ds", s.options.IntervalTime),   // 健康检查间隔
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", s.options.DeregisterTime), //check失败后多少秒删除本服务，注销时间，相当于过期时间
	}
	switch s.options.CheckType {
	case "HTTP":
		check.HTTP = s.options.CheckPath
	case "TCP":
		s.options.CheckPath = fmt.Sprintf("%s:%d", s.options.CheckAddr, s.options.CheckPort)
		check.TCP = s.options.CheckPath
	case "GRPC":
		check.GRPC = fmt.Sprintf("%s/%s", s.options.CheckPath, s.options.Name)

		// 2. 注册健康检查服务
		healthServer := NewServer(s.options.Name, s)
		grpc_health_v1.RegisterHealthServer(s.options.GrpcService, healthServer)
		// 设置服务的健康状态为SERVING（健康）
		healthServer.SetServingStatus(s.options.Name, grpc_health_v1.HealthCheckResponse_SERVING)
	}

	svcReg := &consulApi.AgentServiceRegistration{
		ID:                s.options.Id,
		Name:              s.options.Name,
		Tags:              s.options.Tags,
		Port:              s.options.CheckPort,
		Address:           s.options.CheckAddr,
		EnableTagOverride: true,
		Check:             check,
		Checks:            nil,
	}
	err := s.client.Agent().ServiceRegister(svcReg)
	if err != nil {
		return errors.Wrap(err, "register service error")
	}
	return nil
}

func (s *Client) Deregister() error {
	var err error
	if s.options.Nodes > 1 {
		catalogServices, _, errService := s.client.Catalog().Service(s.options.Name, "", nil)
		if errService == nil {
			for _, service := range catalogServices {
				if service.ServiceID == s.options.Id {
					address := fmt.Sprintf("%s:%d", service.Address, 8500)
					addrVal, ok := s.options.NodeAddr[service.Node]
					if ok {
						address = addrVal
					}
					client, errClient := consulApi.NewClient(&consulApi.Config{Token: s.options.Token, Address: address})
					if errClient == nil {
						errServiceDeregister := client.Agent().ServiceDeregister(service.ServiceID)
						if errServiceDeregister != nil {
							err = errors.Wrapf(errServiceDeregister, "deregister service error[key=%s]", s.options.Id)
						}
					}
				}
			}
		}
	} else {
		errServiceDeregister := s.client.Agent().ServiceDeregister(s.options.Id)
		if errServiceDeregister != nil {
			err = errors.Wrapf(errServiceDeregister, "deregister service error[key=%s]", s.options.Id)
		}
	}

	return err
}
