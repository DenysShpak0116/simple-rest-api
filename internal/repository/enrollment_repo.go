package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"simple-rest-api/internal/models"

	"github.com/lib/pq"
)

type EnrollmentRepository struct {
	db *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment *models.Enrollment) error {
	query := `
		INSERT INTO enrollments (student_id, course_id)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(query, enrollment.StudentID, enrollment.CourseID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "enrollments_pkey") || strings.Contains(pqErr.Constraint, "primary") {
				return ErrEnrollmentConflict
			}
		}
		return fmt.Errorf("failed to create enrollment: %w", err)
	}

	return nil
}

func (r *EnrollmentRepository) GetEnrollmentByID(studentID, courseID int) (*models.Enrollment, error) {
	query := `
		SELECT student_id, course_id
		FROM enrollments
		WHERE student_id = $1 AND course_id = $2
	`

	enrollment := &models.Enrollment{}
	err := r.db.QueryRow(query, studentID, courseID).Scan(&enrollment.StudentID, &enrollment.CourseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrEnrollmentNotFound
		}
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) GetStudentEnrollments(studentID int) ([]*models.Enrollment, error) {
	query := `
		SELECT student_id, course_id
		FROM enrollments
		WHERE student_id = $1
		ORDER BY course_id DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		enrollment := &models.Enrollment{}
		err := rows.Scan(&enrollment.StudentID, &enrollment.CourseID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating enrollments: %w", err)
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) GetCourseEnrollments(courseID int) ([]*models.Enrollment, error) {
	query := `
		SELECT student_id, course_id
		FROM enrollments
		WHERE course_id = $1
		ORDER BY student_id DESC
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*models.Enrollment
	for rows.Next() {
		enrollment := &models.Enrollment{}
		err := rows.Scan(&enrollment.StudentID, &enrollment.CourseID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating enrollments: %w", err)
	}

	return enrollments, nil
}

func (r *EnrollmentRepository) DeleteEnrollment(studentID, courseID int) error {
	query := `
		DELETE FROM enrollments
		WHERE student_id = $1 AND course_id = $2
	`

	result, err := r.db.Exec(query, studentID, courseID)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrEnrollmentNotFound
	}

	return nil
}

func (r *EnrollmentRepository) IsStudentEnrolled(studentID, courseID int) (bool, error) {
	query := `
		SELECT 1
		FROM enrollments
		WHERE student_id = $1 AND course_id = $2
	`

	var exists int
	err := r.db.QueryRow(query, studentID, courseID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check enrollment: %w", err)
	}

	return true, nil
}
