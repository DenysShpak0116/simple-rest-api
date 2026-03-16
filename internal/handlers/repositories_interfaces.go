package httpserver

import (
	"simple-rest-api/internal/handlers/courses"
	"simple-rest-api/internal/handlers/enrollments"
	"simple-rest-api/internal/handlers/students"
	"simple-rest-api/internal/handlers/teachers"
)

type StudentRepository interface {
	students.StudentCreator
	students.StudentGetter
	students.StudentsLister
	students.StudentUpdater
	students.StudentDeleter
}

type TeacherRepository interface {
	teachers.TeacherCreator
	teachers.TeacherGetter
	teachers.TeachersLister
	teachers.TeacherUpdater
	teachers.TeacherDeleter
}

type CourseRepository interface {
	courses.CourseCreator
	courses.CourseGetter
	courses.CoursesLister
	courses.CourseUpdater
	courses.CourseDeleter
}

type EnrollmentRepository interface {
	enrollments.EnrollmentCreator
	enrollments.EnrollmentDeleter
}
