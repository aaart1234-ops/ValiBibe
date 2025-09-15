package dto

type FolderNode struct {
    ID       string       `json:"id"`
    Name     string       `json:"name"`
    ParentID *string      `json:"parent_id,omitempty"`
    Children []FolderNode `json:"children"`
}
