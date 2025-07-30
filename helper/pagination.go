package helper

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaginationData struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	NextPage     *int  `json:"next_page"`
	PrevPage     *int  `json:"prev_page"`
}

func Paginate(db *gorm.DB, c *fiber.Ctx, result interface{}) (*PaginationData, error) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if page <= 0 {
		page = 1
	}

	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	var totalRecords int64

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	if err := db.Limit(perPage).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}
	totalPages := int(math.Ceil(float64(totalRecords) / float64(perPage)))

	pagination := &PaginationData{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PerPage:      perPage,
		NextPage:     nil,
		PrevPage:     nil,
	}

	if page < totalPages {
		next := page + 1
		pagination.NextPage = &next
	}

	if page > 1 {
		prev := page - 1
		pagination.PrevPage = &prev
	}

	return pagination, nil
}
