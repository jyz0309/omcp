/*
package main

import (

	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

)

	func main() {
		go func() {
			// Create MCP server
			s := server.NewMCPServer(
				"Demo 🚀",
				"1.0.0",
				server.WithToolCapabilities(true),
			)

			// Add tool
			tool := mcp.NewTool("hello_world",
				mcp.WithDescription("Say hello to someone"),
				mcp.WithString("name",
					mcp.Required(),
					mcp.Description("Name of the person to greet"),
				),
			)

			// Add tool handler
			s.AddTool(tool, helloHandler)
			sse := server.NewSSEServer(s, server.WithBaseURL("http://127.0.0.1:5000"))
			err := sse.Start(":5000")
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(1 * time.Second)
		cli, err := client.NewSSEMCPClient("http://127.0.0.1:5000/sse")
		if err != nil {
			panic(err)
		}

		err = cli.Start(context.Background())
		if err != nil {
			panic(err)
		}

		// Initialize
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "test-client",
			Version: "1.0.0",
		}

		result, err := cli.Initialize(context.Background(), initRequest)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)

		tools, err := cli.ListTools(context.Background(), mcp.ListToolsRequest{})
		if err != nil {
			panic(err)
		}
		fmt.Println(tools)
	}

	func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return nil, errors.New("name must be a string")
		}

		return mcp.NewToolResultText(fmt.Sprintf("Hello, sth %s!", name)), nil
	}
*/
package main

import (
	"github.com/jyz0309/omcp/cmd"
)

func main() {
	cmd.NewCli().Execute()
}
