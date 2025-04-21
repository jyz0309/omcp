package mcp

import "time"

type MCPResource struct {
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
