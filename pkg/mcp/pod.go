package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *Server) initPod() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("pod_logs",
				mcp.WithDescription("get pod logs"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to get pods in"),
					mcp.Required(),
				),
				mcp.WithString("pod",
					mcp.Description("the pod to get"),
					mcp.Required(),
				),
				mcp.WithString("container",
					mcp.Description("the container to get"),
					mcp.Required(),
				),
			),
			Handler: s.podLogs,
		},
		{
			Tool: mcp.NewTool("pod_analyze",
				mcp.WithDescription("analyze pod"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to get pods in"),
					mcp.Required(),
				),
			),
			Handler: s.podAnalyze,
		},
	}
	if s.write {
		tools = append(tools, server.ServerTool{
			Tool: mcp.NewTool("pod_exec",
				mcp.WithDescription("execute a command in a pod"),
				mcp.WithString("namespace",
					mcp.Description("the namespace to get pods in"),
				),
				mcp.WithString("pod",
					mcp.Description("the pod to get"),
				),
				mcp.WithString("command",
					mcp.Description("the command to execute"),
				),
			),
			Handler: s.podExec,
		})
	}
	return tools
}

func (s *Server) podLogs(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	pod := ctr.GetString("pod", "")
	container := ctr.GetString("container", "")
	res, err := s.k8s.PodLogs(ctx, ns, pod, container)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get logs for pod %s/%s: %v", ns, pod, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) podExec(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	pod := ctr.GetString("pod", "")
	cmd := ctr.GetString("command", "")
	res, err := s.k8s.PodExec(ctx, ns, pod, cmd)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to execute command %s on pod %s/%s: %v", cmd, ns, pod, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) podAnalyze(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	res, err := s.k8s.AnalyzePod(ctx, ns)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to analyze pods in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
