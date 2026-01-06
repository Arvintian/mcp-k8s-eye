package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initNode() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("node_analyze",
				mcp.WithDescription("filter unhealthy nodes and analyze it"),
				mcp.WithString("name",
					mcp.Description("the node name to analyze, can be a \"\" string to all nodes"),
				),
			),
			Handler: s.nodeAnalyze,
		},
	}
}
func (s *Server) nodeAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := ctr.GetString("name", "")
	res, err := s.k8s.AnalyzeNode(ctx, name)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze node %s: %v", name, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
