package repository

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrStudentNotFound    = errors.New("student not found")
	ErrTeacherNotFound    = errors.New("teacher not found")
	ErrCourseNotFound     = errors.New("course not found")
	ErrEnrollmentNotFound = errors.New("enrollment not found")

	ErrConflict             = errors.New("resource already exists")
	ErrEnrollmentConflict   = errors.New("student already enrolled in this course")
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrForeignKeyConstraint = errors.New("foreign key constraint violation")

	ErrInvalidID = errors.New("invalid id")
)
