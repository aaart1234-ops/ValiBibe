package dto

type TagCreateInput struct {
    Name string `json:"name" binding:"required"`
}

type TagUpdateInput struct {
    Name string `json:"name" binding:"required"`
}

type NoteTagInput struct {
	NoteID string `json:"note_id" binding:"required"`
	TagID  string `json:"tag_id" binding:"required"`
}
