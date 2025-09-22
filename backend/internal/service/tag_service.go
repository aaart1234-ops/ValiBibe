package service

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "valibibe/internal/controller/dto"
    "valibibe/internal/models"
    "valibibe/internal/repository/interfaces"
)

type TagService struct {
    repo interfaces.TagRepository
}

func NewTagService(repo interfaces.TagRepository) *TagService {
    return &TagService{repo: repo}
}

func (s *TagService) CreateTag(ctx context.Context, userID string, input dto.TagCreateInput) (*models.Tag, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, errors.New("invalid userID")
    }

    exists, err := s.repo.ExistsByName(ctx, uid, input.Name)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("tag with this name already exists")
    }

    tag := &models.Tag{
        UserID: uid,
        Name:   input.Name,
    }

    if err := s.repo.Create(ctx, tag); err != nil {
        return nil, err
    }
    return tag, nil
}

func (s *TagService) GetTags(ctx context.Context, userID string) ([]models.Tag, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, errors.New("invalid userID")
    }
    return s.repo.ListByUser(ctx, uid)
}

func (s *TagService) GetTag(ctx context.Context, userID, tagID string) (*models.Tag, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, errors.New("invalid userID")
    }
    tid, err := uuid.Parse(tagID)
    if err != nil {
        return nil, errors.New("invalid tagID")
    }
    return s.repo.GetByID(ctx, uid, tid)
}

func (s *TagService) UpdateTag(ctx context.Context, userID, tagID string, input dto.TagUpdateInput) (*models.Tag, error) {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return nil, errors.New("invalid userID")
    }
    tid, err := uuid.Parse(tagID)
    if err != nil {
        return nil, errors.New("invalid tagID")
    }

    tag, err := s.repo.GetByID(ctx, uid, tid)
    if err != nil {
        return nil, err
    }
    if tag == nil {
        return nil, errors.New("tag not found")
    }

    // проверка дубля
    exists, err := s.repo.ExistsByName(ctx, uid, input.Name)
    if err != nil {
        return nil, err
    }
    if exists && tag.Name != input.Name { // чтобы не ругалось, если имя не изменилось
        return nil, errors.New("tag with this name already exists")
    }

    tag.Name = input.Name
    if err := s.repo.Update(ctx, tag); err != nil {
        return nil, err
    }

    return tag, nil
}

func (s *TagService) DeleteTag(ctx context.Context, userID, tagID string) error {
    uid, err := uuid.Parse(userID)
    if err != nil {
        return errors.New("invalid userID")
    }
    tid, err := uuid.Parse(tagID)
    if err != nil {
        return errors.New("invalid tagID")
    }
    return s.repo.Delete(ctx, uid, tid)
}
