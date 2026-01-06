package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/common"
)

// Register webhook analysis tools
func (s *Server) initWebhook() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.NewTool("validatingwebhook_analyze",
				mcp.WithDescription("analyze validating webhook configurations"),
				mcp.WithString("name",
					mcp.Description("the name of the validating webhook configuration to analyze"),
				),
				mcp.WithString("label-selector",
					mcp.Description("label selector to filter resources (optional)"),
				),
			),
			Handler: s.validatingWebhookAnalyze,
		},
		{
			Tool: mcp.NewTool("mutatingwebhook_analyze",
				mcp.WithDescription("analyze mutating webhook configurations"),
				mcp.WithString("name",
					mcp.Description("the name of the mutating webhook configuration to analyze"),
				),
				mcp.WithString("label-selector",
					mcp.Description("label selector to filter resources (optional)"),
				),
			),
			Handler: s.mutatingWebhookAnalyze,
		},
	}
}

// Handler for validating webhook analysis
func (s *Server) validatingWebhookAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := ctr.GetString("name", "")
	labelSelector := ctr.GetString("label-selector", "")
	res, err := s.k8s.AnalyzeValidatingWebhook(common.Request{
		Context:       ctx,
		Name:          name,
		LabelSelector: labelSelector,
	})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze validating webhooks: %v", err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

// Handler for mutating webhook analysis
func (s *Server) mutatingWebhookAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := ctr.GetString("name", "")
	labelSelector := ctr.GetString("label-selector", "")

	res, err := s.k8s.AnalyzeMutatingWebhook(common.Request{
		Context:       ctx,
		Name:          name,
		LabelSelector: labelSelector,
	})
	if err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("failed to analyze mutating webhooks: %v", err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
