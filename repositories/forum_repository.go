package repositories

import (
	"classroom_itats_api/entities"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type forumRepository struct {
	db *gorm.DB
}

type ForumRepository interface {
	StudentForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error)
	LecturerForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error)
	Forums(ctx context.Context, masterActivityID string) ([]entities.Announcement, error)
	StoreAnnouncement(ctx context.Context, announcement entities.AnnouncementStore) error
	StoreComment(ctx context.Context, comment entities.CommentStore) error
	UpdateAnnouncement(ctx context.Context, announcement entities.AnnouncementUpdate) error
	DeleteAnnouncement(ctx context.Context, announcementID int) error
	UpdateComment(ctx context.Context, comment entities.CommentUpdate) error
	DeleteComment(ctx context.Context, commentID int) error
}

func NewForumRepository(db *gorm.DB) *forumRepository {
	return &forumRepository{
		db: db,
	}
}

func (f *forumRepository) StudentForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error) {
	announcement := []entities.Announcement{}

	err := f.db.WithContext(ctx).Raw(`select 
				post_klstw.id_post_klstw, post_klstw.content_post,
				post_klstw.created_at, mhs.mhsnama as nama mhs.foto
				from post_klstw
				join vw_kelas_tawar on cast(vw_kelas_tawar.id_master_kegiatan as character varying) = post_klstw.master_kegiatan_id
				join mk on mk.id_mata_kuliah = vw_kelas_tawar.id_mata_kuliah
				join users on users.name = post_klstw.author_id
				join mhs on users.name = mhs.mhsid
				where mk.mkid = ? and kelas = ?;`, SubjectID, SubjectClass).
		Find(&announcement).Error

	return announcement, err
}

func (f *forumRepository) LecturerForum(ctx context.Context, SubjectID string, SubjectClass string) ([]entities.Announcement, error) {
	announcement := []entities.Announcement{}

	err := f.db.WithContext(ctx).Raw(`select 
				post_klstw.id_post_klstw, post_klstw.content_post,
				post_klstw.created_at, dos.dosnama as nama, dos.foto
				from post_klstw
				join vw_kelas_tawar on cast(vw_kelas_tawar.id_master_kegiatan as character varying) = post_klstw.master_kegiatan_id
				join mk on mk.id_mata_kuliah = vw_kelas_tawar.id_mata_kuliah
				join users on users.name = post_klstw.author_id
				join dos on users.name = dos.dosid
				where mk.mkid = ? and kelas = ?;`, SubjectID, SubjectClass).
		Find(&announcement).Error

	return announcement, err
}

func (f *forumRepository) Forums(ctx context.Context, masterActivityID string) ([]entities.Announcement, error) {
	announcement := []entities.Announcement{}

	err := f.db.WithContext(ctx).Raw(`
			select forum.id_post_klstw, forum.content_post,
			forum.created_at, forum.nama, forum.foto, forum.author_id, jsonb_agg(post_materi.*) as post_materi, jsonb_agg(comment.*) as comments 
			from (select * from (select post_klstw.id_post_klstw, post_klstw.content_post,
			post_klstw.created_at, dos.dosnama as nama, dos.foto, post_klstw.author_id
			from post_klstw
			join vw_kelas_tawar on cast(vw_kelas_tawar.id_master_kegiatan as character varying) = post_klstw.master_kegiatan_id
			join mk on mk.id_mata_kuliah = vw_kelas_tawar.id_mata_kuliah
			join users on users.name = post_klstw.author_id
			join dos on users.name = dos.dosid
			where master_kegiatan_id = ?) as lecturer_forum
			UNION select * from (select 
			post_klstw.id_post_klstw, post_klstw.content_post,
			post_klstw.created_at, mhs.mhsnama as nama, mhs.foto, post_klstw.author_id
			from post_klstw
			join vw_kelas_tawar on cast(vw_kelas_tawar.id_master_kegiatan as character varying) = post_klstw.master_kegiatan_id
			join mk on mk.id_mata_kuliah = vw_kelas_tawar.id_mata_kuliah
			join users on users.name = post_klstw.author_id
			join mhs on users.name = mhs.mhsid
			where master_kegiatan_id = ?) as student_forum) as forum
			
			left join (
				select * from materi
				join post_materi on post_materi.materi_id = materi.materi_id
				where hidden_status = ?
				and deleted_at is null
			) post_materi on post_materi.post_klstw_id = forum.id_post_klstw

			left join (
			select * from (
				select post_comment.id_post_comment, post_comment.post_klstw_id, post_comment.content_comment,
				post_comment.created_at, dos.dosnama as nama, dos.foto, post_comment.author_id from post_comment
				join users on users.name = post_comment.author_id
				join dos on users.name = dos.dosid) as dos_comment
				UNION select * from (
				select post_comment.id_post_comment, post_comment.post_klstw_id, post_comment.content_comment,
				post_comment.created_at, mhs.mhsnama as nama, mhs.foto, post_comment.author_id from post_comment
				join users on users.name = post_comment.author_id
				join mhs on users.name = mhs.mhsid) as mhs_comment
			) comment on comment.post_klstw_id = forum.id_post_klstw

			group by (forum.id_post_klstw, forum.content_post,
				forum.created_at, forum.nama, forum.foto, forum.author_id)
			order by id_post_klstw DESC;
			`, masterActivityID, masterActivityID, 0).
		Find(&announcement).Error

	return announcement, err
}

func (f *forumRepository) StoreAnnouncement(ctx context.Context, announcement entities.AnnouncementStore) error {
	postJSON, err := json.Marshal(announcement.PostContent)

	if err != nil {
		return err
	}

	// return f.db.WithContext(ctx).Table("post_klstw").Create(&announcement).Error
	return f.db.WithContext(ctx).Exec("insert into post_klstw(master_kegiatan_id, content_post,author_id,flag_author,created_at,updated_at) values(?,?,?,?,?,?)", announcement.ActivityMasterId, postJSON, announcement.AuthorId, announcement.FlagAuthor, announcement.CreatedAt, announcement.UpdatedAt).Error
}

func (f *forumRepository) UpdateAnnouncement(ctx context.Context, announcement entities.AnnouncementUpdate) error {
	postJSON, err := json.Marshal(announcement.PostContent)

	if err != nil {
		return err
	}

	return f.db.WithContext(ctx).Exec("update post_klstw set content_post = ?, updated_at = ? WHERE id_post_klstw = ?", postJSON, announcement.UpdatedAt, announcement.AnnouncementID).Error
}

func (f *forumRepository) DeleteAnnouncement(ctx context.Context, announcementID int) error {
	return f.db.Transaction(func(tx *gorm.DB) error {
		err := f.db.WithContext(ctx).Table("post_comment").Where("post_klstw_id = ?", announcementID).Delete(&entities.Comment{}).Error

		if err != nil {
			return err
		}

		return f.db.WithContext(ctx).Exec("DELETE FROM post_klstw where id_post_klstw = ?", announcementID).Error
	})
}

func (f *forumRepository) StoreComment(ctx context.Context, comment entities.CommentStore) error {
	return f.db.WithContext(ctx).Table("post_comment").Create(&comment).Error
}

func (f *forumRepository) UpdateComment(ctx context.Context, comment entities.CommentUpdate) error {
	return f.db.WithContext(ctx).Table("post_comment").Omit("post_klstw_id", "id_post_comment", "author_id", "flag_author").Where("id_post_comment = ?", comment.CommentID).Updates(&comment).Error
}

func (f *forumRepository) DeleteComment(ctx context.Context, commentID int) error {
	return f.db.WithContext(ctx).Table("post_comment").Where("id_post_comment = ?", commentID).Delete(&entities.Comment{}).Error
}
