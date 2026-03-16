package repository

import (
	"database/sql"
	"fmt"
	"simple-rest-api/internal/models"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(course *models.Course) (int, error) {
	query := `
		INSERT INTO courses (title, description, teacher_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(query, course.Title, course.Description, course.TeacherID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create course: %w", err)
	}

	return id, nil
}

func (r *CourseRepository) GetCourseByID(id int) (*models.Course, error) {
	query := `
		SELECT id, title, description, teacher_id
		FROM courses
		WHERE id = $1
	`

	course := &models.Course{}
	var teacherID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(&course.ID, &course.Title, &course.Description, &teacherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCourseNotFound
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if teacherID.Valid {
		tid := int(teacherID.Int64)
		course.TeacherID = &tid
	}

	return course, nil
}

func (r *CourseRepository) GetAllCourses() ([]*models.Course, error) {
	query := `
		SELECT id, title, description, teacher_id
		FROM courses
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		var teacherID sql.NullInt64
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &teacherID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		if teacherID.Valid {
			tid := int(teacherID.Int64)
			course.TeacherID = &tid
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating courses: %w", err)
	}

	return courses, nil
}

func (r *CourseRepository) UpdateCourse(course *models.Course) error {
	query := `
		UPDATE courses
		SET title = $1, description = $2, teacher_id = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(query, course.Title, course.Description, course.TeacherID, course.ID)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCourseNotFound
	}

	return nil
}

func (r *CourseRepository) DeleteCourse(id int) error {
	query := `
		DELETE FROM courses
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCourseNotFound
	}

	return nil
}

func (r *CourseRepository) GetCoursesByTeacherID(teacherID int) ([]*models.Course, error) {
	query := `
		SELECT id, title, description, teacher_id
		FROM courses
		WHERE teacher_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses for teacher: %w", err)
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		course := &models.Course{}
		var tid sql.NullInt64
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &tid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		if tid.Valid {
			teacherid := int(tid.Int64)
			course.TeacherID = &teacherid
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating courses: %w", err)
	}

	return courses, nil
}
