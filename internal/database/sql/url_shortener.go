package sql

import (
	"context"
	"errors"

	"github.com/Parzival-05/url-shortener/internal/domain"

	"gorm.io/gorm"
)

type UrlRepositoryPG struct {
	db service
}

func NewUrlRepositoryPG(db service) *UrlRepositoryPG {
	return &UrlRepositoryPG{db: db}
}

func (u *UrlRepositoryPG) GetID(ctx context.Context, fullUrl string) (id int64, err error) {
	url, err := gorm.G[Url](u.db.db).Where("full_url = ?", fullUrl).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, domain.ErrUrlNotFound
		}
		return 0, err
	}
	return url.Id, nil
}

func (u *UrlRepositoryPG) GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error) {
	url, err := gorm.G[Url](u.db.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return "", err
	}
	return url.FullUrl, nil
}

func (u *UrlRepositoryPG) SaveUrl(ctx context.Context, fullUrl string) (err error) {
	url := Url{FullUrl: fullUrl}
	err = gorm.G[Url](u.db.db).Create(ctx, &url)
	return err
}
