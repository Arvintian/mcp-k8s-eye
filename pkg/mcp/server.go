package mcp

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/k8s"
)

type Server struct {
	server  *server.MCPServer
	k8s     *k8s.Kubernetes
	write   bool
	analyze bool
}

func NewServer(name, version string, write, extend, analyze bool) (*Server, error) {
	s := &Server{
		server: server.NewMCPServer(
			name,
			version,
			server.WithResourceCapabilities(true, true),
			server.WithPromptCapabilities(true),
			server.WithLogging(),
			server.WithRecovery(),
		),
		write:   write,
		analyze: analyze,
	}
	k8s, err := k8s.NewKubernetes()
	if err != nil {
		return nil, err
	}
	s.k8s = k8s

	tools := slices.Concat(
		s.initResource(),
		s.initPod(),
		s.initDeployment(),
		s.initService(),
		s.initStatefulSet(),
		s.initNode(),
		s.initIngress(),
		s.initCronJob(),
	)
	if extend {
		tools = append(tools, s.initNetworkPolicy()...)
		tools = append(tools, s.initWebhook()...)
	}
	s.server.AddTools(tools...)
	for _, item := range tools {
		os.Stderr.WriteString(fmt.Sprintf("add tool %s\n", item.Tool.Name))
	}
	return s, nil
}

func (s *Server) ServeStdio() error {
	options := []server.StdioOption{}
	return server.ServeStdio(s.server, options...)
}

func (s *Server) ServeSSE() *server.SSEServer {
	options := []server.SSEOption{
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(time.Second * 120),
	}
	return server.NewSSEServer(s.server, options...)
}
