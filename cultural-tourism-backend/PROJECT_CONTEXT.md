# Project Master Context: Digital Cultural Tourism Backend

## 1. 维护协议 (Protocol)

> **给 AI 的指令 (CRITICAL)**：
>
> 1. 本文档是项目的“唯一事实来源 (Source of Truth)”。
> 2. **核心红线**：
>    - **禁止臆造底层逻辑**：`tcb/client.go` 必须保持通用性（支持 Map 传参构造 Filter），**严禁**在底层 SDK 中硬编码业务逻辑（如 `$eq`）。
>    - **禁止使用 Mongo Driver**：必须且只能使用 `tcb/client.go` 封装的 HTTP API。
>    - **生产级标准**：所有错误处理必须规范，配置缺失必须 Panic，禁止“能跑就行”的将就心态。
> 3. **操作规范**：
>    - **更新(Update)**: 必须用 `PUT` + `/update` + `filter` + `data` (data 内含_id)。
>    - **删除(Delete)**: 必须用 `POST` + `/delete` + `filter`。
> 4. **业务边界**：严格遵循 PRD，不准臆造需求（如：不做支付、不做在线修图）。

## 2. 项目概况 (Overview)

- **项目名称**: 数字文旅一体化小程序平台
- **核心价值**: 内容发现(LBS) -> 参与上传(UGC) -> 线下取图(旅拍机) -> 互动评价
- **技术架构**: Go (Gin) + TCB HTTP API + Swagger + Air
- **业务边界**:
  - **不做**: 站内交易/支付 (只做小程序跳转导流)
  - **不做**: 在线修图/AI 合成
  - **不做**: 跨平台跳转 (只跳微信小程序)

## 3. 架构规范 (Constraints)

1. **数据交互**:
    - 所有数据库操作统一封装在 `tcb/client.go`。
    - `ListData` 方法必须支持 `map[string]interface{}` 类型的通用 Filter，以支持 `$eq`, `$regex`, `$gt` 等复杂查询。
    - 列表接口默认支持分页 (`page`, `size`)。
2. **API 设计**:
    - RESTful 路由 (`/api/resource/:id`)。
    - 必须提供 Swagger 注释。

## 4. 开发进度与路线图 (Roadmap based on PRD)

### Phase 1: 基础设施 [已完成 ✅]

- [x] TCB HTTP Client 核心封装 (通用 Filter 模式) `[Audit: Passed]`
- [x] Swagger 文档集成 `[Audit: Passed]`
- [x] 统一错误处理 `[Audit: Passed]`

### Phase 2: 基础资源管理 [已完成 ✅]

- [x] **区域管理 (Regions)** (PRD 6.2, 7.1) `[Audit: Passed]`
  - [x] 模型: Name, Sort, Status
  - [x] 接口: 增删改查(RESTful)

### Phase 3: LBS 线下体验 [已完成 ✅]

- [x] **景点点位 (POIs)** (PRD 5.2, 7.1) `[Audit: Passed]`
  - [x] **模型定义**:
    - `type`: 枚举 (scenic/food/hotel/booth)
    - `location`: 经纬度 (Go层计算距离)
    - `info`: 轮播图, 简介, 营业时间, 电话
    - `system`: 包含 _openid, owner 等完整系统字段
  - [x] **接口开发**:
    - [x] 增删改查 (CRUD 核心闭环)
    - [x] 列表查询 (支持分页与区域/类型筛选)
    - [x] LBS 距离计算 (Haversine Algorithm)

### Phase 4: UGC 旅拍社区 [进行中 🔄]

> *对应 PRD 5.1, 6.3, 7.1*

- [x] **旅拍主题 (Themes)** `[Audit: Passed]`
  - [x] **模型**: Name, Cover, RegionID, Sort, Status, Desc
  - [x] **接口**:
    - [x] 主题列表: 支持按 **“区域优先”** (region_id) 筛选
    - [x] 主题详情: 展示封面、简介
- [ ] **照片管理 (Photos)** `[Audit: Pending]`
  - [ ] **模型**: ThemeID, UserID, URL, Status(待审/通过/下架)
  - [ ] **接口**:
    - [ ] 照片上传 (仅手机相册)
    - [ ] 瀑布流列表 (关联 Theme)
    - [ ] 审核状态流转 (UGC 核心: 待审->通过/拒绝)

### Phase 5: 互动与导流

> *对应 PRD 5.4, 5.5, 6.1, 7.1*

- [ ] **评论互动 (Comments)** `[Audit: Pending]`
  - [ ] **模型**: POI_ID, Content, UserID, Status, ParentID(盖楼)
  - [ ] **接口**: 发布评论(默认待审)
- [ ] **商品导流 (Products)** `[Audit: Pending]`
  - [ ] **模型**: Name, Image, Price, JumpAppID, JumpPath
  - [ ] **接口**: 列表(区域优先), 详情(无支付直接跳转)

### Phase 6: 用户资产与旅拍机

> *对应 PRD 5.3, 5.6, 7.1*

- [ ] **旅拍机联动 (Booth)** `[Audit: Pending]`
  - [ ] **接口**: 扫码取图 (关联 BoothOrder/Scan)
- [ ] **电子相册 (Album)** `[Audit: Pending]`
  - [ ] **逻辑**: 聚合“旅拍机照片”+“线上上传照片”
- [ ] **收藏体系 (Favorites)** `[Audit: Pending]`
  - [ ] **接口**: 收藏/取消收藏 (对象: Theme/POI/Product)

### Phase 7: 后台管理与审核 (Admin API)

> *对应 PRD 6.1, 6.4*

- [ ] **内容审核**: 照片/评论的 批量通过/拒绝 `[Audit: Pending]`
- [ ] **数据统计**: 扫码量, UGC 上传量 `[Audit: Pending]`

## 5. 核心数据模型 (Schema Snapshot from PRD 7.1)

### `regions` (已上线)

```json
{ "name": "string", "sort": "number", "status": "number", "_id": "..." }

```

### `pois` (已上线)

```json
{
  "name": "string",
  "type": "string (enum: scenic, food, hotel, booth)",
  "region_id": "string (ref: regions._id)",
  "latitude": "number",
  "longitude": "number",
  "images": "array<string>",
  "desc": "string",
  "address": "string",
  "phone": "string",
  "open_time": "string",
  "status": "number",
  "_openid": "string (system)",
  "_id": "..."
}

```

### `themes` (Collection Name: `theme`) (已上线)

```json
{
  "name": "string",
  "cover": "string",
  "desc": "string",
  "region_id": "string (ref: regions._id)",
  "sort": "number",
  "status": "number",
  "_openid": "string (system)",
  "_id": "..."
}

```

## 6. 当前项目目录结构

```text
/
├── config
│   └── config.go
├── controllers
│   ├── poi_controller.go
│   ├── region_controller.go
│   └── theme_controller.go  # [新增]
├── database
│   └── db.go
├── model-json
│   ├── regions_model.json
│   ├── pois_model.json
│   └── themes_model.json    # [新增]
├── models
│   ├── poi.go
│   ├── region.go
│   └── theme.go             # [新增]
├── routes
│   └── router.go
├── tcb
│   └── client.go
├── go.mod
├── main.go
└── PROJECT_CONTEXT.md

```
