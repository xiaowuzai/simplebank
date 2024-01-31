package gapi

import (
	"context"
	"log"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lib/pq"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
	"github.com/xiaowuzai/simplebank/pb"
	"github.com/xiaowuzai/simplebank/util"
	"github.com/xiaowuzai/simplebank/validator"
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
