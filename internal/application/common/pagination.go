package common

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type Sorting struct {
	SortBy  string `form:"sort_by"`
	OrderBy string `form:"order_by"`
}

type QueryOptions struct {
	PaginationRequest
	Sorting
	Search string `form:"search"`
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

func (s Sorting) Normalize(allowed map[string]bool, defaultSort string) Sorting {
	if !allowed[s.SortBy] {
		s.SortBy = defaultSort
	}

	switch s.OrderBy {
	case "asc", "ASC":
		s.OrderBy = "ASC"
	case "desc", "DESC":
		s.OrderBy = "DESC"
	default:
		s.OrderBy = "DESC"
	}

	return s
}

func (q QueryOptions) Normalize(allowedSorts map[string]bool, defaultSort string) QueryOptions {
	q.PaginationRequest = q.PaginationRequest.Normalize()
	q.Sorting = q.Sorting.Normalize(allowedSorts, defaultSort)

	return q
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
