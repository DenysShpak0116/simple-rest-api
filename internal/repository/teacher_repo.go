package repository

import (
	"database/sql"
	"fmt"
	"simple-rest-api/internal/models"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) CreateTeacher(teacher *models.Teacher) (int, error) {
	query := `
		INSERT INTO teachers (first_name, last_name, department)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, teacher.FirstName, teacher.LastName, teacher.Department).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create teacher: %w", err)
	}

	return id, nil
}

func (r *TeacherRepository) GetTeacherByID(id int) (*models.Teacher, error) {
	query := `
		SELECT id, first_name, last_name, department
		FROM teachers
		WHERE id = $1
	`

	teacher := &models.Teacher{}
	err := r.db.QueryRow(query, id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTeacherNotFound
		}
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return teacher, nil
}

func (r *TeacherRepository) GetAllTeachers() ([]*models.Teacher, error) {
	query := `
		SELECT id, first_name, last_name, department
		FROM teachers
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teachers: %w", err)
	}
	defer rows.Close()

	var teachers []*models.Teacher
	for rows.Next() {
		teacher := &models.Teacher{}
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Department)
		if err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teachers: %w", err)
	}

	return teachers, nil
}

func (r *TeacherRepository) UpdateTeacher(teacher *models.Teacher) error {
	query := `
		UPDATE teachers
		SET first_name = $1, last_name = $2, department = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(query, teacher.FirstName, teacher.LastName, teacher.Department, teacher.ID)
	if err != nil {
		return fmt.Errorf("failed to update teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTeacherNotFound
	}

	return nil
}

func (r *TeacherRepository) DeleteTeacher(id int) error {
	query := `
		DELETE FROM teachers
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTeacherNotFound
	}

	return nil
}
