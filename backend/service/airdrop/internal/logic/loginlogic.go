// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"time"

	"airdrop/internal/entity"
	"airdrop/internal/security"
	"airdrop/internal/svc"
	"airdrop/internal/types"
	"airdrop/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const loginMessagePrefix = "airdrop-login"
const loginSkew = 5 * time.Minute

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	wallet, err := util.NormalizeWallet(req.Wallet)
	if err != nil {
		return nil, err
	}
	if !l.validateTimestamp(req.Timestamp) {
		return nil, errors.New("timestamp skew too large")
	}
	// message := fmt.Sprintf("%s:%s:%d", loginMessagePrefix, wallet, req.Timestamp)
	// if err := util.VerifyPersonalSignature(wallet, req.Signature, message); err != nil {
	// 	return nil, err
	// }

	now := time.Now().UTC()
	var user entity.User
	if err := l.svcCtx.RunTx(l.ctx, func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("wallet = ?", wallet).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				user = entity.User{
					Wallet:      wallet,
					LoginStreak: 0,
					LoginDays:   0,
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		newDay := user.LastLoginAt == nil || !util.SameDay(*user.LastLoginAt, now)
		if newDay {
			if user.LastLoginAt != nil && util.IsYesterday(*user.LastLoginAt, now) {
				if user.LoginStreak < 5 {
					user.LoginStreak++
				}
			} else {
				user.LoginStreak = 1
			}
			if user.LoginStreak > 5 {
				user.LoginStreak = 5
			}
			user.LoginDays++
			user.LastLoginAt = &now
			reward := int64(user.LoginStreak) * 100
			if reward > 0 {
				if err := l.svcCtx.AwardPointsInTx(l.ctx, &user, reward, "login-streak", map[string]interface{}{
					"streak": user.LoginStreak,
				}, tx); err != nil {
					return err
				}
			}
			return tx.Save(&user).Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	role := security.RoleUser
	if l.svcCtx.IsAdminWallet(wallet) {
		role = security.RoleAdmin
	}
	token, expiresAt, err := l.svcCtx.JWTManager.Generate(user.ID, wallet, role)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		BaseResp: types.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Data: types.LoginData{
			AccessToken: token,
			ExpiresAt:   expiresAt,
			LoginDays:   int64(user.LoginDays),
			Points:      user.PointsBalance,
		},
	}, nil
}

func (l *LoginLogic) validateTimestamp(ts int64) bool {
	if ts == 0 {
		return false
	}
	t := time.Unix(ts, 0)
	diff := time.Since(t)
	if diff < 0 {
		diff = -diff
	}
	return diff <= loginSkew
}
