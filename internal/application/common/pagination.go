package common

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (p PaginationRequest) Normalize() PaginationRequest {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 10
	}
	return p
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}
