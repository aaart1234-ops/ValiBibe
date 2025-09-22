package dto

type TagCreateInput struct {
    Name string `json:"name" binding:"required"`
}

type TagUpdateInput struct {
    Name string `json:"name" binding:"required"`
}
