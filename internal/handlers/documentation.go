package httpserver

import "time"

type ApiResponse[T any] struct {
	Message string `json:"message" example:"Operation successful"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Resource not found"`
}

type ValidationError map[string]interface{}

type StudentResponse struct {
	ID        int       `json:"id" example:"1"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStudentRequest struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
}

type TeacherResponse struct {
	ID         int       `json:"id" example:"1"`
	FirstName  string    `json:"first_name" example:"John"`
	LastName   string    `json:"last_name" example:"Smith"`
	Department string    `json:"department" example:"Mathematics"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTeacherRequest struct {
	FirstName  string `json:"first_name" example:"John"`
	LastName   string `json:"last_name" example:"Smith"`
	Department string `json:"department" example:"Mathematics"`
}

type CourseResponse struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Introduction to Mathematics"`
	Description string    `json:"description" example:"Basic mathematics concepts"`
	TeacherID   int       `json:"teacher_id" example:"1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCourseRequest struct {
	Title       string `json:"title" example:"Introduction to Mathematics"`
	Description string `json:"description" example:"Basic mathematics concepts"`
	TeacherID   int    `json:"teacher_id" example:"1"`
}

type EnrollmentResponse struct {
	StudentID  int       `json:"student_id" example:"1"`
	CourseID   int       `json:"course_id" example:"1"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

// @title           Simple REST API
// @version         1.0
// @description     This API provides endpoints for managing students, teachers, courses, and enrollments.
// @host            localhost:8081
// @BasePath        /

// _listStudents godoc
// @Summary      List all students
// @Description  Retrieve a list of all students in the system
// @Tags         Students
// @Produce      json
// @Success      200  {object}  ApiResponse[[]StudentResponse]
// @Failure      500  {object}  ErrorResponse
// @Router       /students [get]
func _listStudents() {}

// _getStudent godoc
// @Summary      Get a student by ID
// @Description  Retrieve details of a specific student
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  ApiResponse[StudentResponse]
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /students/{id} [get]
func _getStudent() {}

// _createStudent godoc
// @Summary      Create a new student
// @Description  Create a new student with provided information
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        request body CreateStudentRequest true "Student creation request"
// @Success      201  {object}  ApiResponse[StudentResponse]
// @Failure      400  {object}  ValidationError
// @Failure      500  {object}  ErrorResponse
// @Router       /students [post]
func _createStudent() {}

// _updateStudent godoc
// @Summary      Update a student
// @Description  Update an existing student's information
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Student ID"
// @Param        request body      CreateStudentRequest true  "Student update request"
// @Success      200     {object}  ApiResponse[StudentResponse]
// @Failure      400     {object}  ValidationError
// @Failure      404     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /students/{id} [put]
func _updateStudent() {}

// _deleteStudent godoc
// @Summary      Delete a student
// @Description  Delete a student by ID
// @Tags         Students
// @Param        id   path      int  true  "Student ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /students/{id} [delete]
func _deleteStudent() {}

// _listTeachers godoc
// @Summary      List all teachers
// @Description  Retrieve a list of all teachers in the system
// @Tags         Teachers
// @Produce      json
// @Success      200  {object}  ApiResponse[[]TeacherResponse]
// @Failure      500  {object}  ErrorResponse
// @Router       /teachers [get]
func _listTeachers() {}

// _getTeacher godoc
// @Summary      Get a teacher by ID
// @Description  Retrieve details of a specific teacher
// @Tags         Teachers
// @Produce      json
// @Param        id   path      int  true  "Teacher ID"
// @Success      200  {object}  ApiResponse[TeacherResponse]
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /teachers/{id} [get]
func _getTeacher() {}

// _createTeacher godoc
// @Summary      Create a new teacher
// @Description  Create a new teacher with provided information
// @Tags         Teachers
// @Accept       json
// @Produce      json
// @Param        request body CreateTeacherRequest true "Teacher creation request"
// @Success      201  {object}  ApiResponse[TeacherResponse]
// @Failure      400  {object}  ValidationError
// @Failure      500  {object}  ErrorResponse
// @Router       /teachers [post]
func _createTeacher() {}

// _updateTeacher godoc
// @Summary      Update a teacher
// @Description  Update an existing teacher's information
// @Tags         Teachers
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Teacher ID"
// @Param        request body      CreateTeacherRequest true  "Teacher update request"
// @Success      200     {object}  ApiResponse[TeacherResponse]
// @Failure      400     {object}  ValidationError
// @Failure      404     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /teachers/{id} [put]
func _updateTeacher() {}

// _deleteTeacher godoc
// @Summary      Delete a teacher
// @Description  Delete a teacher by ID
// @Tags         Teachers
// @Param        id   path      int  true  "Teacher ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /teachers/{id} [delete]
func _deleteTeacher() {}

// _listCourses godoc
// @Summary      List all courses
// @Description  Retrieve a list of all courses in the system
// @Tags         Courses
// @Produce      json
// @Success      200  {object}  ApiResponse[[]CourseResponse]
// @Failure      500  {object}  ErrorResponse
// @Router       /courses [get]
func _listCourses() {}

// _getCourse godoc
// @Summary      Get a course by ID
// @Description  Retrieve details of a specific course
// @Tags         Courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  ApiResponse[CourseResponse]
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /courses/{id} [get]
func _getCourse() {}

// _createCourse godoc
// @Summary      Create a new course
// @Description  Create a new course with provided information. Validates that the teacher_id exists.
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        request body CreateCourseRequest true "Course creation request"
// @Success      201  {object}  ApiResponse[CourseResponse]
// @Failure      400  {object}  ValidationError "Validation error or teacher not found"
// @Failure      500  {object}  ErrorResponse
// @Router       /courses [post]
func _createCourse() {}

// _updateCourse godoc
// @Summary      Update a course
// @Description  Update an existing course's information. Validates that the course exists and that the teacher_id exists.
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        id      path      int                 true  "Course ID"
// @Param        request body      CreateCourseRequest true  "Course update request"
// @Success      200     {object}  ApiResponse[CourseResponse]
// @Failure      400     {object}  ValidationError "Validation error or teacher not found"
// @Failure      404     {object}  ErrorResponse "Course not found"
// @Failure      500     {object}  ErrorResponse
// @Router       /courses/{id} [put]
func _updateCourse() {}

// _deleteCourse godoc
// @Summary      Delete a course
// @Description  Delete a course by ID
// @Tags         Courses
// @Param        id   path      int  true  "Course ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /courses/{id} [delete]
func _deleteCourse() {}

// _enrollStudent godoc
// @Summary      Enroll a student in a course
// @Description  Add a student enrollment in a course. Validates that both student and course exist.
// @Tags         Enrollments
// @Produce      json
// @Param        id         path      int  true  "Student ID"
// @Param        course_id  path      int  true  "Course ID"
// @Success      201        {object}  ApiResponse[EnrollmentResponse]
// @Failure      400        {object}  ErrorResponse "Student or course not found"
// @Failure      409        {object}  ErrorResponse "Student already enrolled in course"
// @Failure      500        {object}  ErrorResponse
// @Router       /students/{id}/courses/{course_id} [post]
func _enrollStudent() {}

// _unenrollStudent godoc
// @Summary      Remove student enrollment from a course
// @Description  Delete an enrollment, removing a student from a course. Validates that both student and course exist.
// @Tags         Enrollments
// @Param        id         path      int  true  "Student ID"
// @Param        course_id  path      int  true  "Course ID"
// @Success      204        "No Content"
// @Failure      400        {object}  ErrorResponse "Student or course not found"
// @Failure      404        {object}  ErrorResponse "Enrollment not found"
// @Failure      500        {object}  ErrorResponse
// @Router       /students/{id}/courses/{course_id} [delete]
func _unenrollStudent() {}
