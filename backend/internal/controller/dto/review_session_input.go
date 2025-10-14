package dto

// ReviewSessionInput представляет входные данные для создания сессии повторения
type ReviewSessionInput struct {
	FolderID *string  `json:"folder_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TagIDs   []string `json:"tag_ids" example:"550e8400-e29b-41d4-a716-446655440001,550e8400-e29b-41d4-a716-446655440002"`
	Limit    int      `json:"limit" example:"10" minimum:"1" maximum:"100"`
}

// ReviewSessionResponse представляет ответ с заметками для повторения
type ReviewSessionResponse struct {
	Notes []ReviewSessionNote `json:"notes"`
	Total int                 `json:"total"`
}

// ReviewSessionNote представляет заметку в сессии повторения
type ReviewSessionNote struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	MemoryLevel   int    `json:"memory_level"`
	NextReviewAt  string `json:"next_review_at,omitempty"`
	FolderID      string `json:"folder_id,omitempty"`
	FolderName    string `json:"folder_name,omitempty"`
	Tags          []Tag  `json:"tags"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Tag представляет тег в контексте сессии повторения
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

