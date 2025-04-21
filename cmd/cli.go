package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jyz0309/omcp/config"
	web "github.com/jyz0309/omcp/web"

	"github.com/jyz0309/omcp/mcp"
)

type OmcpServerCli struct {
	url string

	cli *http.Client
}

func NewOmcpServerCli(host string) *OmcpServerCli {
	// default host and port
	if host == "" {
		host = config.Host()
	}
	// TODO
	// else {
	//	 config.SetHost(host)s
	// }

	return &OmcpServerCli{
		url: host,
		cli: &http.Client{},
	}
}

func (c *OmcpServerCli) CreateMcpServer(name, desc, version string) error {
	body := web.CreateMcpServerReq{
		Name:    name,
		Desc:    desc,
		Version: version,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/server/create", c.url), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create mcp server, status code: %d", resp.StatusCode)
	} else {
		var respBody web.CreateMcpServerResp
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		} else if !respBody.Success {
			return fmt.Errorf("failed to create mcp server, message: %s", respBody.Message)
		}
	}

	return nil
}

func (c *OmcpServerCli) ListMcpServers() ([]*mcp.MCPServer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/server/list", c.url), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list mcp servers, status code: %d", resp.StatusCode)
	}

	var respBody web.ListMcpServerResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Servers, nil
}

func (c *OmcpServerCli) DeleteMcpServer(name string) error {
	body := web.DeleteMcpServerReq{
		Name: name,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/server/delete", c.url), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete mcp server, status code: %d", resp.StatusCode)
	} else {
		var respBody web.ServerResp
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		} else if !respBody.Success {
			return fmt.Errorf("failed to delete mcp server, message: %s", respBody.Message)
		}
	}

	return nil
}

func (c *OmcpServerCli) StartMcpServer(name string) error {
	body := web.StartMcpServerReq{
		Name: name,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/server/start", c.url), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start mcp server, status code: %d", resp.StatusCode)
	}

	var respBody web.ServerResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	if !respBody.Success {
		return fmt.Errorf("failed to start mcp server, message: %s", respBody.Message)
	}

	return nil
}

func (c *OmcpServerCli) StopMcpServer(name string) error {
	body := web.StopMcpServerReq{
		Name: name,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/server/stop", c.url), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop mcp server, status code: %d", resp.StatusCode)
	}

	var respBody web.ServerResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	if !respBody.Success {
		return fmt.Errorf("failed to stop mcp server, message: %s", respBody.Message)
	}

	return nil
}

func (c *OmcpServerCli) Load(filepath string, plugins []mcp.Plugin) error {
	body := web.LoadReq{
		Plugins: plugins,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/load", c.url), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
