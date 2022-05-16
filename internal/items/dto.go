package items

type CreateItemDTO struct {
	Name         string `form:"name" json:"name" binding:"required"`
	ProductName  string `form:"productName" json:"productName" binding:"required"`
	SerialNumber string `form:"serialNumber" json:"serialNumber" binding:"required"`
}

type EditItemDTO struct {
	ID           uint   `form:"id" json:"id" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	ProductName  string `form:"productName" json:"productName" binding:"required"`
	SerialNumber string `form:"serialNumber" json:"serialNumber" binding:"required"`
	OwnerID      *int   `form:"ownerid" json:"ownerid"`
}
