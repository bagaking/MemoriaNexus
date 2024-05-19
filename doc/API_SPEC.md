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
    - POST /api/v1/book：创建册子
    - GET /api/v1/book：获取册子列表
    - GET /api/v1/book/:id：获取册子详情
    - PUT /api/v1/book/:id：更新册子信息
    - DELETE /api/v1/book/:id：删除册子
- 学习材料管理
    - POST /api/v1/item：创建学习材料
    - GET /api/v1/item：获取学习材料列表
    - GET /api/v1/item/:id：获取学习材料详情
    - PUT /api/v1/item/:id：更新学习材料信息
    - DELETE /api/v1/item/:id：删除学习材料
- 复习计划管理
    - GET /api/v1/dungeon/schedules：获取复习计划列表
    - GET /api/v1/dungeon/schedules/:id：获取复习计划详情
    - POST /api/v1/dungeon/schedules：创建复习计划
    - PUT /api/v1/dungeon/schedules/:id：更新复习计划
    - DELETE /api/v1/dungeon/schedules/:id：删除复习计划
    - GET /api/v1/dungeon/instances：获取复习实例
    - GET /api/v1/dungeon/instances/:id：获取复习实例详情
- NFT管理
    - GET /api/v1/nft：获取用户 NFT
    - GET /api/v1/nft/:id：获取 NFT 详情
    - GET /api/v1/nft/shop：查看商店
    - POST /api/v1/nft/draw：抽卡
- NFT交易管理
    - GET /api/v1/trade：获取市场交易
    - POST /api/v1/trade：创建交易
    - GET /api/v1/trade/:id：获取交易详情
    - DELETE /api/v1/trade/:id：取消交易
    - POST /api/v1/trade/:id/purchase：购买交易
- 成就系统
    - GET /api/v1/achievements：获取所有成就
    - GET /api/v1/achievements/:id：获取成就详情
- 运营管理
    - GET /api/v1/operation/task：当前任务获取
    - GET /api/v1/operation/task/:id：获取任务详情
    - GET /api/v1/operation/activity：获取活动
    - GET /api/v1/operation/activity/:id：获取活动详情