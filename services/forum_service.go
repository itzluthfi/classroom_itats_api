package services

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	"context"
	"encoding/json"
)

type forumService struct {
	forumRepository repositories.ForumRepository
}

type ForumService interface {
	StudentForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error)
	LecturerForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error)
	Forums(ctx context.Context, masterActivityID string) ([]entities.AnnouncementJSON, error)
	StoreAnnouncement(ctx context.Context, announcement entities.AnnouncementStore) error
	StoreComment(ctx context.Context, comment entities.CommentStore) error
	UpdateAnnouncement(ctx context.Context, announcement entities.AnnouncementUpdate) error
	DeleteAnnouncement(ctx context.Context, announcementID int) error
	UpdateComment(ctx context.Context, comment entities.CommentUpdate) error
	DeleteComment(ctx context.Context, commentID int) error
}

func NewForumService(forumRepository repositories.ForumRepository) *forumService {
	return &forumService{
		forumRepository: forumRepository,
	}
}

func (f *forumService) StudentForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error) {
	return f.forumRepository.StudentForum(ctx, SubjectID, SubjectClass)
}

func (f *forumService) LecturerForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error) {
	return f.forumRepository.LecturerForum(ctx, SubjectID, SubjectClass)
}

func (f *forumService) Forums(ctx context.Context, masterActivityID string) ([]entities.AnnouncementJSON, error) {
	rawAnnouncement, err := f.forumRepository.Forums(ctx, masterActivityID)

	if err != nil {
		return []entities.AnnouncementJSON{}, err
	}

	announcements := []entities.AnnouncementJSON{}

	for _, v := range rawAnnouncement {
		materials := []entities.AnnouncementMaterial{}
		err = json.Unmarshal([]byte(v.Materials), &materials)
		if err != nil {
			return nil, err
		}
		comments := []entities.Comment{}
		err := json.Unmarshal([]byte(v.Comments), &comments)
		if err != nil {
			return nil, err
		}
		if materials[0].MaterialID != "" {
			announcements = append(announcements, entities.AnnouncementJSON{
				AnnouncementID: v.AnnouncementID,
				PostContent:    v.PostContent,
				CreatedAt:      v.CreatedAt,
				Author:         v.Author,
				AuthorID:       v.AuthorID,
				Photo:          v.Photo,
				Materials:      materials,
				Comments:       comments,
			})
		} else {
			announcements = append(announcements, entities.AnnouncementJSON{
				AnnouncementID: v.AnnouncementID,
				PostContent:    v.PostContent,
				CreatedAt:      v.CreatedAt,
				Author:         v.Author,
				AuthorID:       v.AuthorID,
				Photo:          v.Photo,
				Comments:       comments,
				Materials:      []entities.AnnouncementMaterial{},
			})
		}

	}

	return announcements, nil
}

func (f *forumService) StoreAnnouncement(ctx context.Context, announcement entities.AnnouncementStore) error {
	return f.forumRepository.StoreAnnouncement(ctx, announcement)
}

func (f *forumService) UpdateAnnouncement(ctx context.Context, announcement entities.AnnouncementUpdate) error {
	return f.forumRepository.UpdateAnnouncement(ctx, announcement)
}

func (f *forumService) DeleteAnnouncement(ctx context.Context, announcementID int) error {
	return f.forumRepository.DeleteAnnouncement(ctx, announcementID)
}

func (f *forumService) StoreComment(ctx context.Context, comment entities.CommentStore) error {
	return f.forumRepository.StoreComment(ctx, comment)
}

func (f *forumService) UpdateComment(ctx context.Context, comment entities.CommentUpdate) error {
	return f.forumRepository.UpdateComment(ctx, comment)
}

func (f *forumService) DeleteComment(ctx context.Context, commentID int) error {
	return f.forumRepository.DeleteComment(ctx, commentID)
}
