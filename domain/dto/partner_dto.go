package dto

import "github.com/hology8/hology-be/domain/entity"

type PartnerTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"partner_type_name"`
}

type PartnerCreate struct {
	Name          string `json:"name"`
	PartnerTypeID int    `json:"partner_type_id"`
}

type PartnerUpdate struct {
	Name          string `json:"name"`
	PartnerTypeID int    `json:"partner_type_id"`
}

type PartnerParams struct {
	PartnerTypeID int `json:"partner_type_id"`
}

type PartnerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type PartnersResponse struct {
	Partners []PartnerResponse `json:"partners"`
}

func PartnerEntityToResponse(partner entity.Partner) PartnerResponse {
	return PartnerResponse{
		ID:   partner.ID,
		Name: partner.Name,
		Link: partner.ImageLink,
	}
}

func PartnerSliceToPartnersReponse(partners []entity.Partner) PartnersResponse {
	res := make([]PartnerResponse, len(partners))

	for i, partner := range partners {
		res[i] = PartnerEntityToResponse(partner)
	}

	return PartnersResponse{
		Partners: res,
	}
}

func PartnerTypeEntityToDto(partnerType *entity.PartnerType) PartnerTypeResponse {
	return PartnerTypeResponse{
		ID:   partnerType.ID,
		Name: partnerType.Name,
	}
}
