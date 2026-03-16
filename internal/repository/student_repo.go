package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"simple-rest-api/internal/models"

	"github.com/lib/pq"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(student *models.Student) (int, error) {
	query := `
		INSERT INTO students (first_name, last_name, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, student.FirstName, student.LastName, student.Email).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, fmt.Errorf("failed to create student: %w", err)
	}

	return id, nil
}

func (r *StudentRepository) GetStudentByID(id int) (*models.Student, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM students
		WHERE id = $1
	`

	student := &models.Student{}
	err := r.db.QueryRow(query, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrStudentNotFound
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return student, nil
}

func (r *StudentRepository) GetAllStudents() ([]*models.Student, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM students
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %w", err)
	}
	defer rows.Close()

	var students []*models.Student
	for rows.Next() {
		student := &models.Student{}
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating students: %w", err)
	}

	return students, nil
}

func (r *StudentRepository) UpdateStudent(student *models.Student) error {
	query := `
		UPDATE students
		SET first_name = $1, last_name = $2, email = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(query, student.FirstName, student.LastName, student.Email, student.ID)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

func (r *StudentRepository) DeleteStudent(id int) error {
	query := `
		DELETE FROM students
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrStudentNotFound
	}

	return nil
}
