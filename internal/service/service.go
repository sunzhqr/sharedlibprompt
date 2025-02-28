package service

import test "github.com/sunzhqr/sharedlibprompt/pkg/api/test/api"

type Service struct {
	test.OrderServiceServer
}

func New() *Service {
	return &Service{}
}