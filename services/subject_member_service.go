package services

import (
	"classroom_itats_api/entities"
	"classroom_itats_api/repositories"
	"context"
)

type subjectMemberService struct {
	subjectMemberRepository repositories.SubjectMemberRepository
}

type SubjectMemberService interface {
	GetSubjectMember(ctx context.Context, pakID string, mkID string, class string, jurID string, name string, role string) ([]entities.SubjectMember, error)
}

func NewSubjectMemberService(subjectMemberRepository repositories.SubjectMemberRepository) *subjectMemberService {
	return &subjectMemberService{
		subjectMemberRepository: subjectMemberRepository,
	}
}

func (s *subjectMemberService) GetSubjectMember(ctx context.Context, pakID string, mkID string, class string, jurID string, name string, role string) ([]entities.SubjectMember, error) {
	data, err := s.subjectMemberRepository.GetSubjectMember(ctx, pakID, mkID, class, jurID)
	temp := data[len(data)-1]
	copy(data[1:], data)
	data[0] = temp

	if role == "Mahasiswa" {
		for i, v := range data {
			if v.UserID == name {
				data = append(data[:i], data[i+1:]...)

				data = append(data[:1], append([]entities.SubjectMember{v}, data[1:]...)...)

				break
			}
		}
	}

	return data, err
}
