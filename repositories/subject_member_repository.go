package repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type subjectMemberRepository struct {
	db *gorm.DB
}

type SubjectMemberRepository interface {
	GetSubjectMember(ctx context.Context, pakID string, mkID string, class string, jurID string) ([]entities.SubjectMember, error)
}

func NewSubjectMemberRepository(db *gorm.DB) *subjectMemberRepository {
	return &subjectMemberRepository{
		db: db,
	}
}

func (s *subjectMemberRepository) GetSubjectMember(ctx context.Context, pakID string, mkID string, class string, jurID string) ([]entities.SubjectMember, error) {
	students := []entities.SubjectMember{}

	err := s.db.WithContext(ctx).
		Raw(`select mhs.mhsid as uid, mhs.mhsnama as nama, mhs.mobile as mobile, (select count(*) from absen where mhsid = mhs.mhsid and pakid = ? and mkid = ? and kelas = ? and kultype = 'M') as kehadiran from mhs 
			join krs on krs.mhsid = mhs.mhsid
			join mk on mk.mkid = krs.mkid
			where krs.pakid = ? and krs.mkid = ? and krs.kelaskrs = ? and krs.jurid = ?

			UNION select * from (
				select dos.dosid as uid, dos.dosnama as nama, dos.mobile as mobile, 0 as kehadiran from dos 
				join vw_jadwal on vw_jadwal.dosid = dos.dosid
				where vw_jadwal.pakid = ? and vw_jadwal.mkid = ? and vw_jadwal.kelas = ?
			) as dos
			`, pakID, mkID, class, pakID, mkID, class, jurID, pakID, mkID, class).Find(&students).
		Error

	return students, err
}
