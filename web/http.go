package web

import (
	"os"

	"github.com/jyz0309/omcp/mcp"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/client"
	"github.com/sirupsen/logrus"
)

type OmcpServer struct {
	*gin.Engine

	logger       *logrus.Logger
	MCPServerMap map[string]*mcp.MCPServer
}

func NewHttpServer() *OmcpServer {
	r := gin.Default()
	// TODO: test, remove it
	logger := logrus.New()
	logger.Out = os.Stdout
	gin.DefaultWriter = logger.Out

	omcpServer := OmcpServer{
		Engine:       r,
		logger:       logger,
		MCPServerMap: make(map[string]*mcp.MCPServer),
	}
	// test
	mcpServer := mcp.NewMcpSSEServer("hello", "hello", "1.0.0")
	omcpServer.MCPServerMap["hello"] = mcpServer

	r.GET("/ready", omcpServer.HandleReady)
	// server api
	r.GET("/api/server/list", omcpServer.ListMcpServer)
	r.POST("/api/server/create", omcpServer.CreateMcpServer)
	r.POST("/api/server/delete", omcpServer.DeleteMcpServer)
	r.POST("/api/server/start", omcpServer.StartMcpServer)
	r.POST("/api/server/stop", omcpServer.StopMcpServer)

	// tool api
	r.GET("/api/tool/list", omcpServer.ListTool)

	// load plugin api
	r.POST("/api/load", omcpServer.Load)
	// sse api
	r.GET("/mcp/ping", omcpServer.HandlePing)
	r.GET("/mcp/:name/sse", omcpServer.HandleSSE)
	r.POST("/mcp/:name/message", omcpServer.HandleMessage)

	return &omcpServer
}

func (s *OmcpServer) Run(addr string) error {
	return s.Engine.Run(addr)
}

// HandleReady checks if OMCP server is ready
func (s *OmcpServer) HandleReady(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ready",
	})
}

// HandlePing checks if all the SSE server is ready
func (s *OmcpServer) HandlePing(c *gin.Context) {
	for _, sseServer := range s.MCPServerMap {
		cli, err := client.NewSSEMCPClient(sseServer.CompleteSsePath())
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error",
			})
			return
		}
		err = cli.Ping(c)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "error",
			})
			return
		}
		cli.Close()
	}

	c.JSON(200, gin.H{
		"message": "ready",
	})
}

// HandleSSE handles the MCP server SSE request
func (s *OmcpServer) HandleSSE(c *gin.Context) {
	name := c.Param("name")
	sseServer, exist := s.MCPServerMap[name]

	if !exist {
		// TODO: 转发到别的port
		c.JSON(200, ServerResp{
			Success: false,
			Message: "server not found",
		})
		return
	}
	if sseServer.State == mcp.McpServerStateRunning {
		sseServer.ServeHTTP(c.Writer, c.Request)
	} else {
		c.JSON(200, ServerResp{
			Success: false,
			Message: "server is not running",
		})
		return
	}
}

// HandleMessage handles the MCP server message request
func (s *OmcpServer) HandleMessage(c *gin.Context) {
	s.logger.Info(c.Request.URL.Path)
	name := c.Param("name")
	sseServer, exist := s.MCPServerMap[name]
	if !exist {
		c.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	if sseServer.State == mcp.McpServerStateRunning {
		sseServer.ServeHTTP(c.Writer, c.Request)
	} else {
		c.JSON(404, gin.H{
			"message": "server is not running",
		})
		return
	}
}

func (s *OmcpServer) CreateMcpServer(c *gin.Context) {
	var req CreateMcpServerReq
	s.logger.Error(c.Request.URL.Path)
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, CreateMcpServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}

	if _, exist := s.MCPServerMap[req.Name]; exist {
		s.logger.Error("mcp server already exists")
		c.JSON(200, CreateMcpServerResp{
			Success: false,
			Message: "mcp server already exists",
		})
		return
	}

	mcpServer := mcp.NewMcpSSEServer(req.Name, req.Desc, req.Version)
	s.MCPServerMap[req.Name] = mcpServer
	c.JSON(200, CreateMcpServerResp{
		Success: true,
		Message: "success",
		Server:  mcpServer,
	})
}

func (s *OmcpServer) DeleteMcpServer(c *gin.Context) {
	var req DeleteMcpServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, ServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	/*
		we don't need to shutdown the sse server here,
		the sse server have not start, we just use the ServeHTTP() func
		sseServer, exist := s.SSEServerMcp[req.Name]
		if !exist {
			c.JSON(404, gin.H{
				"message": "not found",
			})
			return
		}
		sseServer.Shutdown()
	*/

	delete(s.MCPServerMap, req.Name)
	c.JSON(200, ServerResp{
		Success: true,
		Message: "success",
	})
}

func (s *OmcpServer) ListMcpServer(c *gin.Context) {
	var req ListMcpServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, ServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	resp := &ListMcpServerResp{}
	for _, sseServer := range s.MCPServerMap {
		if req.IsAlive {
			if sseServer.State == mcp.McpServerStateRunning {
				resp.Servers = append(resp.Servers, sseServer)
			}
		} else {
			resp.Servers = append(resp.Servers, sseServer)
		}
	}
	c.JSON(200, resp)
}

func (s *OmcpServer) StartMcpServer(c *gin.Context) {
	var req StartMcpServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, ServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	sseServer, exist := s.MCPServerMap[req.Name]
	if !exist {
		c.JSON(200, ServerResp{
			Success: false,
			Message: "not found",
		})
		return
	}
	sseServer.Start()
	c.JSON(200, ServerResp{
		Success: true,
		Message: "success",
	})
}

func (s *OmcpServer) StopMcpServer(c *gin.Context) {
	var req StopMcpServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, ServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	sseServer, exist := s.MCPServerMap[req.Name]
	if !exist {
		c.JSON(200, ServerResp{
			Success: false,
			Message: "not found",
		})
		return
	}
	sseServer.Stop()
	s.logger.Info("stop mcp server", req.Name)
	c.JSON(200, ServerResp{
		Success: true,
		Message: "success",
	})
}

func (s *OmcpServer) ListTool(c *gin.Context) {
	var req ListToolReq
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error(err)
		c.JSON(200, ServerResp{
			Success: false,
			Message: "invalid request",
		})
		return
	}
	mcpServer, exist := s.MCPServerMap[req.Server]
	if !exist {
		c.JSON(200, ServerResp{
			Success: false,
			Message: "not found",
		})
		return
	}
	tools, err := mcpServer.ListTools()
	if err != nil {
		c.JSON(200, ServerResp{
			Success: false,
			Message: "error",
		})
		return
	}
	c.JSON(200, ListToolResp{
		Total: int64(len(tools)),
		Tools: tools,
	})
}

func (s *OmcpServer) Load(c *gin.Context) {
	pluginFile, err := c.FormFile("plugin_file")
	if err != nil {
		c.JSON(200, LoadResp{
			Success: false,
			Message: "error",
		})
		return
	}
	dst := "./plugins/" + pluginFile.Filename
	if err := c.SaveUploadedFile(pluginFile, dst); err != nil {
		c.JSON(200, LoadResp{
			Success: false,
			Message: "error",
		})
		return
	}
	c.JSON(200, LoadResp{
		Success: true,
		Message: "success",
	})
}
