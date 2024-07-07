package def

// QuizMode 决定了如何处理 familiarity 为 0 的项
type QuizMode string

const (
	// QuizModeAlwaysNew - always_new: 优先学新 (总是先选 familiarity 为 0 的)
	QuizModeAlwaysNew QuizMode = "always_new" // 总是优先复习，但从熟练度低的开始
	// QuizModeAlwaysOld - always_old: 优先复习 (总是先选 familiarity 不为 0 的)
	QuizModeAlwaysOld QuizMode = "always_old" //
	// QuizModeBalance - balance: 平衡随机 (先选 familiarity 不为 0 的，但一定概率会选中 familiarity 为 0 的)
	QuizModeBalance QuizMode = "balance" //
	// QuizModeThreshold - threshold: 阀门 (根据 familiarity 不为 0 的数量决定，familiarity 不为 0 的项超过一定数量时，优先复习 familiarity 为 0 的项)
	QuizModeThreshold QuizMode = "threshold" //
	// QuizModeDynamic - dynamic: 动态调权 (根据当天加入 familiarity 为 0 的数量决定，保证至少学习一定数量的新项)
	QuizModeDynamic QuizMode = "dynamic" //

)
