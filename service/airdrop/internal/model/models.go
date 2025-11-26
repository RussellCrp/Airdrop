package model

import "time"

// Users 表
type User struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	WalletAddr string    `gorm:"column:wallet_addr;size:64;uniqueIndex;not null"`
	Nickname   string    `gorm:"column:nickname;size:64"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// Tasks 表
type Task struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string    `gorm:"column:code;size:64;uniqueIndex;not null"`
	Name        string    `gorm:"column:name;size:128;not null"`
	Description string    `gorm:"column:description;type:text"`
	ScoreWeight int       `gorm:"column:score_weight;not null"`
	Enabled     bool      `gorm:"column:enabled;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Task) TableName() string {
	return "tasks"
}

// UserTasks 表
type UserTask struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      uint64    `gorm:"column:user_id;index;not null"`
	TaskID      uint64    `gorm:"column:task_id;index;not null"`
	Status      int8      `gorm:"column:status;not null"` // 0 未完成 1 已完成
	ExtraData   string    `gorm:"column:extra_data;type:json"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserTask) TableName() string {
	return "user_tasks"
}

// UserScores 表
type UserScore struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     uint64    `gorm:"column:user_id;uniqueIndex;not null"`
	TotalScore int       `gorm:"column:total_score;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserScore) TableName() string {
	return "user_scores"
}

// ScoreTiers 表
type ScoreTier struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	MinScore      int       `gorm:"column:min_score;not null"`
	MaxScore      int       `gorm:"column:max_score;not null"`
	AirdropAmount string    `gorm:"column:airdrop_amount;type:decimal(38,0);not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ScoreTier) TableName() string {
	return "score_tiers"
}

// AirdropSnapshots 表
type AirdropSnapshot struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;size:128;not null"`
	MerkleRoot  string    `gorm:"column:merkle_root;size:66"`
	TokenAddr   string    `gorm:"column:token_address;size:64;not null"`
	TotalUsers  int       `gorm:"column:total_users;not null"`
	TotalAmount string    `gorm:"column:total_amount;type:decimal(38,0);not null"`
	Status      int8      `gorm:"column:status;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AirdropSnapshot) TableName() string {
	return "airdrop_snapshots"
}

// AirdropSnapshotItems 表
type AirdropSnapshotItem struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	SnapshotID uint64    `gorm:"column:snapshot_id;index;not null"`
	UserID     uint64    `gorm:"column:user_id;index;not null"`
	WalletAddr string    `gorm:"column:wallet_addr;size:64;not null"`
	Idx        int       `gorm:"column:idx;not null"`
	Score      int       `gorm:"column:score;not null"`
	Amount     string    `gorm:"column:amount;type:decimal(38,0);not null"`
	LeafHash   string    `gorm:"column:leaf_hash;size:66"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (AirdropSnapshotItem) TableName() string {
	return "airdrop_snapshot_items"
}


