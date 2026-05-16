package cmd

import (
	"clash-cli/client"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	apiURL  string
	secret  string
)

// rootCmd 是 clash-cli 的根命令
var rootCmd = &cobra.Command{
	Use:   "clash-cli",
	Short: "Clash 代理管理命令行工具",
	Long: `clash-cli 是一个基于 Clash RESTful API 的命令行管理工具。

它允许你通过命令行轻松管理本地或远程运行的 Clash 实例，包括：
  • 查看所有代理节点和代理组
  • 切换代理组中的活动节点
  • 测试单个节点或整组节点的延迟
  • 查看当前活动的代理节点

所有命令均通过 Clash RESTful API 与 Clash 核心通信，
请确保 Clash 已启动且 RESTful API 已开启。`,
	Example: `  # 查看所有代理节点
  clash-cli proxies

  # 切换代理节点
  clash-cli switch "Proxy Group" "HK Node 01"

  # 测试节点延迟
  clash-cli delay "HK Node 01"

  # 测试整组节点延迟
  clash-cli delay-group "Proxy Group"`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://127.0.0.1:9090", "Clash RESTful API 地址")
	rootCmd.PersistentFlags().StringVar(&secret, "secret", "", "Clash RESTful API 认证密钥 (Secret)")
}

// newClient 根据全局 flag 创建 Clash API 客户端
func newClient() *client.Client {
	return client.NewClient(apiURL, secret)
}

// ==================== proxies 命令 ====================

var proxiesCmd = &cobra.Command{
	Use:   "proxies",
	Short: "获取所有代理服务器列表",
	Long: `获取 Clash 中所有可用的代理服务器和代理组信息。

该命令会列出所有代理组和节点，以树状结构展示：
  • 代理组会显示其类型（Selector、URLTest、Fallback 等）和当前活动节点
  • 每个代理组下会列出所有可选的节点
  • 非 Selector 类型的代理组（如 Shadowsocks、VMess 等）也会列出`,
	Example: `  # 查看所有代理节点
  clash-cli proxies

  # 指定 Clash API 地址
  clash-cli proxies --api-url http://192.168.1.100:9090

  # 指定认证密钥
  clash-cli proxies --secret your-secret-key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		resp, err := c.GetProxies()
		if err != nil {
			return fmt.Errorf("获取代理列表失败: %w", err)
		}

		// 区分代理组和普通节点
		var groups []string
		var nodes []string
		for name, detail := range resp.Proxies {
			if len(detail.All) > 0 {
				groups = append(groups, name)
			} else {
				nodes = append(nodes, name)
			}
		}
		sort.Strings(groups)
		sort.Strings(nodes)

		fmt.Println("╔══════════════════════════════════════════════════════════════")
		fmt.Println("║                    🌐 Clash 代理列表                        ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════")
		fmt.Println()

		// 显示代理组
		if len(groups) > 0 {
			fmt.Println("📂 代理组:")
			fmt.Println("──────────────────────────────────────────────────────────────")
			for _, groupName := range groups {
				detail := resp.Proxies[groupName]
				fmt.Printf("  📁 %s [%s]\n", groupName, detail.Type)
				if detail.Now != "" {
					fmt.Printf("     ├─ 当前节点: ✅ %s\n", detail.Now)
				}
				sort.Strings(detail.All)
				for i, node := range detail.All {
					prefix := "     ├─"
					if i == len(detail.All)-1 {
						prefix = "     └─"
					}
					marker := "  "
					if node == detail.Now {
						marker = "✅"
					}
					fmt.Printf("     %s %s %s\n", prefix, marker, node)
				}
				fmt.Println()
			}
		}

		// 显示普通节点
		if len(nodes) > 0 {
			fmt.Println("🔌 节点:")
			fmt.Println("──────────────────────────────────────────────────────────────")
			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(tw, "  名称\t类型\tUDP\n")
			fmt.Fprintf(tw, "  ─────────────────────────────────────\t──────────────\t─────\n")
			for _, nodeName := range nodes {
				detail := resp.Proxies[nodeName]
				udpStr := "❌"
				if detail.UDP {
					udpStr = "✅"
				}
				fmt.Fprintf(tw, "  %s\t%s\t%s\n", nodeName, detail.Type, udpStr)
			}
			tw.Flush()
		}

		fmt.Println()
		fmt.Printf("共 %d 个代理组, %d 个节点\n", len(groups), len(nodes))
		return nil
	},
}

// ==================== switch 命令 ====================

var switchCmd = &cobra.Command{
	Use:   "switch <代理组名> <节点名>",
	Short: "切换代理组的活动节点",
	Long: `切换指定代理组中的活动节点。

该命令会将指定代理组的当前选择切换到新的节点。
只有类型为 Selector 的代理组才支持手动切换。

参数：
  • 代理组名: 要切换的代理组名称（如 "Proxy"、"节点选择" 等）
  • 节点名: 要切换到的节点名称（必须是该代理组中的有效节点）`,
	Example: `  # 切换 "Proxy" 组的节点到 "HK Node 01"
  clash-cli switch "Proxy" "HK Node 01"

  # 切换中文命名的代理组
  clash-cli switch "节点选择" "香港 01"

  # 如果节点名或组名不含空格，引号可省略
  clash-cli switch Proxy HKNode01`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		group := args[0]
		node := args[1]

		c := newClient()

		// 先验证代理组是否存在
		detail, err := c.GetProxy(group)
		if err != nil {
			return fmt.Errorf("代理组 '%s' 不存在或无法访问: %w", group, err)
		}

		// 检查节点是否在该代理组中
		found := false
		for _, n := range detail.All {
			if n == node {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("节点 '%s' 不在代理组 '%s' 中\n\n可用节点:\n%s",
				node, group, formatNodeList(detail.All))
		}

		// 执行切换
		if err := c.SwitchProxy(group, node); err != nil {
			return fmt.Errorf("切换代理节点失败: %w", err)
		}

		fmt.Printf("✅ 成功切换代理组 '%s' 的节点:\n", group)
		fmt.Printf("   %s → %s\n", detail.Now, node)
		return nil
	},
}

func formatNodeList(nodes []string) string {
	sort.Strings(nodes)
	var sb strings.Builder
	for _, n := range nodes {
		sb.WriteString("  • " + n + "\n")
	}
	return sb.String()
}

// ==================== delay 命令 ====================

var (
	testURL string
	timeout int
)

var delayCmd = &cobra.Command{
	Use:   "delay <节点名>",
	Short: "测试指定代理节点的延迟",
	Long: `测试指定代理节点的网络延迟。

该命令会通过指定节点发送 HTTP 请求来测量延迟时间（单位：毫秒）。
可以通过 --url 参数自定义测试 URL，通过 --timeout 设置超时时间。`,
	Example: `  # 测试 "HK Node 01" 的延迟
  clash-cli delay "HK Node 01"

  # 使用自定义测试 URL
  clash-cli delay "HK Node 01" --url https://www.gstatic.com/generate_204

  # 设置超时时间为 5000ms
  clash-cli delay "HK Node 01" --timeout 5000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		node := args[0]
		c := newClient()

		resp, err := c.TestDelay(node, testURL, timeout)
		if err != nil {
			return fmt.Errorf("测试节点 '%s' 延迟失败: %w", node, err)
		}

		if resp.Delay == 0 {
			fmt.Printf("❌ 节点 '%s' 超时或不可达\n", node)
			if resp.Message != "" {
				fmt.Printf("   原因: %s\n", resp.Message)
			}
		} else {
			quality := getDelayQuality(resp.Delay)
			fmt.Printf("📡 节点 '%s' 延迟: %dms %s\n", node, resp.Delay, quality)
		}
		return nil
	},
}

// delayGroupCmd 测试整个代理组的延迟
var delayGroupCmd = &cobra.Command{
	Use:   "delay-group <代理组名>",
	Short: "测试指定代理组中所有节点的延迟",
	Long: `测试指定代理组中所有节点的网络延迟。

该命令会同时对代理组中的所有节点进行延迟测试，
并以表格形式按延迟从低到高排序输出结果。

节点质量评定标准：
  • 🟢 优秀: < 200ms
  • 🟡 一般: 200-500ms
  • 🟠 较差: 500-1000ms
  • 🔴 超时: > 1000ms 或不可达`,
	Example: `  # 测试 "Proxy" 组中所有节点的延迟
  clash-cli delay-group "Proxy"

  # 测试中文命名的代理组
  clash-cli delay-group "节点选择"

  # 使用自定义测试 URL 和超时时间
  clash-cli delay-group "Proxy" --url https://www.gstatic.com/generate_204 --timeout 5000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		group := args[0]
		c := newClient()

		// 先获取代理组信息
		detail, err := c.GetProxy(group)
		if err != nil {
			return fmt.Errorf("获取代理组 '%s' 失败: %w", group, err)
		}
		if len(detail.All) == 0 {
			return fmt.Errorf("'%s' 不是一个代理组", group)
		}

		fmt.Printf("⏳ 正在测试代理组 '%s' 中 %d 个节点的延迟...\n\n", group, len(detail.All))

		// 逐个测试节点延迟
		delays := make(map[string]int)
		for _, node := range detail.All {
			resp, err := c.TestDelay(node, testURL, timeout)
			if err != nil || resp.Delay == 0 {
				delays[node] = 0
			} else {
				delays[node] = resp.Delay
			}
		}

		// 按延迟排序
		type nodeDelay struct {
			name  string
			delay int
		}
		var results []nodeDelay
		for name, d := range delays {
			results = append(results, nodeDelay{name, d})
		}
		sort.Slice(results, func(i, j int) bool {
			// 超时(0)排到最后
			if results[i].delay == 0 {
				return false
			}
			if results[j].delay == 0 {
				return true
			}
			return results[i].delay < results[j].delay
		})

		fmt.Println("┌────────────────────────────────────────────────────────────┐")
		fmt.Printf("│  %-50s │\n", fmt.Sprintf("📡 代理组 '%s' 延迟测试结果", group))
		fmt.Println("├──────────────────────────────────┬───────────┬────────────┤")
		fmt.Println("│ 节点名称                         │ 延迟(ms)  │ 质量       │")
		fmt.Println("├──────────────────────────────────┼───────────┼────────────┤")

		for _, r := range results {
			name := r.name
			if len(name) > 30 {
				name = name[:27] + "..."
			}

			if r.delay == 0 {
				fmt.Printf("│ %-30s  │  超时     │ 🔴 不可达  │\n", name)
			} else {
				quality := getDelayQuality(r.delay)
				fmt.Printf("│ %-30s  │  %5d    │ %s   │\n", name, r.delay, quality)
			}
		}

		fmt.Println("└──────────────────────────────────┴───────────┴────────────┘")
		return nil
	},
}

func getDelayQuality(delay int) string {
	switch {
	case delay < 200:
		return "🟢 优秀"
	case delay < 500:
		return "🟡 一般"
	case delay < 1000:
		return "🟠 较差"
	default:
		return "🔴 很差"
	}
}

// ==================== current 命令 ====================

var currentCmd = &cobra.Command{
	Use:   "current [代理组名]",
	Short: "查看当前活动的代理节点",
	Long: `查看指定代理组或所有代理组的当前活动节点。

如果不指定代理组名，则显示所有代理组的当前活动节点。
如果指定代理组名，则只显示该代理组的当前活动节点。`,
	Example: `  # 查看所有代理组的当前活动节点
  clash-cli current

  # 查看指定代理组的当前活动节点
  clash-cli current "Proxy"
  clash-cli current "节点选择"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		resp, err := c.GetProxies()
		if err != nil {
			return fmt.Errorf("获取代理列表失败: %w", err)
		}

		fmt.Println("📌 当前活动节点:")
		fmt.Println("──────────────────────────────────────────────────────────────")

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "  代理组\t类型\t当前节点\n")
		fmt.Fprintf(tw, "  ───────────────────\t────────────\t──────────────────\n")

		// 过滤和排序
		var groupNames []string
		for name, detail := range resp.Proxies {
			if len(detail.All) > 0 {
				if len(args) > 0 && name != args[0] {
					continue
				}
				groupNames = append(groupNames, name)
			}
		}
		sort.Strings(groupNames)

		if len(groupNames) == 0 && len(args) > 0 {
			return fmt.Errorf("代理组 '%s' 不存在", args[0])
		}

		for _, name := range groupNames {
			detail := resp.Proxies[name]
			now := detail.Now
			if now == "" {
				now = "(无)"
			}
			fmt.Fprintf(tw, "  %s\t%s\t✅ %s\n", name, detail.Type, now)
		}
		tw.Flush()

		return nil
	},
}

func init() {
	// delay 和 delay-group 共用的 flags
	delayCmd.Flags().StringVar(&testURL, "url", "https://www.gstatic.com/generate_204", "测速使用的 URL")
	delayCmd.Flags().IntVar(&timeout, "timeout", 3000, "测速超时时间 (毫秒)")
	delayGroupCmd.Flags().StringVar(&testURL, "url", "https://www.gstatic.com/generate_204", "测速使用的 URL")
	delayGroupCmd.Flags().IntVar(&timeout, "timeout", 5000, "测速超时时间 (毫秒)")

	// 注册子命令
	rootCmd.AddCommand(proxiesCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(delayCmd)
	rootCmd.AddCommand(delayGroupCmd)
	rootCmd.AddCommand(currentCmd)
}
