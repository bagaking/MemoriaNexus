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

### Item 的过程属性和 Monster

**与用户无关的属性直接记录在 item 中，包括：**
难度（Difficulty）: 使用十六进制分级方式表示，从 Novice 到 Master，每级别有不同难度。
合理的 Difficulty 应该会因人而异，所以按照 “小白-容易” “小白-普通” “小白-困难” “业余-容易” ... “职业-” ... "专家-" ... "大师-" ... 来分类
重要程度（Importance）: 表示学习材料在其影响领域范围的经典或重要程度，和 Difficulty 类似的是，Importance 包含了范围 + 奠基程度，比如 “Domain-General、Domain-Key ... Area-MasterPiece ... Global-Essential ...”。

**与用户相关的属性，被定义为 monster 主要是**
熟练度（Familiarity）: 记录用户对学习材料的熟悉程度, 取值 0-100，表示熟练度的百分位。

**与用户和 dungeon 都相关的属性，被定义为 dungeon_monster**
包含了在某个具体 dungeon 中的攻略信息

### 为什么复习计划叫做 Dungeon

Dungeon 是把一个复习计划包装成了游戏里打怪升级副本的概念，完成卡片记忆就是挑战怪物的过程，综合测验就是大 boss。

Dungeon 共有三个大类
- 战役地牢（Campaign Dungeon），由用户配置，可以指定 item，也可以通过 Book 关联或者 Tag 进行导入，配置后系统对根据 Dungeon 中 item 的具体数量、难度、熟练度等情况来分配关卡和奖励。当 Book 或 Tag 对应的内容增加和减少时，item 不会变化。
- 无限地牢（Endless Dungeon），由用户配置，可以和 Book、Tag 关联，但不会导入 item，当 Book 和 Tag 增减时，计划会发生相应的变化
- 副本地牢（Instance Dungeon），则是会由系统主动创建，会有很多种类型，比如限时副本，突发副本等，其包装的概念是怪物可能会随时出现， 对应到记忆中，就是指之前记熟过的内容会间隔重复的出现。

他们具备不同的特点，分别对应了不同的需清洗类型
- Campaign Dungeon 和 item 由于是导入和配置的，因此通常是成套的，比如都属于一个学科，或者面向一个考试，适合用于新学一项知识。
- Endless Dungeon 中的 item 是变化的，因此通常是用于聚焦于某个具体的话题
- 而Instances 中的则是一些更加开放的规则，比如纯粹随机抽查，按照 tag 抽查，AI 结合最近重大事件进行推荐等。相比 Campaign 和 Endless，它面向全局，目标是长期记忆。

### 各类 Dungeon 的 Monster 插入逻辑

1. Campaign Dungeon（计划类的 Dungeon）: 
   - 固定的 dungeon, 用户通过界面导入 Book、Tag 或者 Item，导入时就创建 DungeonMonster 记录，后续 Book 和 Tag 和 item 的映射发生变化时，Dungeon 对应的 DungeonMonster 不会发生变化。
2. Endless Dungeon（计划类的 Dungeon）:
   - 无限 dungeon, 用户通过界面关联 Book、Tag 或者 Item。
   - 只有关联 Item 会创建 DungeonMonster 记录，Book、Tag 不会。但每次查询都能查询到 DungeonMonster 记录的全集（通过 Book Tag 等逻辑关联）。 
   - 因此，后续 Book 和 Tag 和 item 的映射发生变化时，Endless Dungeon 总能查到最新的 DungeonMonster。表现上 DungeonMonster 会出现和离开。
3. Instances Dungeon（即时类的 Dungeon） 
   - 系统会自动创建突发复习任务，让用户在不定时复习间隔记忆内容，自身不会创建新的 DungeonMonster，只会根据现有的 DungeonMonster 进行组合。

创建 DungeonMonster 时不一定要同时创建 Monster，Monster 中没有记录默认熟悉度为 0 即可

> - AddXXXToCampaignDungeon 时应该要创建 DungeonMonster （创建 DungeonMonster 时不一定要同时创建 Monster，Monster 中没有记录默认熟悉度为 0 即可），在 Get 时只查询现有记录
> - AddXXXToEndlessDungeon 则是只有在 addItem 创建 DungeonMonster，而 book、tag 则只是创建 DungeonBooks、DungeonItems 的记录，在 Get 时，在 Get 时要根据 Book、Tag 的关系去组装

### Dungeon 复习流程的逻辑

Dungeon 的基本复习逻辑，是把可以复习当做对一个 Monster 的攻击。
Monster 其实就是 Item，但是注册阶段 dungeon 也可以注册 item，所以在复习阶段转而使用 Monster 的概念来代替。
Monster 即是一个 Dungeon 可以找到的所有 Item，不光是直接注册的 Item，也包含从 Book 和 Tag 能级联出来且去重后的 Item。

复习包括几个方法
- 获取当前 Dungeon 所有的 Monster
- 获取后 n 个要战斗（复习）的 Monster
- 上报复习的结果：完全陌生，有点印象，牢记 （更具复习的情况，Monster 的显影程度会增加，也就是随着记忆次数变多，知道这个 Monster 的强度）
- 获取复习的结果：共复习了多少张卡片（怪物数量），今天挑战的综合难度（怪物状态）等，这里可以自由发挥
