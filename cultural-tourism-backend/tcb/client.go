// File: tcb/client.go
package tcb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Global Client Instance
var Client *CloudBaseClient

type CloudBaseClient struct {
	EnvID       string
	AccessToken string
	BaseURL     string
	HTTPClient  *http.Client
}

// Init 初始化全局客户端
func Init() {
	_ = godotenv.Load() // 加载 .env (本地开发用)

	envID := os.Getenv("CLOUDBASE_ENV_ID")
	accessToken := os.Getenv("CLOUDBASE_ACCESS_TOKEN")

	if envID == "" {
		fmt.Println("⚠️ 警告: 未找到 CLOUDBASE_ENV_ID 环境变量")
	}

	Client = &CloudBaseClient{
		EnvID:       envID,
		AccessToken: accessToken,
		BaseURL:     fmt.Sprintf("https://%s.api.tcloudbasegateway.com", envID),
		HTTPClient:  &http.Client{},
	}
	fmt.Println("✅ 云开发 HTTP 客户端已初始化 (通用模式)")
}

// Request 通用 HTTP 请求处理
func (c *CloudBaseClient) Request(method, path string, body interface{}, customHeaders map[string]string) (interface{}, error) {
	url := c.BaseURL + path
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("JSON序列化失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// 鉴权核心：云托管内网或本地调试通过 Token 访问
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 错误处理：非 2xx 视为错误
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API错误 [%d]: %s", resp.StatusCode, string(bodyBytes))
	}

	var result interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("JSON解析失败: %w", err)
		}
	}
	return result, nil
}

// CreateData 新增数据
// API: POST /v1/model/prod/{modelName}/create
func (c *CloudBaseClient) CreateData(modelName string, data interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v1/model/prod/%s/create", modelName)
	payload := map[string]interface{}{
		"data": data,
	}

	result, err := c.Request("POST", path, payload, nil)
	if err != nil {
		return nil, err
	}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}
	return nil, fmt.Errorf("返回格式异常")
}

// ListData 查询列表 (核心修正版)
// ⚠️ 严禁在此处硬编码 $eq 等逻辑。
// filter 参数必须由 Controller 层构造成完整的 TCB 查询对象 (包含 where, orderBy 等)
// API: POST /v1/model/prod/{modelName}/list
func (c *CloudBaseClient) ListData(modelName string, filter map[string]interface{}, page, size int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v1/model/prod/%s/list", modelName)

	payload := map[string]interface{}{
		"pageNumber": page,
		"pageSize":   size,
		"getCount":   true,
	}

	// [Audit Fix]: 仅仅透传 filter，不做任何假设或加工
	// 调用方 (Controller) 负责构造 { "where": {...}, "orderBy": [...] }
	if filter != nil && len(filter) > 0 {
		payload["filter"] = filter
	}

	result, err := c.Request("POST", path, payload, nil)
	if err != nil {
		return nil, err
	}
	if resultMap, ok := result.(map[string]interface{}); ok {
		return resultMap, nil
	}
	return nil, fmt.Errorf("返回格式异常")
}

// UpdateData 更新数据
// [Audit Fix]: 使用 PUT 方法 + /update 路径，并正确构造 filter
// API: PUT /v1/model/prod/{modelName}/update
func (c *CloudBaseClient) UpdateData(modelName, id string, data interface{}) error {
	path := fmt.Sprintf("/v1/model/prod/%s/update", modelName)

	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"where": map[string]interface{}{
				"_id": map[string]interface{}{
					"$eq": id,
				},
			},
		},
		"data": data,
	}

	_, err := c.Request("PUT", path, payload, nil)
	return err
}

// DeleteData 删除数据
// API: POST /v1/model/prod/{modelName}/delete
func (c *CloudBaseClient) DeleteData(modelName, id string) error {
	path := fmt.Sprintf("/v1/model/prod/%s/delete", modelName)

	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"where": map[string]interface{}{
				"_id": map[string]interface{}{
					"$eq": id,
				},
			},
		},
	}

	_, err := c.Request("POST", path, payload, nil)
	return err
}

// GetDetail 获取单条详情
// 复用 list 接口，查询 _id
func (c *CloudBaseClient) GetDetail(modelName, id string) (map[string]interface{}, error) {
	// 构造标准 filter
	filter := map[string]interface{}{
		"where": map[string]interface{}{
			"_id": map[string]interface{}{
				"$eq": id,
			},
		},
	}

	// 复用 ListData 逻辑
	result, err := c.ListData(modelName, filter, 1, 1)
	if err != nil {
		return nil, err
	}

	if dataMap, ok := result["data"].(map[string]interface{}); ok {
		if records, ok := dataMap["records"].([]interface{}); ok && len(records) > 0 {
			return records[0].(map[string]interface{}), nil
		}
	}
	return nil, fmt.Errorf("未找到记录")
}
