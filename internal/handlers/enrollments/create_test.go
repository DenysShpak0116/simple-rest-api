package enrollments_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/enrollments"
	"simple-rest-api/internal/handlers/enrollments/mocks"
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEnrollment(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		studentID      string
		courseID       string
		mockSetup      func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Enrollment Created",
			studentID: "1",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
				studentValidator.EXPECT().GetStudentByID(1).Return(&models.Student{ID: 1}, nil)
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{ID: 1}, nil)
				creator.EXPECT().CreateEnrollment(mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "enrollment created successfully",
		},
		{
			name:      "Error - Invalid Student ID",
			studentID: "invalid",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Negative Student ID",
			studentID: "-1",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Zero Student ID",
			studentID: "0",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Invalid Course ID",
			studentID: "1",
			courseID:  "invalid",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:      "Error - Negative Course ID",
			studentID: "1",
			courseID:  "-1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:      "Error - Student Not Found",
			studentID: "99",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
				studentValidator.EXPECT().GetStudentByID(99).Return(nil, repository.ErrStudentNotFound)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "student not found",
		},
		{
			name:      "Error - Course Not Found",
			studentID: "1",
			courseID:  "99",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
				studentValidator.EXPECT().GetStudentByID(1).Return(&models.Student{ID: 1}, nil)
				courseValidator.EXPECT().GetCourseByID(99).Return(nil, repository.ErrCourseNotFound)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "course not found",
		},
		{
			name:      "Error - Student Already Enrolled",
			studentID: "1",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
				studentValidator.EXPECT().GetStudentByID(1).Return(&models.Student{ID: 1}, nil)
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{ID: 1}, nil)
				creator.EXPECT().CreateEnrollment(mock.Anything).Return(repository.ErrEnrollmentConflict)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "student is already enrolled in this course",
		},
		{
			name:      "Error - Database Error",
			studentID: "1",
			courseID:  "1",
			mockSetup: func(creator *mocks.MockEnrollmentCreator, studentValidator *mocks.MockStudentValidator, courseValidator *mocks.MockCourseValidator) {
				studentValidator.EXPECT().GetStudentByID(1).Return(&models.Student{ID: 1}, nil)
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{ID: 1}, nil)
				creator.EXPECT().CreateEnrollment(mock.Anything).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create enrollment",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCreator := mocks.NewMockEnrollmentCreator(t)
			mockStudentValidator := mocks.NewMockStudentValidator(t)
			mockCourseValidator := mocks.NewMockCourseValidator(t)

			tc.mockSetup(mockCreator, mockStudentValidator, mockCourseValidator)

			handler := enrollments.Create(silentLogger, mockCreator, mockStudentValidator, mockCourseValidator)

			req := httptest.NewRequest(http.MethodPost, "/students/"+tc.studentID+"/enrollments/"+tc.courseID, nil)
			req.SetPathValue("id", tc.studentID)
			req.SetPathValue("course_id", tc.courseID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
