package model

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement"`
	Wallet        string     `gorm:"size:64;uniqueIndex"`
	LoginStreak   uint32     `gorm:"column:login_streak"`
	LoginDays     uint64     `gorm:"column:login_days"`
	LastLoginAt   *time.Time `gorm:"column:last_login_at"`
	PointsBalance int64      `gorm:"column:points_balance"`
	FrozenPoints  int64      `gorm:"column:frozen_points"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type Task struct {
	ID          uint64         `gorm:"primaryKey"`
	Code        string         `gorm:"size:64;uniqueIndex"`
	Name        string         `gorm:"size:128"`
	Description string         `gorm:"type:text"`
	BasePoints  int64          `gorm:"column:base_points"`
	MaxPoints   *int64         `gorm:"column:max_points"`
	Stackable   bool           `gorm:"column:stackable"`
	Metadata    datatypes.JSON `gorm:"column:metadata"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

func (Task) TableName() string {
	return "tasks"
}

type UserTask struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"column:user_id;index"`
	TaskCode  string    `gorm:"column:task_code;size:64"`
	Amount    int64     `gorm:"column:amount"`
	Points    int64     `gorm:"column:points"`
	UniqueKey string    `gorm:"column:unique_key;size:128"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserTask) TableName() string {
	return "user_tasks"
}

type PointsLedger struct {
	ID        uint64         `gorm:"primaryKey"`
	UserID    uint64         `gorm:"column:user_id;index"`
	Delta     int64          `gorm:"column:delta"`
	Reason    string         `gorm:"column:reason;size:64"`
	RefID     *uint64        `gorm:"column:ref_id"`
	Meta      datatypes.JSON `gorm:"column:meta"`
	CreatedAt time.Time      `gorm:"column:created_at"`
}

func (PointsLedger) TableName() string {
	return "points_ledger"
}

type AirdropRound struct {
	ID            uint64     `gorm:"primaryKey"`
	Name          string     `gorm:"size:64"`
	MerkleRoot    string     `gorm:"column:merkle_root;size:128"`
	TokenAddress  string     `gorm:"column:token_address;size:64"`
	ClaimDeadline time.Time  `gorm:"column:claim_deadline"`
	Status        string     `gorm:"column:status;size:32"`
	SnapshotAt    *time.Time `gorm:"column:snapshot_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

func (AirdropRound) TableName() string {
	return "airdrop_rounds"
}

type RoundPoint struct {
	ID            uint64    `gorm:"primaryKey"`
	RoundID       uint64    `gorm:"column:round_id;index"`
	UserID        uint64    `gorm:"column:user_id;index"`
	Points        int64     `gorm:"column:points"`
	ClaimedPoints int64     `gorm:"column:claimed_points"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (RoundPoint) TableName() string {
	return "round_points"
}

type Claim struct {
	ID        uint64     `gorm:"primaryKey"`
	RoundID   uint64     `gorm:"column:round_id;index"`
	UserID    uint64     `gorm:"column:user_id;index"`
	Wallet    string     `gorm:"column:wallet;size:64"`
	Amount    int64      `gorm:"column:amount"`
	TxHash    string     `gorm:"column:tx_hash;size:128"`
	Status    string     `gorm:"column:status;size:32"`
	ClaimedAt *time.Time `gorm:"column:claimed_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

func (Claim) TableName() string {
	return "claims"
}
