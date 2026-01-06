package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initService() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("service_analyze",
				mcp.WithDescription("filter unhealthy services and analyze it"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to analyze services in"),
					mcp.Required(),
				),
			),
			Handler: s.serviceAnalyze,
		},
	}
}

func (s *Server) serviceAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	res, err := s.k8s.AnalyzeService(ctx, ns)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze services in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
