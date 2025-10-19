package dto

type UserResponse struct {
	ID              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	NickName        string `json:"nickName"`
	ImageProfileUrl string `json:"imageProfileUrl"`
}

type UpdateProfileRequest struct {
	NickName        string `json:"nickname,omitempty"`
	ImageProfileUrl string `json:"imageProfileUrl,omitempty"`
}
