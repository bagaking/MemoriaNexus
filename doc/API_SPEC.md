# 功能描述

Memoria Nexus 利用艾宾浩斯遗忘曲线（间隔重复法）优化学习记忆过程。
它智能地为学习材料安排复习时间，确保你能在长期内以最小的努力保留知识。

- 自定义学习材料：支持用户自行创建、编辑、管理学习材料，包含笔记和抽认卡。学习材料可以按照 “册” 进行组织，并且支持标签和双向链接功能。
- 智能复习调度：利用艾宾浩斯的遗忘曲线确定最佳的复习时间，以加强记忆保留。
- 用户进度分析：通过详细的报告和分析，了解你的学习进步。
- 游戏化激励：完成记忆时会根据任务难度发放积分，积分可以抽卡获取 NFT，NFT 可以用于各种特效，比如
    - 追加类：增加完成任务时的积分、额外获得抽卡机会
    - 挂机类：获得挂机收益
    - 通道类：解锁收益更高的复习计划，解锁开屏任务等
- 多平台支持：我们的响应式网页设计兼容桌面和移动平台，随时随地访问你的学习材料。
- 通知和提醒：通过多设备的及时通知，永不错过复习。
- AI 功能：支持 AI 优化卡片内容，AI 根据关联重新组织卡片内容量，AI 根据卡片推荐记忆技巧等。

> DB: MySQL, Cache: Redis, AuthZ: RangeIAM

## 一些细节

### Book Item Tag 之间的关系

Item 和 Book 是 n 对 m 关系，但是他们都同时只会属于一个 user。
item 和 book 都支持添加任意多个 tag。
tag 是全局的，可以对应到任意多个 book 或者 item。

# API 设计 (实现见 src/module) 

- 用户账户服务
    - GET /api/v1/profile/me：获取用户个人资料
    - PUT /api/v1/profile/me：更新用户个人资料
    - GET /api/v1/profile/points：获取用户 points (积分，金币，钻石)
    - GET /api/v1/profile/settings：获取用户设置
    - PUT /api/v1/profile/settings：设置用户配置信息
- 系统操作
    - GET /api/v1/system/notifications：获取所有通知
    - POST /api/v1/system/notifications/markAsRead：标记通知为已读
    - GET /api/v1/system/announcements：获取所有公告
    - POST /api/v1/system/announcements/markAsRead：标记公告为已读
    - GET /api/v1/system/configs：获取全局配置
- 册子管理
    - POST /api/v1/books：创建册子
    - GET /api/v1/books：获取册子列表
    - GET /api/v1/books/:id：获取册子详情
    - PUT /api/v1/books/:id：更新册子信息
    - DELETE /api/v1/books/:id：删除册子
- 学习材料管理
    - POST /api/v1/items：创建学习材料
    - GET /api/v1/items：获取学习材料列表
    - GET /api/v1/items/:id：获取学习材料详情
    - PUT /api/v1/items/:id：更新学习材料信息
    - DELETE /api/v1/items/:id：删除学习材料
- 复习计划管理
    - GET /api/v1/dungeon/schedules：获取复习计划列表
    - GET /api/v1/dungeon/schedules/:id：获取复习计划详情
    - POST /api/v1/dungeon/schedules：创建复习计划
    - PUT /api/v1/dungeon/schedules/:id：更新复习计划
    - DELETE /api/v1/dungeon/schedules/:id：删除复习计划
    - GET /api/v1/dungeon/instances：获取复习即时副本
    - GET /api/v1/dungeon/instances/:id：获取复习即时副本详情
- NFT管理
    - GET /api/v1/nft/nfts：获取用户 NFT
    - GET /api/v1/nft/nfts/:id：获取 NFT 详情
    - POST /api/v1/nft/draw_card：以抽卡的方式创建 nft
    - GET /api/v1/nft/shops：查看所有商店
    - GET /api/v1/nft/shops/:id：查看某个商店
    - POST /api/v1/nft/transfer：赠予
- NFT交易管理
    - GET /api/v1/nft/trades：获取市场交易对
    - POST /api/v1/nft/trades：创建交易对
    - GET /api/v1/nft/trades/:id：获取交易详情
    - DELETE /api/v1/nft/trades/:id：取消交易
    - POST /api/v1/nft/trades/:id/buy：创建购买订单 (会直接生效、因此就是购买)
- 成就系统
    - GET /api/v1/achievements：获取所有成就
    - GET /api/v1/achievements/:id：获取成就详情
- 运营管理
    - GET /api/v1/operation/tasks：当前任务获取
    - GET /api/v1/operation/tasks/:id：获取任务详情
    - GET /api/v1/operation/activities：获取活动
    - GET /api/v1/operation/activities/:id：获取活动详情