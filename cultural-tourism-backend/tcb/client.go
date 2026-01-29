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

// 全局客户端实例
var Client *CloudBaseClient

type CloudBaseClient struct {
	EnvID       string
	AccessToken string
	BaseURL     string
	HTTPClient  *http.Client
}

// Init 替换原来的 NewCloudBaseClient，改为全局初始化
func Init() {
	// 加载 .env 文件 (本地开发用)
	_ = godotenv.Load()

	envID := os.Getenv("CLOUDBASE_ENV_ID")
	// 注意：本地开发时，你需要手动获取 AccessToken 填入 .env
	// 部署到云托管后，通常可以通过元数据服务获取，或者继续使用持久化的配置
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
	fmt.Println("✅ 云开发 HTTP 客户端已初始化")
}

// Request 通用请求方法
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
	req.Header.Set("Authorization", "Bearer "+c.AccessToken) // 鉴权核心

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 调试日志：看看腾讯云返回了什么
	// fmt.Printf("Debug Response: %s\n", string(bodyBytes))

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

// CreateData 封装：新增数据 (对应数据库 Insert)
// 使用数据模型 API: /v1/model/{envType}/{modelName}/create
func (c *CloudBaseClient) CreateData(modelName string, data interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v1/model/prod/%s/create", modelName) // 默认使用 prod 环境

	// 构造成数据模型要求的格式
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

// ListData 封装：查询数据列表 (支持筛选)
// 对应 API: /v1/model/{envType}/{modelName}/list
// 修改：增加了 filter 参数 map[string]interface{}
func (c *CloudBaseClient) ListData(modelName string, filter map[string]interface{}, page, size int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v1/model/prod/%s/list", modelName)

	payload := map[string]interface{}{
		"pageNumber": page,
		"pageSize":   size,
		"getCount":   true,
	}

	// 如果有筛选条件，并且不为空，则构造 filter 结构
	// TCB 的 list 接口要求 filter 格式为: { "where": { "字段": { "$eq": 值 } } }
	if len(filter) > 0 {
		where := make(map[string]interface{})
		for k, v := range filter {
			// 这里做一个简单的处理，默认所有筛选都是“等于” ($eq)
			// 如果需要更复杂的查询，需要在 controller 层构造完整的 TCB 查询语法
			where[k] = map[string]interface{}{
				"$eq": v,
			}
		}
		payload["filter"] = map[string]interface{}{
			"where": where,
		}
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

// DeleteData 封装：删除数据
// 对应 API: /v1/model/{envType}/{modelName}/delete
// DeleteData 封装：删除数据
func (c *CloudBaseClient) DeleteData(modelName, id string) error {
	path := fmt.Sprintf("/v1/model/prod/%s/delete", modelName)

	// 【修正】按照报错提示，构建严格的 filter 结构
	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"where": map[string]interface{}{
				"_id": map[string]interface{}{
					"$eq": id, // 意思是: delete where _id == id
				},
			},
		},
	}

	_, err := c.Request("POST", path, payload, nil)
	return err
}

// UpdateData 封装：更新数据 (最终修正版)
// 结合了之前的报错信息：
// 1. 必须用 /update 路径 (wedaUpdate 报 404)
// 2. 必须用 PUT 方法 (POST 报 not supported)
// 3. 必须带 filter 参数 (否则报 filter不能为空)
func (c *CloudBaseClient) UpdateData(modelName, id string, data interface{}) error {
	path := fmt.Sprintf("/v1/model/prod/%s/update", modelName)

	payload := map[string]interface{}{
		// 1. 必须提供筛选条件
		"filter": map[string]interface{}{
			"where": map[string]interface{}{
				"_id": map[string]interface{}{
					"$eq": id,
				},
			},
		},
		// 2. 提供要更新的数据
		"data": data,
	}

	// 3. 必须使用 PUT 方法
	_, err := c.Request("PUT", path, payload, nil)
	return err
}

// GetDetail 封装：查询单条详情
// 实际上复用 list 接口，通过 ID 过滤来实现
func (c *CloudBaseClient) GetDetail(modelName, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/v1/model/prod/%s/list", modelName)

	payload := map[string]interface{}{
		"filter": map[string]interface{}{
			"where": map[string]interface{}{
				"_id": map[string]interface{}{
					"$eq": id,
				},
			},
		},
		"pageNumber": 1,
		"pageSize":   1,
	}

	result, err := c.Request("POST", path, payload, nil)
	if err != nil {
		return nil, err
	}

	// 解析返回结果，提取第一条记录
	if dataMap, ok := result.(map[string]interface{}); ok {
		if records, ok := dataMap["data"].(map[string]interface{})["records"].([]interface{}); ok && len(records) > 0 {
			return records[0].(map[string]interface{}), nil
		}
	}

	return nil, fmt.Errorf("未找到该记录")
}
