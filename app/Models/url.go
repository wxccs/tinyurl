package Models

import "time"

type Url struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ShortCode   string    `gorm:"uniqueIndex;size:20;not null" json:"short_code"`
	OriginalUrl string    `gorm:"type:varchar(2048);not null" json:"original_url"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
