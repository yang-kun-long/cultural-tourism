// File: models/poi.go
package models

// POI 类型枚举 (对应 PRD 7.1 类型: 景点/饭店/酒店/旅拍机)
const (
	POITypeScenic = "scenic" // 景点
	POITypeFood   = "food"   // 饭店
	POITypeHotel  = "hotel"  // 酒店
	POITypeBooth  = "booth"  // 旅拍机
)

// POI 点位数据模型
type POI struct {
	ID        string   `json:"_id,omitempty"` // TCB 自动生成的 ID
	Name      string   `json:"name"`          // 名称
	Type      string   `json:"type"`          // 类型 (枚举值)
	RegionID  string   `json:"region_id"`     // 所属区域 ID
	Latitude  float64  `json:"latitude"`      // 纬度
	Longitude float64  `json:"longitude"`     // 经度
	Images    []string `json:"images"`        // 图片 URL 列表
	Desc      string   `json:"desc"`          // 简介
	Address   string   `json:"address"`       // 详细地址
	Phone     string   `json:"phone"`         // 联系电话
	OpenTime  string   `json:"open_time"`     // 营业时间
	Status    int      `json:"status"`        // 状态 1:启用 0:禁用
	CreatedAt string   `json:"created_at"`    // 创建时间 (ISO string)
	UpdatedAt string   `json:"updated_at"`    // 更新时间 (ISO string)
}

// POIQuery 用于列表筛选
type POIQuery struct {
	RegionID string `form:"region_id"`
	Type     string `form:"type"`
	Page     int    `form:"page,default=1"`
	Size     int    `form:"size,default=10"`
}
