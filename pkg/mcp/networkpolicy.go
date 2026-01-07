package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/common"
)

// Register networkpolicy analyze tool
func (s *Server) initNetworkPolicy() []server.ServerTool {
	tools := []server.ServerTool{}
	if s.analyze {
		tools = append(tools, []server.ServerTool{{
			Tool: mcp.NewTool("networkpolicy_analyze",
				mcp.WithDescription("filter unhealthy network policies and analyze it"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to analyze network policies in, can be a \"\" string to all namespace"),
				),
			),
			Handler: s.networkPolicyAnalyze,
		}}...,
		)
	}
	return tools
}

func (s *Server) networkPolicyAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	res, err := s.k8s.AnalyzeNetworkPolicy(common.Request{
		Context:   ctx,
		Namespace: ns,
	})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze network policies in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
