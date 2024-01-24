package gapi

// 转换数据结构

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/pb"
)

func convertUserToPb(du db.User) *pb.User {
	return &pb.User{
		Username:          du.Username,
		FullName:          du.FullName,
		Email:             du.Email,
		PasswordChangedAt: timestamppb.New(du.PasswordChangedAt),
		CreatedAt:         timestamppb.New(du.CreatedAt),
	}
}
