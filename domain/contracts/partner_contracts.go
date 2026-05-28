package contracts

import (
	"context"
	"mime/multipart"

	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type PartnerRepository interface {
	FetchAll(ctx context.Context, params *dto.PartnerParams) ([]entity.Partner, error)
	FetchAllTypes(ctx context.Context) ([]entity.PartnerType, error)
	FetchOneByID(ctx context.Context, id int) (entity.Partner, error)
	InsertPartner(ctx context.Context, partner *entity.Partner) error
	UpdatePartner(ctx context.Context, id int, partner *entity.Partner) error
	DeletePartner(ctx context.Context, id int) error
}

type PartnerService interface {
	FetchAll(ctx context.Context, params *dto.PartnerParams) (dto.PartnersResponse, error)
	FetchAllTypes(ctx context.Context) ([]dto.PartnerTypeResponse, error)
	FetchOneByID(ctx context.Context, id int) (dto.PartnerResponse, error)
	CreatePartner(ctx context.Context, partner *dto.PartnerCreate, imageFile *multipart.FileHeader) error
	UpdatePartner(ctx context.Context, id int, partner *dto.PartnerUpdate, image *multipart.FileHeader) error
	DeletePartner(ctx context.Context, id int) error
}
