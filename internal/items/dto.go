package items

type CreateItemDTO struct {
	Name         string `form:"name" json:"name" binding:"required"`
	ProductName  string `form:"productName" json:"productName" binding:"required"`
	SerialNumber string `form:"serialNumber" json:"serialNumber" binding:"required"`
}

type EditItemDTO struct {
	Name         *string `form:"name" json:"name"`
	ProductName  *string `form:"productName" json:"productName"`
	SerialNumber *string `form:"serialNumber" json:"serialNumber"`
	OwnerID      *int    `form:"ownerid" json:"ownerid"`
}
