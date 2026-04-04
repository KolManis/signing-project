package service

import (
	"context"

	"github.com/KolManis/signing-project/sign-service/internal/model"
)

type SignService interface {
	SignDocument(ctx context.Context, userId string, document []byte) (*model.SignedDocument, error)
}
