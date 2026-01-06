package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/common"
)

func (s *Server) initStatefulSet() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("statefulset_analyze",
				mcp.WithDescription("filter unhealthy statefulset and analyze it"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to analyze statefulset in"),
					mcp.Required(),
				),
			),
			Handler: s.statefulSetAnalyze,
		},
	}
}
func (s *Server) statefulSetAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	r := common.Request{
		Context:   ctx,
		Namespace: ns,
	}
	res, err := s.k8s.AnalyzeStatefulSet(r)
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze statefulsets in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
