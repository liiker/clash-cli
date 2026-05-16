package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ==================== 数据结构定义 ====================

// ProxiesResponse 是 GET /proxies 的响应结构
type ProxiesResponse struct {
	Proxies map[string]ProxyDetail `json:"proxies"`
}

// ProxyDetail 是单个代理节点的详细信息
type ProxyDetail struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	All     []string      `json:"all"`
	Now     string        `json:"now"`
	History []DelayHistory `json:"history"`
	UDP     bool          `json:"udp"`
}

// DelayHistory 是延迟测试历史记录
type DelayHistory struct {
	Time  string `json:"time"`
	Delay int    `json:"delay"`
}

// DelayResponse 是 GET /proxies/{name}/delay 的响应结构
type DelayResponse struct {
	Delay   int    `json:"delay"`
	Message string `json:"message"`
}

// ConfigResponse 是 GET /configs 的响应结构
type ConfigResponse struct {
	Port       int           `json:"port"`
	SocksPort  int           `json:"socks-port"`
	RedirPort  int           `json:"redir-port"`
	TProxyPort int           `json:"tproxy-port"`
	MixedPort  int           `json:"mixed-port"`
	AllowLan   bool          `json:"allow-lan"`
	BindAddress string       `json:"bind-address"`
	Mode       string        `json:"mode"`
	LogLevel   string        `json:"log-level"`
	IPv6       bool          `json:"ipv6"`
}

// ConnectionsResponse 是 GET /connections 的响应结构
type ConnectionsResponse struct {
	DownloadTotal int64        `json:"downloadTotal"`
	UploadTotal   int64        `json:"uploadTotal"`
	Connections   []Connection `json:"connections"`
}

// Connection 是单个连接的详细信息
type Connection struct {
	ID         string            `json:"id"`
	Metadata   ConnectionMetadata `json:"metadata"`
	Upload     int64             `json:"upload"`
	Download   int64             `json:"download"`
	Start      string            `json:"start"`
	Chains     []string          `json:"chains"`
	Rule       string            `json:"rule"`
	RulePayload string           `json:"rulePayload"`
}

// ConnectionMetadata 是连接的元数据
type ConnectionMetadata struct {
	Network  string `json:"network"`
	Type     string `json:"type"`
	SourceIP string `json:"sourceIP"`
	DstIP    string `json:"destinationIP"`
	SourcePort string `json:"sourcePort"`
	DstPort  string `json:"destinationPort"`
	Host     string `json:"host"`
	DNSMode  string `json:"dnsMode"`
	Process  string `json:"process"`
	ProcessPath string `json:"processPath"`
}

// ProvidersResponse 是 GET /providers/proxies 的响应结构
type ProvidersResponse struct {
	Providers map[string]ProviderDetail `json:"providers"`
}

// ProviderDetail 是单个代理提供者的详细信息
type ProviderDetail struct {
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	VehicleType string              `json:"vehicleType"`
	Proxies     []ProviderProxyItem `json:"proxies"`
	UpdatedAt   string              `json:"updatedAt"`
	SubscriptionInfo struct {
		Upload   int64 `json:"Upload"`
		Download int64 `json:"Download"`
		Total    int64 `json:"Total"`
		Expire   int64 `json:"Expire"`
	} `json:"subscriptionInfo"`
}

// ProviderProxyItem 是代理提供者中的节点条目
type ProviderProxyItem struct {
	Name string `json:"name"`
}

// VersionResponse 是 GET /version 的响应结构
type VersionResponse struct {
	Version string `json:"version"`
}

// ==================== API 方法 ====================

// GetProxies 获取所有代理节点列表
func (c *Client) GetProxies() (*ProxiesResponse, error) {
	data, err := c.Get("/proxies")
	if err != nil {
		return nil, err
	}

	var resp ProxiesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// GetProxy 获取指定代理节点的详细信息
func (c *Client) GetProxy(name string) (*ProxyDetail, error) {
	encodedName := url.PathEscape(name)
	data, err := c.Get("/proxies/" + encodedName)
	if err != nil {
		return nil, err
	}

	var detail ProxyDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &detail, nil
}

// SwitchProxy 切换指定代理组中的节点
func (c *Client) SwitchProxy(group, node string) error {
	encodedGroup := url.PathEscape(group)
	payload := map[string]string{"name": node}
	_, err := c.Put("/proxies/"+encodedGroup, payload)
	return err
}

// TestDelay 测试指定节点的延迟
func (c *Client) TestDelay(name string, testURL string, timeout int) (*DelayResponse, error) {
	encodedName := url.PathEscape(name)
	path := fmt.Sprintf("/proxies/%s/delay?timeout=%d&url=%s", encodedName, timeout, url.QueryEscape(testURL))

	data, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var resp DelayResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// GetConfigs 获取 Clash 运行配置
func (c *Client) GetConfigs() (*ConfigResponse, error) {
	data, err := c.Get("/configs")
	if err != nil {
		return nil, err
	}

	var resp ConfigResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// PatchConfigs 修改 Clash 运行配置（如 mode、log-level 等）
func (c *Client) PatchConfigs(payload map[string]interface{}) error {
	_, err := c.doPatch("/configs", payload)
	return err
}

// ReloadConfigs 热重载配置文件
func (c *Client) ReloadConfigs(configPath string) error {
	payload := map[string]string{"path": configPath}
	_, err := c.Put("/configs?force=true", payload)
	return err
}

// GetConnections 获取所有活跃连接
func (c *Client) GetConnections() (*ConnectionsResponse, error) {
	data, err := c.Get("/connections")
	if err != nil {
		return nil, err
	}

	var resp ConnectionsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// CloseAllConnections 关闭所有连接
func (c *Client) CloseAllConnections() error {
	_, err := c.doDelete("/connections")
	return err
}

// CloseConnection 关闭指定连接
func (c *Client) CloseConnection(id string) error {
	_, err := c.doDelete("/connections/" + url.PathEscape(id))
	return err
}

// GetProviders 获取所有代理提供者
func (c *Client) GetProviders() (*ProvidersResponse, error) {
	data, err := c.Get("/providers/proxies")
	if err != nil {
		return nil, err
	}

	var resp ProvidersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// UpdateProvider 更新（刷新）指定代理提供者
func (c *Client) UpdateProvider(name string) error {
	encodedName := url.PathEscape(name)
	_, err := c.Put("/providers/proxies/"+encodedName, nil)
	return err
}

// HealthCheckProvider 触发指定代理提供者的健康检查
func (c *Client) HealthCheckProvider(name string) error {
	encodedName := url.PathEscape(name)
	_, err := c.Get("/providers/proxies/" + encodedName + "/healthcheck")
	return err
}

// GetVersion 获取 Clash 版本
func (c *Client) GetVersion() (*VersionResponse, error) {
	data, err := c.Get("/version")
	if err != nil {
		return nil, err
	}

	var resp VersionResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}
