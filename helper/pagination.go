package helper

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	NextPage     *int  `json:"next_page"`
	PrevPage     *int  `json:"prev_page"`
}

func Paginate() {

}
