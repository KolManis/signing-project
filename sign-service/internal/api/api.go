package api

import (
	signV1 "github.com/KolManis/signing-project/shared/pkg/proto/sign/v1"
	"github.com/KolManis/signing-project/sign-service/internal/service"
)

type api struct {
	signV1.UnimplementedSignServiceServer

	signService service.SignService
}

func NewAPI(signService service.SignService) *api {
	return &api{
		signService: signService,
	}
}
