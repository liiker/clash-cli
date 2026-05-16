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
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	All       []string `json:"all"`
	Now       string   `json:"now"`
	History   []DelayHistory `json:"history"`
	UDP       bool     `json:"udp"`
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

// TestGroupDelay 测试指定代理组中所有节点的延迟
func (c *Client) TestGroupDelay(group string, testURL string, timeout int) (map[string]int, error) {
	encodedGroup := url.PathEscape(group)
	path := fmt.Sprintf("/group/%s/delay?timeout=%d&url=%s", encodedGroup, timeout, url.QueryEscape(testURL))

	data, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var resp map[string]int
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return resp, nil
}
