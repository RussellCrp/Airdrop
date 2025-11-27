// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"

	"airdrop/internal/entity"
	"airdrop/internal/svc"
	"airdrop/internal/tasks"
	"airdrop/internal/types"
	"airdrop/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AwardTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAwardTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AwardTaskLogic {
	return &AwardTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AwardTaskLogic) AwardTask(req *types.AwardTaskRequest) (*types.AwardTaskResponse, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	wallet, err := util.NormalizeWallet(req.Wallet)
	if err != nil {
		return nil, err
	}
	taskCode := tasks.Normalize(req.Task)
	uniqueKey, err := tasks.UniqueKey(taskCode, req.Extra)
	if err != nil {
		return nil, err
	}
	points, err := tasks.PointsFor(taskCode, req.Amount)
	if err != nil {
		return nil, err
	}
	if points <= 0 {
		return nil, errors.New("no points to award")
	}

	var user entity.User
	err = l.svcCtx.RunTx(l.ctx, func(tx *gorm.DB) error {
		if err := tx.Where("wallet = ?", wallet).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = entity.User{Wallet: wallet}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		record := &entity.UserTask{
			UserID:    user.ID,
			TaskCode:  taskCode,
			Amount:    req.Amount,
			Points:    points,
			UniqueKey: uniqueKey,
		}
		if err := tx.Create(record).Error; err != nil {
			if strings.Contains(err.Error(), "uk_user_task_once") {
				return errors.New("task already claimed")
			}
			return err
		}
		return l.svcCtx.AwardPointsInTx(l.ctx, &user, points, "task-"+taskCode, map[string]interface{}{"task": taskCode}, tx)
	})
	if err != nil {
		return nil, err
	}

	return &types.AwardTaskResponse{
		Points: user.PointsBalance,
	}, nil
}
