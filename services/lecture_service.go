package services

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	"context"
)

type lectureService struct {
	lectureRepository repositories.LectureRepository
}

type LectureService interface {
	GetLecture(ctx context.Context, pakID string, subjectID string, class string) ([]entities.Lecture, error)
	GetStudentLecture(ctx context.Context, pakID string, subjectID, class string) (entities.StudentLectureJSON, error)
	GetLecturerLecture(ctx context.Context, pakID string, subjectID string, class string) ([][]entities.Lecture, error)
}

func NewLectureService(lectureRepository repositories.LectureRepository) *lectureService {
	return &lectureService{
		lectureRepository: lectureRepository,
	}
}

func (l *lectureService) GetLecture(ctx context.Context, pakID string, subjectID string, class string) ([]entities.Lecture, error) {
	return l.lectureRepository.GetLecture(ctx, pakID, subjectID, class)
}

func (l *lectureService) GetLecturerLecture(ctx context.Context, pakID string, subjectID string, class string) ([][]entities.Lecture, error) {
	lectures, err := l.GetLecture(ctx, pakID, subjectID, class)

	if err != nil {
		return nil, err
	}

	materialLectures := []entities.Lecture{}
	allResponsiLectures := [][]entities.Lecture{}
	responsiLectures := []entities.Lecture{}
	allLectures := [][]entities.Lecture{}

	for key, v := range lectures {
		if key > 0 {
			if v.LectureType != lectures[key-1].LectureType {
				materialLectures = append(materialLectures, lectures[:key]...)
				responsiLectures = append(responsiLectures, lectures[key:]...)
				break
			}
		}
		// if v.LectureType != "M" {
		// }
	}

	if len(materialLectures) == 0 {
		materialLectures = append(materialLectures, lectures...)
	}

	temp := 0
	for key := range responsiLectures {
		if key+1 == len(responsiLectures) {
			allResponsiLectures = append(allResponsiLectures, responsiLectures[temp:])
		} else {
			if responsiLectures[key+1].WeekID == 1 {
				allResponsiLectures = append(allResponsiLectures, responsiLectures[:key])
				temp = key
			}
		}
	}

	for i := 0; i < len(materialLectures); i++ {
		flag := false
		temp := []entities.Lecture{}
		if len(allResponsiLectures) != 0 {
			for k, v := range allResponsiLectures {
				for j := 0; j < len(v); j++ {
					if materialLectures[i].WeekID == v[j].WeekID {
						if k == 0 {
							temp = append(temp, materialLectures[i])
						}
						temp = append(temp, v[j])
						allLectures = append(allLectures, temp)
						flag = true
						break
					}
				}
			}
		}
		if !flag {
			temp = append(temp, materialLectures[i])
			allLectures = append(allLectures, temp)
		}
	}

	return allLectures, nil
}

func (l *lectureService) GetStudentLecture(ctx context.Context, pakID string, subjectID, class string) (entities.StudentLectureJSON, error) {
	lectures, err := l.lectureRepository.GetLecture(ctx, pakID, subjectID, class)

	if err != nil {
		return entities.StudentLectureJSON{}, err
	}

	materialLectures := []entities.Lecture{}
	allResponsiLectures := [][]entities.Lecture{}
	responsiLectures := []entities.Lecture{}

	for key, v := range lectures {
		if key > 0 {
			if v.LectureType != lectures[key-1].LectureType {
				materialLectures = append(materialLectures, lectures[:key]...)
				responsiLectures = append(responsiLectures, lectures[key:]...)
				break
			}
		}
	}

	temp := 0
	for key := range responsiLectures {
		if key+1 == len(responsiLectures) {
			allResponsiLectures = append(allResponsiLectures, responsiLectures[temp:])
		} else {
			if responsiLectures[key+1].WeekID == 1 {
				allResponsiLectures = append(allResponsiLectures, responsiLectures[:key])
				temp = key
			}
		}
	}

	if len(materialLectures) == 0 {
		return entities.StudentLectureJSON{
			MaterialLectures: lectures,
			ResponsiLectures: allResponsiLectures,
		}, err
	}

	return entities.StudentLectureJSON{
		MaterialLectures: materialLectures,
		ResponsiLectures: allResponsiLectures,
	}, err
}
