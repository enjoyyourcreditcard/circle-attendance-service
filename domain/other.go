package domain

type PaginatedResponse struct {
	TotalItems 	int64 	`json:"total_items"`
	TotalPages 	int		`json:"total_pages"`
	CurrentPage int		`json:"current_page"`
	NextPage 	int		`json:"next_page"`
	PrevPage 	int		`json:"prev_page"`
}