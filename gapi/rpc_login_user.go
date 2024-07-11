package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "error: User NOT found, %s", err)
		}

		return nil, status.Errorf(codes.InvalidArgument, "error: Could not search user, %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error: User with provided credentials does not exist, %s", err)
	}

	token, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: Failed to create access token, %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: Failed to create refresh token, %s", err)
	}

	mtdata := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdata.userAgent,
		ClientIp:     mtdata.clientIp,
		ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: Failed to create session, %s", err)
	}

	response := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           token,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt.Time),
		User:                  convertUser(user),
	}

	return response, nil
}
