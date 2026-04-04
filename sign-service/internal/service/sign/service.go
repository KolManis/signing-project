package sign

import (
	def "github.com/KolManis/signing-project/sign-service/internal/service"
)

var _ def.SignService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
