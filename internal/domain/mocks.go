package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UrlRepositoryMock struct {
	mock.Mock
}

func (u *UrlRepositoryMock) GetID(ctx context.Context, fullUrl string) (id int64, err error) {
	args := u.Called(ctx, fullUrl)
	return int64(args.Int(0)), args.Error(1)
}

func (u *UrlRepositoryMock) GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error) {
	args := u.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (u *UrlRepositoryMock) SaveUrl(ctx context.Context, fullUrl string) (err error) {
	args := u.Called(ctx, fullUrl)
	return args.Error(0)
}
