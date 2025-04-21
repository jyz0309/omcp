package mcp

import (
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type McpServerState string

const (
	McpServerStateRunning McpServerState = "running"
	McpServerStateStopped McpServerState = "stopped"
)

type MCPServer struct {
	baseServer *server.MCPServer
	*server.SSEServer
	Name      string         `json:"name"`
	Desc      string         `json:"desc"`
	Version   string         `json:"version"`
	State     McpServerState `json:"state"`
	Tools     []MCPTool      `json:"tools"`
	Resources []MCPResource  `json:"resources"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func NewMcpSSEServer(name, desc, version string) *MCPServer {
	mcpServer := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	sseServer := server.NewSSEServer(mcpServer, server.WithBasePath(fmt.Sprintf("/mcp/%s", name)))
	return &MCPServer{
		baseServer: mcpServer,
		SSEServer:  sseServer,
		Name:       name,
		Desc:       desc,
		Version:    version,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		State:      McpServerStateStopped,
	}
}

func (s *MCPServer) Start() {
	s.State = McpServerStateRunning
}

func (s *MCPServer) Stop() {
	s.State = McpServerStateStopped
}

func (s *MCPServer) ListTools() ([]MCPTool, error) {
	return s.Tools, nil
}

func (s *MCPServer) AddTools(tools []MCPTool) {
	for _, tool := range tools {
		tool.Option = append(tool.Option, mcp.WithDescription(tool.Desc))
		s.baseServer.AddTool(mcp.NewTool(tool.Name, tool.Option...), tool.Handler)
	}
	s.Tools = append(s.Tools, tools...)
}

func (s *MCPServer) DeleteTool(name string) {
	s.baseServer.DeleteTools(name)
}

func (s *MCPServer) AddResources(resources []MCPResource) {
	s.Resources = append(s.Resources, resources...)
}
