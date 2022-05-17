package user

type EditUserDTO struct {
	Password *string `form:"password" json:"password"`
	Familia  *string `form:"familia" json:"familia"`
	Name     *string `form:"name" json:"name"`
}
