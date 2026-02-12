package input

type UserLogin struct {
	Name string `json:"name" binding:"required"`
	Pass string `json:"pass" binding:"required"`
}
