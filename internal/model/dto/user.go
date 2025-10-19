package dto

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
}

type UpdateNicknameRequest struct {
	NickName string `json:"nickname" binding:"required,min=2"`
}
