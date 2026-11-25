package services

import (
	"errors"
	"strings"

	"pkl-attendance/backend/models"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

type StudentManagementService struct {
	users       *repositories.UserRepository
	students    *repositories.StudentRepository
	attendances *repositories.AttendanceRepository
	logs        *ActivityLogService
}

func NewStudentManagementService(
	u *repositories.UserRepository,
	s *repositories.StudentRepository,
	a *repositories.AttendanceRepository,
	l *ActivityLogService,
) *StudentManagementService {
	return &StudentManagementService{users: u, students: s, attendances: a, logs: l}
}

type CreateStudentInput struct {
	Username     string
	Password     string
	Name         string
	NIS          string
	Class        string
	Major        string
	Phone        string
	Address      string
	PKLStartDate string
	PKLEndDate   string
}

func (s *StudentManagementService) Create(adminID uint, in CreateStudentInput) (*models.User, error) {
	username := strings.TrimSpace(in.Username)
	name := strings.TrimSpace(in.Name)
	if len(username) < 3 {
		return nil, errors.New("username minimal 3 karakter")
	}
	if len(in.Password) < 6 {
		return nil, errors.New("password minimal 6 karakter")
	}
	if name == "" {
		return nil, errors.New("nama wajib diisi")
	}
	exists, _ := s.users.UsernameExists(username)
	if exists {
		return nil, errors.New("username sudah digunakan")
	}

	var start, end *models.DateOnly
	if in.PKLStartDate != "" {
		d, err := models.ParseDateOnly(in.PKLStartDate)
		if err != nil {
			return nil, errors.New("pkl_start_date harus YYYY-MM-DD")
		}
		start = &d
	}
	if in.PKLEndDate != "" {
		d, err := models.ParseDateOnly(in.PKLEndDate)
		if err != nil {
			return nil, errors.New("pkl_end_date harus YYYY-MM-DD")
		}
		end = &d
	}
	if start != nil && end != nil && end.Time().Before(start.Time()) {
		return nil, errors.New("pkl_end_date tidak boleh sebelum pkl_start_date")
	}

	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Username: username,
		Password: hash,
		Name:     name,
		Role:     models.RoleSiswa,
		IsActive: true,
	}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	st := &models.Student{
		UserID:       u.ID,
		NIS:          strings.TrimSpace(in.NIS), // <-- NIS DISIMPAN KE DB
		Class:        strings.TrimSpace(in.Class),
		Major:        strings.TrimSpace(in.Major),
		Phone:        strings.TrimSpace(in.Phone),
		Address:      strings.TrimSpace(in.Address),
		PKLStartDate: start,
		PKLEndDate:   end,
	}
	if err := s.students.Create(st); err != nil {
		return nil, err
	}
	s.logs.Log(adminID, "create_student", "Membuat akun siswa "+username, map[string]interface{}{"student_user_id": u.ID})
	u.Student = st
	return u, nil
}

type UpdateStudentInput struct {
	UserID       uint
	Name         string
	NIS          string
	Class        string
	Major        string
	Phone        string
	Address      string
	PKLStartDate string
	PKLEndDate   string
}

func (s *StudentManagementService) Update(adminID uint, in UpdateStudentInput) (*models.User, error) {
	u, err := s.users.FindByID(in.UserID)
	if err != nil {
		return nil, errors.New("siswa tidak ditemukan")
	}
	if u.Role != models.RoleSiswa {
		return nil, errors.New("hanya dapat mengubah akun siswa")
	}
	st, err := s.students.FindByUserID(in.UserID)
	if err != nil {
		return nil, errors.New("profil siswa tidak ditemukan")
	}

	if name := strings.TrimSpace(in.Name); name != "" {
		u.Name = name
	}
	st.NIS = strings.TrimSpace(in.NIS)
	st.Class = strings.TrimSpace(in.Class)
	st.Major = strings.TrimSpace(in.Major)
	st.Phone = strings.TrimSpace(in.Phone)
	st.Address = strings.TrimSpace(in.Address)

	if in.PKLStartDate != "" {
		d, err := models.ParseDateOnly(in.PKLStartDate)
		if err != nil {
			return nil, errors.New("pkl_start_date harus YYYY-MM-DD")
		}
		st.PKLStartDate = &d
	}
	if in.PKLEndDate != "" {
		d, err := models.ParseDateOnly(in.PKLEndDate)
		if err != nil {
			return nil, errors.New("pkl_end_date harus YYYY-MM-DD")
		}
		st.PKLEndDate = &d
	}
	if st.PKLStartDate != nil && st.PKLEndDate != nil && st.PKLEndDate.Time().Before(st.PKLStartDate.Time()) {
		return nil, errors.New("pkl_end_date tidak boleh sebelum pkl_start_date")
	}

	if err := s.users.Update(u); err != nil {
		return nil, err
	}
	if err := s.students.Update(st); err != nil {
		return nil, err
	}
	s.logs.Log(adminID, "update_student", "Mengubah data siswa "+u.Username, map[string]interface{}{"student_user_id": u.ID})
	u.Student = st
	return u, nil
}

func (s *StudentManagementService) ResetPassword(adminID, userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password baru minimal 6 karakter")
	}
	u, err := s.users.FindByID(userID)
	if err != nil || u.Role != models.RoleSiswa {
		return errors.New("siswa tidak ditemukan")
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.Password = hash
	if err := s.users.Update(u); err != nil {
		return err
	}
	s.logs.Log(adminID, "reset_password", "Reset password siswa "+u.Username, map[string]interface{}{"student_user_id": u.ID})
	return nil
}

func (s *StudentManagementService) ToggleActive(adminID, userID uint, active bool) error {
	u, err := s.users.FindByID(userID)
	if err != nil || u.Role != models.RoleSiswa {
		return errors.New("siswa tidak ditemukan")
	}
	u.IsActive = active
	if err := s.users.Update(u); err != nil {
		return err
	}
	action := "disable_student"
	desc := "Menonaktifkan akun siswa "
	if active {
		action = "enable_student"
		desc = "Mengaktifkan akun siswa "
	}
	s.logs.Log(adminID, action, desc+u.Username, map[string]interface{}{"student_user_id": u.ID})
	return nil
}

func (s *StudentManagementService) Delete(adminID, userID uint) error {
	u, err := s.users.FindByID(userID)
	if err != nil || u.Role != models.RoleSiswa {
		return errors.New("siswa tidak ditemukan")
	}
	if err := s.users.DeleteWithAssociations(userID); err != nil {
		return err
	}
	s.logs.Log(adminID, "delete_student", "Menghapus akun siswa "+u.Username, map[string]interface{}{"student_user_id": userID})
	return nil
}

func (s *StudentManagementService) List() ([]models.User, error) {
	return s.users.ListStudents()
}