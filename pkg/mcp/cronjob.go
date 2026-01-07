package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/common"
)

// Register cronjob analyze tool
func (s *Server) initCronJob() []server.ServerTool {
	tools := []server.ServerTool{}
	if s.analyze {
		tools = append(tools, []server.ServerTool{{
			Tool: mcp.NewTool("cronjob_analyze",
				mcp.WithDescription("filter unhealthy cronjob and analyze it"),
				mcp.WithString("namespace",
					mcp.Description("the cronjob namespace to analyze, can be a \"\" string to all namespace"),
				),
			),
			Handler: s.cronjobAnalyze,
		}}...,
		)
	}
	return tools
}

func (s *Server) cronjobAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	res, err := s.k8s.AnalyzeCronJob(common.Request{
		Context:   ctx,
		Namespace: ns,
	})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze cronjob in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
