package dto

const DEFAULT_SIZE = 10

type PaginationRequest struct {
	Offset int 
	Limit int 
	Page int 
}

type PaginationResponse struct {
	TotalPages int `json:"total_pages"`
	Page int `json:"page"`
}

func NewPaginationRequest(page int) *PaginationRequest {
	p := &PaginationRequest{}

	page -= 1

	p.Offset = page * DEFAULT_SIZE
	p.Limit = DEFAULT_SIZE
	p.Page = page + 1
	
	return p 
}