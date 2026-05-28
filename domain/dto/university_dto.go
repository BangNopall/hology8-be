package dto

import "github.com/hology8/hology-be/domain/entity"

type UniversityResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"university_name"`
}

type UniversityParam struct {
	Name string 
}

func UniversityEntityToDto(uni *entity.University) UniversityResponse {
	return UniversityResponse{
		ID: uni.ID,
		Name: uni.Name,
	}
}