package service

import "github.com/xiaoxue1272/club-5fw-backend/model"

func Sign(userSign *model.UserSign) (*model.UserJwt, error) {
	// todo 真正的用户登录
	return &model.UserJwt{Account: "test"}, nil
}
