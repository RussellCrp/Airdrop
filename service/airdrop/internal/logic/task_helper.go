package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"airdrop/internal/model"
	"airdrop/internal/svc"

	"gorm.io/gorm"
)

const (
	taskCodePromote      = "PROMOTE"
	taskCodeQuantInvest  = "QUANT_INVEST"
	taskCodeTrade10k     = "TRADE_VOLUME_10K"
	taskCodeReferral10   = "REFERRAL_10"
	taskCodeLogin7Days   = "LOGIN_7_DAYS"
	userTaskStatusDone   = 1
)

// 完成一个任务并累计积分
func completeTaskAndUpdateScore(s *svc.ServiceContext, userID uint64, taskCode string, extra map[string]any) error {
	var task model.Task
	if err := s.DB.Where("code = ? AND enabled = 1", taskCode).First(&task).Error; err != nil {
		return err
	}

	// upsert user_tasks
	var ut model.UserTask
	err := s.DB.Where("user_id = ? AND task_id = ?", userID, task.ID).First(&ut).Error
	now := time.Now()
	extraJSON := ""
	if extra != nil {
		bz, _ := json.Marshal(extra)
		extraJSON = string(bz)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ut = model.UserTask{
			UserID:      userID,
			TaskID:      task.ID,
			Status:      userTaskStatusDone,
			ExtraData:   extraJSON,
			CompletedAt: &now,
		}
		if err = s.DB.Create(&ut).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if ut.Status != userTaskStatusDone {
			ut.Status = userTaskStatusDone
			ut.ExtraData = extraJSON
			ut.CompletedAt = &now
			if err = s.DB.Save(&ut).Error; err != nil {
				return err
			}
		}
	}

	// 更新 user_scores
	var us model.UserScore
	if err := s.DB.Where("user_id = ?", userID).First(&us).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			us = model.UserScore{
				UserID:     userID,
				TotalScore: task.ScoreWeight,
			}
			return s.DB.Create(&us).Error
		}
		return err
	}

	us.TotalScore += task.ScoreWeight
	return s.DB.Save(&us).Error
}

// 查询用户积分和所在档位
func getUserScoreAndTier(s *svc.ServiceContext, userID uint64) (score int64, tierName string, err error) {
	var us model.UserScore
	if err = s.DB.Where("user_id = ?", userID).First(&us).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", nil
		}
		return 0, "", err
	}

	score = int64(us.TotalScore)
	var tier model.ScoreTier
	if err = s.DB.
		Where("min_score <= ? AND max_score >= ?", us.TotalScore, us.TotalScore).
		First(&tier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return score, "", nil
		}
		return 0, "", err
	}
	return score, fmt.Sprintf("%d-%d", tier.MinScore, tier.MaxScore), nil
}


