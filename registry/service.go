package registry

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Nodes   []*Node `json:"nodes"`
}

type Node struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

func (srv *Service) Run(server *grpc.Server) error {
	//自动注册
	rgt := GetCurrentRegistry()
	rgt.Register(context.Background(), srv)

	if srv.Name == "" {
		return fmt.Errorf("service name can not be empty")
	}
	if len(srv.Nodes) == 0 {
		srv.Nodes = []*Node{{"127.0.0.1", 8080, 0}}
	}
	if len(srv.Nodes) != 1 {
		return fmt.Errorf("you can not give more nodes than one when you are running a service")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Nodes[0].Port))
	if err != nil {
		return err
	}

	err = server.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

type OptionService func(service *Service)

func NewService(options ...OptionService) *Service {
	service := &Service{}
	for _, opt := range options {
		opt(service)
	}

	return service
}

func WithName(name string) OptionService {
	return func(service *Service) {
		service.Name = name
	}
}

func WithVersion(version string) OptionService {
	return func(service *Service) {
		service.Version = version
	}
}

func WithNode(ip string, port int) OptionService {
	node := &Node{
		Ip:   ip,
		Port: port,
	}
	return func(service *Service) {
		service.Nodes = append(service.Nodes, node)
	}
}

func WithNodeWeight(ip string, port int, weight int) OptionService {
	node := &Node{
		Ip:     ip,
		Port:   port,
		Weight: weight,
	}
	return func(service *Service) {
		service.Nodes = append(service.Nodes, node)
	}
}
