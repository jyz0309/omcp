package web

import (
	"github.com/jyz0309/omcp/mcp"
)

type ServerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AddToolReq struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	FilePath string `json:"filePath"`
	Func     string `json:"func"`
}

type DeleteToolReq struct {
	ToolName string `json:"tool_name"`
	Server   string `json:"server"`
}

type CreateMcpServerReq struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Desc    string `json:"desc"`
}

type CreateMcpServerResp struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Server  *mcp.MCPServer `json:"server"`
}

type DeleteMcpServerReq struct {
	Name string `json:"name"`
}

type ListMcpServerReq struct {
	IsAlive bool `json:"is_alive"`
}

type ListMcpServerResp struct {
	Servers []*mcp.MCPServer `json:"servers"`
}

type UpdateMcpServerReq struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type UpdateMcpServerResp struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Server  *mcp.MCPServer `json:"server"`
}

type StartMcpServerReq struct {
	Name string `json:"name"`
}

type StopMcpServerReq struct {
	Name string `json:"name"`
}

// Tool
type ListToolReq struct {
	Server string `json:"server"`
	IsRepo bool   `json:"is_repo"`
}

type ListToolResp struct {
	Total int64         `json:"total"`
	Tools []mcp.MCPTool `json:"tools"`
}

// Load
type LoadReq struct {
	Plugins []mcp.Plugin `json:"plugins"`
}

type LoadResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
