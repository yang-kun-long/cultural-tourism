// File: models/poi.go
package models

// POI 类型枚举 (对应 PRD 7.1)
const (
	POITypeScenic = "scenic" // 景点
	POITypeFood   = "food"   // 饭店
	POITypeHotel  = "hotel"  // 酒店
	POITypeBooth  = "booth"  // 旅拍机
)

// POI 点位数据模型
type POI struct {
	ID        string   `json:"_id,omitempty"`     // TCB 自动 ID
	OpenID    string   `json:"_openid,omitempty"` // [Audit Fix] 系统字段：记录创建者
	Name      string   `json:"name"`              // 名称
	Type      string   `json:"type"`              // 类型
	RegionID  string   `json:"region_id"`         // 所属区域
	Latitude  float64  `json:"latitude"`          // 纬度
	Longitude float64  `json:"longitude"`         // 经度
	Images    []string `json:"images"`            // 轮播图
	Desc      string   `json:"desc"`              // 简介
	Address   string   `json:"address"`           // 地址
	Phone     string   `json:"phone"`             // 电话
	OpenTime  string   `json:"open_time"`         // 营业时间
	Status    int      `json:"status"`            // 状态 1:启用 0:禁用
	CreatedAt string   `json:"created_at"`        // 业务创建时间
	UpdatedAt string   `json:"updated_at"`        // 业务更新时间

	// [Audit Fix] 扩展字段：仅用于返回给前端，不存库
	Distance float64 `json:"_distance,omitempty"` // 距离(米)
}

// POIQuery 列表筛选参数
type POIQuery struct {
	RegionID string `form:"region_id"`
	Type     string `form:"type"`

	// [Audit Fix] LBS 核心参数：用户的当前位置
	UserLat float64 `form:"lat"` // 用户纬度
	UserLng float64 `form:"lng"` // 用户经度

	Page int `form:"page,default=1"`
	Size int `form:"size,default=10"`
}
