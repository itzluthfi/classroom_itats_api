package student_repositories

import (
	"classroom_itats_api/entities"
	"context"

	"gorm.io/gorm"
)

type studentProfileRepository struct {
	db *gorm.DB
}

type StudentProfileRepository interface {
	GetStudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfile, error)
	GetStudentSubjectPresence(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubjectPresence, error)
}

func NewStudentProfileRepository(db *gorm.DB) *studentProfileRepository {
	return &studentProfileRepository{db: db}
}

func (s *studentProfileRepository) GetStudentProfile(ctx context.Context, mhsID string, pakID string) (entities.StudentProfile, error) {
	studentProfile := entities.StudentProfile{}

	err := s.db.WithContext(ctx).Raw(`
		select mhs.mhsid, mhs.mhsnama, mhs.mobile, mhs.email, mhs.foto,
		(select count(absen.*) from absen where absen.mhsid = ? and absen.pakid = ?) as absen_mhs,
		(select SUM(vw_kelas_tawar.target_pertemuan)) + 2 as total_absen,
		(select count((select count(tugas_submission.*) from tugas_submission join tugas_kul on tugas_kul.id_tugas_kul = tugas_submission.tugas_kul_id where tugas_submission.mhsid = ? and tugas_kul.master_kegiatan_id = vw_kelas_tawar.id_master_kegiatan))) as tugas_terkumpul,
		(select count((select count(tugas_kul.*) from tugas_kul where tugas_kul.master_kegiatan_id = vw_kelas_tawar.id_master_kegiatan))) as total_tugas

		from mhs 
		join krs on mhs.mhsid = krs.mhsid
		join vw_kelas_tawar on vw_kelas_tawar.mkid = krs.mkid and vw_kelas_tawar.kelas = krs.kelaskrs and vw_kelas_tawar.pakid = krs.pakid and vw_kelas_tawar.jurid = krs.jurid
		where mhs.mhsid = ? and krs.pakid = ? group by mhs.mhsid
	`, mhsID, pakID, mhsID, mhsID, pakID).Find(&studentProfile).Error

	return studentProfile, err
}

func (s *studentProfileRepository) GetStudentSubjectPresence(ctx context.Context, mhsID string, pakID string) ([]entities.StudentSubjectPresence, error) {
	studentSubjectPresences := []entities.StudentSubjectPresence{}

	err := s.db.WithContext(ctx).Raw(`
		select vw_kelas_tawar.mkid, vw_kelas_tawar.mknama, vw_kelas_tawar.kelas, vw_kelas_tawar.id_master_kegiatan,
		(select count(absen.*) from absen where absen.mhsid = ? and absen.pakid = ? and absen.mkid = vw_kelas_tawar.mkid and absen.kelas = vw_kelas_tawar.kelas) as absen_mhs,
		(select SUM(vw_kelas_tawar.target_pertemuan)) + 2 as total_absen

		from krs
		join mhs on mhs.mhsid = krs.mhsid
		join vw_kelas_tawar on vw_kelas_tawar.mkid = krs.mkid and vw_kelas_tawar.kelas = krs.kelaskrs and vw_kelas_tawar.pakid = krs.pakid and vw_kelas_tawar.jurid = krs.jurid
		where mhs.mhsid = ? and krs.pakid = ? group by vw_kelas_tawar.mkid, vw_kelas_tawar.mknama, vw_kelas_tawar.kelas, vw_kelas_tawar.id_master_kegiatan
	`, mhsID, pakID, mhsID, pakID).Find(&studentSubjectPresences).Error

	return studentSubjectPresences, err
}
