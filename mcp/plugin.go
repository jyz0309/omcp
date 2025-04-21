package mcp

type Plugin struct {
	MCPType    string `json:"mcp_type"`
	Name       string `json:"name"`
	VarName    string `json:"var_name"`
	PluginFile string `json:"plugin_file"`
}
