package gapi

import (
	"context"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hibiken/asynq"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/util"
	"github.com/xiaowuzai/simplebank/validator"
	"github.com/xiaowuzai/simplebank/worker"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashPassword,
			FullName:       req.FullName,
			Email:          req.Email,
		},
		AfterCreateUser: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),                // 重试 10 次
				asynq.ProcessIn(10 * time.Second), // 延迟 10 秒处理
				asynq.Queue(worker.QueueCritical), // 队列名称
			}
			return s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := s.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolationCode {
			return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.CreateUserResponse{
		User: convertUserToPb(txResult.User),
	}
	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolations("username", err))
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolations("password", err))
	}

	if err := validator.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, fieldViolations("full_name", err))
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolations("email", err))
	}

	return
}
