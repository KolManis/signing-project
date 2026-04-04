package sign

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/KolManis/signing-project/sign-service/internal/model"
)

func (s *service) SignDocument(ctx context.Context, userId string, document []byte) (*model.SignedDocument, error) {
	if len(document) == 0 {
		return nil, model.ErrInvalidDocument
	}

	hash := sha256.Sum256(document)
	signature := fmt.Sprintf("%x", hash[:])
	documentID := fmt.Sprintf("%x-%s", hash[:4], userId)

	return &model.SignedDocument{
		DocumentID: documentID,
		Signature:  signature,
	}, nil
}
