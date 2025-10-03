package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "valibibe/internal/models"
	"valibibe/internal/controller/dto"
    "valibibe/internal/repository/interfaces"
)

type AssignFolderService struct {
    
}

func (s *NoteService) AssignFolder(ctx context.Context, userID, noteID, folderID string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return errors.New("invalid userID")
    }
    nid, err := uuid.Parse(noteID)
    if err != nil {
        return errors.New("invalid noteID")
    }
    fid, err := uuid.Parse(folderID)
    if err != nil {
        return errors.New("invalid folderID")
    }

    return s.noteRepo.UpdateFolder(ctx, uid, nid, &fid)
}

func (s *NoteService) RemoveFolder(ctx context.Context, userID, noteID, folderID string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return errors.New("invalid userID")
    }
    nid, err := uuid.Parse(noteID)
    if err != nil {
        return errors.New("invalid noteID")
    }
    fid, err := uuid.Parse(folderID)
    if err != nil {
        return errors.New("invalid folderID")
    }

    // Доп проверка: принадлежала ли заметка этой папке
    note, err := s.noteRepo.GetNoteByID(ctx, uid, nid)
    if err != nil {
        return err
    }
    if note == nil || note.FolderID == nil || *note.FolderID != fid {
        return errors.New("note is not in this folder")
    }

    return s.noteRepo.UpdateFolder(ctx, uid, nid, nil)
}

func (s *NoteService) BatchAssignFolder(ctx context.Context, userID string, noteIDs []string, folderID *string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return errors.New("invalid userID")
    }

    var fUUID *uuid.UUID
    if folderID != nil {
        fid, err := uuid.Parse(*folderID)
        if err != nil {
            return errors.New("invalid folderID")
        }
        fUUID = &fid
    }

    var nUUIDs []uuid.UUID
    for _, id := range noteIDs {
        nid, err := uuid.Parse(id)
        if err != nil {
            return errors.New("invalid noteID in batch")
        }
        nUUIDs = append(nUUIDs, nid)
    }

    return s.noteRepo.BatchUpdateFolder(ctx, uid, nUUIDs, fUUID)
}