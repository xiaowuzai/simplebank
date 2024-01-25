package gapi

import (
	"context"
	"database/sql"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/lib/pq"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/util"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.CreateUserResponse{
		User: convertUserToPb(user),
	}
	return res, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	token, tokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	md := getMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    md.UserAgent,
		ClientIp:     md.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		Token:                 token,
		RefreshToken:          refreshToken,
		TokenExpiredAt:        timestamppb.New(tokenPayload.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUserToPb(user),
	}
	return res, nil
}
