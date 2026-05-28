package dto

import "github.com/BangNopall/hology8-be/domain/entity"

type ProvinceResponse struct {
	ID   int    `json:"id" `
	Name string `json:"province_name" `
}

func ProvinceEntityToDto(province *entity.Province) ProvinceResponse {
	return ProvinceResponse{
		ID:   province.ID,
		Name: province.Name,
	}
}
