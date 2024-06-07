### API 设计 (实现见 src/module)

#### 用户账户服务

- **GET /profile/me**：获取用户个人资料（无需参数）
- **PUT /profile/me**：更新用户个人资料（body 支持用户的详细信息更新）
- **GET /profile/points**：获取用户积分（金币，钻石）（无需参数）
- **GET /profile/settings/memorization**：获取用户记忆设置（无需参数）
- **PUT /profile/settings/memorization**：更新用户记忆设置（body 支持记忆设置的详细信息更新）
- **GET /profile/settings/advance**：获取用户高级设置（无需参数）
- **PUT /profile/settings/advance**：更新用户高级设置（body 支持高级设置的详细信息更新）

#### 系统操作

- **GET /system/notifications**：获取所有通知（无需参数）
- **POST /system/notifications/markAsRead**：标记通知为已读（body 支持通知 ID 列表）
- **GET /system/announcements**：获取所有公告（无需参数）
- **POST /system/announcements/markAsRead**：标记公告为已读（body 支持公告 ID 列表）
- **GET /system/configs**：获取全局配置（无需参数）

#### 册子管理

- **POST /books**：创建册子（body 支持册子的详细信息）
- **GET /books**：获取册子列表（query 支持分页参数 page 和 limit）
- **GET /books/:id**：获取册子详情
- **PUT /books/:id**：更新册子信息（body 支持册子的详细信息更新）
- **DELETE /books/:id**：删除册子

#### 学习材料管理

- **POST /items**：创建学习材料（body 支持学习材料的详细信息）
- **GET /items**：获取学习材料列表（query 支持分页参数 page 和 limit，以及可选的 book_id 和 type 过滤）
- **GET /items/:id**：获取学习材料详情
- **PUT /items/:id**：更新学习材料信息（body 支持学习材料的详细信息更新）
- **DELETE /items/:id**：删除学习材料

#### 复习计划管理

- **POST /dungeon/dungeons**：创建复习计划（body 支持复习计划的详细信息）
- **GET /dungeon/dungeons**：获取复习计划列表（无需参数）
- **GET /dungeon/dungeons/:id**：获取复习计划详情
- **PUT /dungeon/dungeons/:id**：更新复习计划（body 支持复习计划的详细信息更新）
- **DELETE /dungeon/dungeons/:id**：删除复习计划

- **POST /dungeon/dungeons/:id/books**：添加复习计划的 Books（body 支持书籍 ID 列表）
- **POST /dungeon/dungeons/:id/items**：添加复习计划的 Items（body 支持学习材料 ID 列表）
- **POST /dungeon/dungeons/:id/tags**：添加复习计划的 Tags（body 支持标签 ID 列表）
- **GET /dungeon/dungeons/:id/books**：获取复习计划的 Books
- **GET /dungeon/dungeons/:id/items**：获取复习计划的 Items
- **GET /dungeon/dungeons/:id/tags**：获取复习计划的 Tags
- **DELETE /dungeon/dungeons/:id/books**：删除复习计划的 Books（body 支持书籍 ID 列表）
- **DELETE /dungeon/dungeons/:id/items**：删除复习计划的 Items（body 支持学习材料 ID 列表）
- **DELETE /dungeon/dungeons/:id/tags**：删除复习计划的 Tags（body 支持标签 ID 列表）

- **GET /dungeon/campaigns/:id/monsters**：获取战役副本的所有 Monsters（query 支持排序字段 sort_by 和分页参数 offset 和 limit）
- **GET /dungeon/campaigns/:id/practice**：获取战役副本的后 n 个 Monsters（query 支持获取数量 count 和排序字段 sort_by）
- **POST /dungeon/campaigns/:id/submit**：上报战役副本的 Monster 结果（body 支持结果数据）
- **GET /dungeon/campaigns/:id/conclusion/today**：获取战役副本的结果 (当日)

- **GET /dungeon/endless/:id/monsters**：获取无限副本的所有 Monsters 及其关联的 Items, Books, Tags（query 支持排序字段 sort_by 和分页参数 offset 和 limit）

#### NFT管理
- **GET /nft/nfts**：获取用户 NFT（无需参数）
- **GET /nft/nfts/:id**：获取 NFT 详情
- **POST /nft/draw_card**：以抽卡的方式创建 NFT（无需参数）
- **POST /nft/transfer**：赠予 NFT（body 支持接收者信息和 NFT ID）
- **GET /nft/shops**：查看所有商店（无需参数）
- **GET /nft/shops/:id**：查看某个商店
- **GET /nft/trades**：获取市场交易对（无需参数）
- **POST /nft/trades**：创建交易对（body 支持交易对的详细信息）
- **GET /nft/trades/:id**：获取交易详情
- **DELETE /nft/trades/:id**：取消交易
- **POST /nft/trades/:id/buy**：创建购买订单

#### 成就系统

- **GET /achievements**：获取所有成就（无需参数）
- **GET /achievements/:id**：获取成就详情

#### 运营管理

- **GET /operation/task**：获取当前任务（无需参数）
- **GET /operation/task/:id**：获取任务详情
- **GET /operation/activity**：获取活动（无需参数）
- **GET /operation/activity/:id**：获取活动详情