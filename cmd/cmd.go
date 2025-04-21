package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jyz0309/omcp/config"
	"github.com/jyz0309/omcp/web"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewCli() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "omcp",
		Short:         "omcp is a tool for managing MCP server deployments",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				versionHandler(cmd, args)
				return
			}

			cmd.Print(cmd.UsageString())
		},
	}
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the OMCP server",
		RunE:  serveHandler,
	}
	rootCmd.AddCommand(serveCmd)

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Manage MCP servers",
	}
	rootCmd.AddCommand(serverCmd)

	var createCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a new MCP server",
		PreRunE: probeServerReady,
		RunE:    createHandler,
	}
	createCmd.Flags().StringP("name", "n", "", "The name of the MCP server")
	createCmd.Flags().StringP("desc", "d", "", "The description of the MCP server")
	createCmd.Flags().StringP("version", "v", "0.0.1", "The version of the MCP server")
	serverCmd.AddCommand(createCmd)

	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Short:   "Delete a MCP server",
		PreRunE: probeServerReady,
		RunE:    deleteHandler,
	}
	deleteCmd.Flags().StringP("name", "n", "", "The name of the MCP server")
	serverCmd.AddCommand(deleteCmd)

	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List MCP servers",
		PreRunE: probeServerReady,
		RunE:    listHandler,
	}
	serverCmd.AddCommand(listCmd)

	var startCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start a MCP server",
		PreRunE: probeServerReady,
		RunE:    startHandler,
	}
	startCmd.Flags().StringP("name", "n", "", "The name of the MCP server")
	serverCmd.AddCommand(startCmd)

	var stopCmd = &cobra.Command{
		Use:     "stop",
		Short:   "Stop a MCP server",
		PreRunE: probeServerReady,
		RunE:    stopHandler,
	}
	stopCmd.Flags().StringP("name", "n", "", "The name of the MCP server")
	serverCmd.AddCommand(stopCmd)

	return rootCmd
}

func serveHandler(cmd *cobra.Command, args []string) error {
	server := web.NewHttpServer()
	err := server.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}

// createHandler creates a new MCP server
func createHandler(cmd *cobra.Command, args []string) error {
	cli := NewOmcpServerCli(config.Host())
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required")
	}
	desc, _ := cmd.Flags().GetString("desc")
	version, _ := cmd.Flags().GetString("version")
	err := cli.CreateMcpServer(name, desc, version)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	return nil
}

func deleteHandler(cmd *cobra.Command, args []string) error {
	cli := NewOmcpServerCli(config.Host())
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required")
	}
	err := cli.DeleteMcpServer(name)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	return nil
}

func listHandler(cmd *cobra.Command, args []string) error {
	cli := NewOmcpServerCli(config.Host())
	servers, err := cli.ListMcpServers()
	if err != nil {
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Version", "Status", "Created_At", "Updated_At"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	for _, server := range servers {
		table.Append([]string{server.Name, server.Desc, server.Version, string(server.State), server.CreatedAt.Format(time.DateTime), server.UpdatedAt.Format(time.DateTime)})
	}
	table.Render()
	return nil
}

func startHandler(cmd *cobra.Command, args []string) error {
	cli := NewOmcpServerCli(config.Host())
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required")
	}
	err := cli.StartMcpServer(name)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	return nil
}

func stopHandler(cmd *cobra.Command, args []string) error {
	cli := NewOmcpServerCli(config.Host())
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		return fmt.Errorf("name is required")
	}
	err := cli.StopMcpServer(name)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	return nil
}

func loadHandler(cmd *cobra.Command, args []string) error {
	return nil
}
func versionHandler(cmd *cobra.Command, args []string) {
	cmd.Println("omcp version 0.0.1")
}

// probeServerReady probes the omcp server to see if it is ready
func probeServerReady(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(fmt.Sprintf("%s/ready", config.Host()))
	if err != nil {
		cmd.PrintErrln("OMCP server is not ready, please start the server first")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	cmd.PrintErrln("OMCP server is not ready, please start the server first")
	return fmt.Errorf("server is not ready")
}
