package dto

type BaseListResponse[T any] struct {
	// Day 6 add omitempty to avoid null in JSON response ตัว get rating ใช้ด้วยแต่ ไม่เอา total
	Total   int64 `json:"total,omitempty"`
	Results T     `json:"results"`
}
