package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wenhuwang/mcp-k8s-eye/pkg/common"
)

func (s *Server) initResource() []server.ServerTool {
	tools := []server.ServerTool{
		{
			Tool: mcp.NewTool("resource_list",
				mcp.WithDescription("list resources in a namespace or all namespaces"),
				mcp.WithString("kind",
					mcp.Description("the kind of resource to list"),
					mcp.Required(),
				),
				mcp.WithString("namespace",
					mcp.Description("the namespace to list resources in, can be a \"\" string to list resources in all namespace"),
				),
				mcp.WithString("labelSelector",
					mcp.Description("the label selector to filter list resource, , can be a \"\" string if no selector"),
				),
			),
			Handler: s.resourceList,
		},
		{
			Tool: mcp.NewTool("resource_get",
				mcp.WithDescription("get resource details"),
				mcp.WithString("kind",
					mcp.Description("the kind of resource to get"),
					mcp.Required(),
				),
				mcp.WithString("namespace",
					mcp.Description("the namespace to get resources in"),
				),
				mcp.WithString("name",
					mcp.Description("the resource name to get"),
					mcp.Required(),
				),
			),
			Handler: s.resourceGet,
		},
		{
			Tool: mcp.NewTool("resource_describe",
				mcp.WithDescription("describe resource"),
				mcp.WithString("kind",
					mcp.Description("the resource kind to describe"),
					mcp.Required(),
				),
				mcp.WithString("namespace",
					mcp.Description("the resource namespace to describe"),
				),
				mcp.WithString("name",
					mcp.Description("the resource name to describe"),
					mcp.Required(),
				),
			),
			Handler: s.ResourceDescribe,
		},
		{
			Tool: mcp.NewTool("workload_resource_usage",
				mcp.WithDescription("workload resource usage"),
				mcp.WithString("kind",
					mcp.Description("the kind of workload"),
					mcp.Required(),
					mcp.Enum("Deployment", "StatefulSet", "DaemonSet", "ReplicaSet", "Pod"),
				),
				mcp.WithString("namespace",
					mcp.Description("the namespace of workload"),
					mcp.Required(),
				),
				mcp.WithString("name",
					mcp.Description("the name of workload"),
				),
			),
			Handler: s.workloadResourceUsage,
		},
	}
	if s.write {
		tools = append(tools,
			server.ServerTool{
				Tool: mcp.NewTool("resource_delete",
					mcp.WithDescription("delete resource"),
					mcp.WithString("kind",
						mcp.Description("the kind of resource to delete"),
						mcp.Required(),
					),
					mcp.WithString("namespace",
						mcp.Description("the namespace to get resources in"),
						mcp.Required(),
					),
					mcp.WithString("name",
						mcp.Description("the resource name to delete"),
						mcp.Required(),
					),
				),
				Handler: s.resourceDelete,
			},
			server.ServerTool{
				Tool: mcp.NewTool("resource_create_or_update",
					mcp.WithDescription("create or update resource"),
					mcp.WithString("resource",
						mcp.Description("the resource to create or update"),
						mcp.Required(),
					),
				),
				Handler: s.resourceCreateOrUpdate,
			})
	}
	return tools
}

func (s *Server) resourceList(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	kind := ctr.GetString("kind", "")
	labelSelector := ctr.GetString("labelSelector", "")
	res, err := s.k8s.ResourceList(ctx, kind, ns, labelSelector)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list resources in namespace %s: %v", ns, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) resourceGet(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	kind := ctr.GetString("kind", "")
	name := ctr.GetString("name", "")
	res, err := s.k8s.ResourceGet(ctx, kind, ns, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get resource %s/%s: %v", ns, name, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) resourceDelete(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ns := ctr.GetString("namespace", "")
	kind := ctr.GetString("kind", "")
	name := ctr.GetString("name", "")
	res, err := s.k8s.ResourceDelete(ctx, kind, ns, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to delete resource %s/%s: %v", ns, name, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) resourceCreateOrUpdate(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resource := ctr.GetString("resource", "")
	res, err := s.k8s.ResourceCreateOrUpdate(ctx, resource)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create/update resource: %v", err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) ResourceDescribe(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	r := common.Request{
		Context:   ctx,
		Kind:      ctr.GetString("kind", ""),
		Name:      ctr.GetString("name", ""),
		Namespace: ctr.GetString("namespace", ""),
	}
	res, err := s.k8s.ResourceDescribe(r)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to describe resource: %v", err)), nil
	}
	return mcp.NewToolResultText(res), nil
}

func (s *Server) workloadResourceUsage(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	namespace := ctr.GetString("namespace", "")
	kind := ctr.GetString("kind", "")
	name := ctr.GetString("name", "")
	res, err := s.k8s.WorkloadResourceUsage(common.Request{
		Context:   ctx,
		Name:      name,
		Namespace: namespace,
		Kind:      kind,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get resource usage in namesoace %s: %v", namespace, err)), nil
	}
	return mcp.NewToolResultText(res), nil
}
