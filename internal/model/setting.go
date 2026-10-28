package model

// Setting 系统设置表（简单 key-value）
type Setting struct {
	Key   string `gorm:"primaryKey;size:100" json:"key"`
	Value string `gorm:"size:500" json:"value"`
}
