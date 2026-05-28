package service

import (
	"context"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/pkg/aws"
)

type partnerService struct {
	partnerRepository contracts.PartnerRepository
	timeout           time.Duration
	aws               aws.CloudStorage
}

func NewPartnerService(partnerRepository contracts.PartnerRepository, timeout time.Duration, aws aws.CloudStorage) contracts.PartnerService {
	return &partnerService{partnerRepository, timeout, aws}
}

func (p *partnerService) FetchAll(ctx context.Context, params *dto.PartnerParams) (dto.PartnersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	partners, err := p.partnerRepository.FetchAll(ctx, params)

	res := dto.PartnerSliceToPartnersReponse(partners)

	select {
	case <-ctx.Done():
		return dto.PartnersResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}

func (p *partnerService) FetchAllTypes(ctx context.Context) ([]dto.PartnerTypeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	partnerTypes, err := p.partnerRepository.FetchAllTypes(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]dto.PartnerTypeResponse, 0)

	for _, partnerType := range partnerTypes {
		res = append(res, dto.PartnerTypeEntityToDto(&partnerType))
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return res, err
	}
}

func (p *partnerService) FetchOneByID(ctx context.Context, id int) (dto.PartnerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	partner, err := p.partnerRepository.FetchOneByID(ctx, id)

	res := dto.PartnerEntityToResponse(partner)

	select {
	case <-ctx.Done():
		return dto.PartnerResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}

func (p *partnerService) CreatePartner(ctx context.Context, partner *dto.PartnerCreate, imageFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if imageFile.Size > 1024*1024 {
		return domain.ErrFileTooBig
	}

	path := "partners/" + url.QueryEscape(partner.Name) + "-" + time.Now().Format("02-Jan-2006 15:04:05")

	link, err := p.aws.Upload(path, imageFile)
	if err != nil {
		return err
	}

	newPartner := entity.NewPartner(partner.Name, link, partner.PartnerTypeID)

	err = p.partnerRepository.InsertPartner(ctx, newPartner)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (p *partnerService) UpdatePartner(ctx context.Context, id int, partner *dto.PartnerUpdate, image *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var (
		link string
		err  error
	)

	if image.Size > 1024*1024 {
		return domain.ErrFileTooBig
	}

	oldPartner, err := p.partnerRepository.FetchOneByID(ctx, id)
	if err != nil {
		return err
	}

	path := "partners/" + url.QueryEscape(partner.Name) + "-" + time.Now().Format("02-Jan-2006 15:04:05")

	if oldPartner.ImageLink != "" {
		link, err = p.aws.Update(path, image, oldPartner.ImageLink)
	} else {
		link, err = p.aws.Upload(path, image)
	}

	if err != nil {
		return err
	}

	partnerToUpdate := entity.NewPartner(partner.Name, link, partner.PartnerTypeID)

	err = p.partnerRepository.UpdatePartner(ctx, id, partnerToUpdate)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (p *partnerService) DeletePartner(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	partner, err := p.partnerRepository.FetchOneByID(ctx, id)
	if err != nil {
		return err
	}

	err = p.aws.Delete(partner.ImageLink)
	if err != nil {
		return err
	}

	err = p.partnerRepository.DeletePartner(ctx, id)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}
