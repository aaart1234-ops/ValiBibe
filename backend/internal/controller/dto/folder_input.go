package dto

type FolderCreateInput struct {
    Name     string  `json:"name" binding:"required,min=1,max=200"`
    ParentID *string `json:"parent_id,omitempty"`
}

type FolderUpdateInput struct {
    Name     *string `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
    ParentID *string `json:"parent_id,omitempty"`
}
