package dto

type AssignFolderInput struct {
    FolderID string `json:"folder_id" binding:"required,uuid" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`
}

type BatchAssignFolderInput struct {
    NoteIDs []string `json:"note_ids" binding:"required,dive,uuid" example:"[\"a12b34cd-5678-90ef-1234-567890abcdef\", \"b34c56de-7890-12gh-3456-789012ijklmn\"]"`
    FolderID *string  `json:"folder_id" binding:"omitempty,uuid" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"` // null → убрать связь
}