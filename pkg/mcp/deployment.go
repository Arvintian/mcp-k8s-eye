package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initDeployment() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("deployment_analyze",
				mcp.WithDescription("analyze deployment status"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to analyze deployments in"),
					mcp.Required(),
				),
			),
			Handler: s.deploymentAnalyze,
		},
	}
	if s.write {
		tools = append(tools, server.ServerTool{
			Tool: mcp.NewTool("deployment_scale",
				mcp.WithDescription("scale deployment replicas"),
				mcp.WithString("namespace",
					mcp.Description("the namespace of the deployment"),
					mcp.Required(),
				),
				mcp.WithString("deployment",
					mcp.Description("the deployment to scale"),
					mcp.Required(),
				),
				mcp.WithNumber("replicas",
					mcp.Description("the number of replicas to scale to"),
					mcp.Required(),
				),
			),
			Handler: s.deploymentScale,
		})
	}
	return tools
}

func (s *Server) deploymentScale(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	deploy := ctr.GetString("deployment", "")
	replicas := int32(ctr.GetInt("replicas", -1))
	if replicas < 0 {
		return mcp.NewToolResultError("scale deployment replicas error"), nil
	}
	res, err := s.k8s.DeploymentScale(ctx, ns, deploy, replicas)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to scale deployment %s/%s: %v", ns, deploy, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) deploymentAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	res, err := s.k8s.AnalyzeDeployment(ctx, ns)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to analyze deployments in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
