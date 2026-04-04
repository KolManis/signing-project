package api

import (
	"context"
	"errors"

	signV1 "github.com/KolManis/signing-project/shared/pkg/proto/sign/v1"
	"github.com/KolManis/signing-project/sign-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) SignDocument(ctx context.Context, req *signV1.SignRequest) (*signV1.SignResponse, error) {

	signed, err := a.signService.SignDocument(ctx, req.UserId, req.Document)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidDocument):
			return nil, status.Error(codes.InvalidArgument, err.Error())

		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &signV1.SignResponse{
		Signature:  signed.Signature,
		DocumentId: signed.DocumentID,
	}, nil
}
