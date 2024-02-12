package model

import (
	"context"
	"gorm.io/gorm"
)

// Repo 提供用户注册的服务
type Repo struct {
	DB *gorm.DB
}

// NewRepo 创建一个新的RegisterService实例
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

// UpdateProfileSettings 创建或者更新用户的昵称和头像URL
func (repo *Repo) UpdateProfileSettings(ctx context.Context, userID uint, nickname string, avatarURL string) error {

	updateData := map[string]any{}

	// Only update fields if they are non-empty to avoid overriding with null.
	if nickname != "" {
		updateData["nickname"] = nickname
	}
	if avatarURL != "" {
		updateData["avatar_url"] = avatarURL
	}

	if len(updateData) == 0 {
		return nil //todo: return param error
	}

	// Start a transaction as a precaution for concurrent database access
	tx := repo.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Attempt to retrieve the user's profile from the database
	var profile Profile
	if err := tx.Where(Profile{UserID: userID}).Attrs(Profile{Nickname: nickname, AvatarURL: avatarURL}).FirstOrCreate(&profile).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Update the profile with new values where they exist
	needsUpdate := false // 在这里添加一个用于确保仅当有实际更改时才运行 Updates 的额外检查
	for key, value := range updateData {
		if (key == "nickname" && profile.Nickname != value) || (key == "avatar_url" && profile.AvatarURL != value) {
			needsUpdate = true
			break
		}
	}

	// 如果确有变化，则进行更新
	if needsUpdate {
		if err := tx.Model(&profile).Updates(updateData).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}
