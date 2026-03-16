package courses_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/courses"
	"simple-rest-api/internal/handlers/courses/mocks"
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateCourse(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	teacherID := 2

	tests := []struct {
		name           string
		courseID       string
		requestBody    map[string]any
		mockSetup      func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success - Course Updated",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
				"teacher_id":  2,
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				teacherValidator.EXPECT().GetTeacherByID(2).Return(&models.Teacher{ID: 2}, nil)
				updater.EXPECT().UpdateCourse(mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "course updated successfully",
		},
		{
			name:     "Success - Course Updated Without Teacher Change",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				updater.EXPECT().UpdateCourse(mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "course updated successfully",
		},
		{
			name:     "Error - Invalid Course ID",
			courseID: "invalid",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Negative Course ID",
			courseID: "-1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Course Not Found During Validation",
			courseID: "99",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(99).Return(nil, repository.ErrCourseNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "course not found",
		},
		{
			name:     "Error - Validation Failed (Title too short)",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "A",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:     "Error - Teacher Not Found",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
				"teacher_id":  99,
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				teacherValidator.EXPECT().GetTeacherByID(99).Return(nil, repository.ErrTeacherNotFound)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "teacher not found",
		},
		{
			name:     "Error - Database Update Failed",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
				"teacher_id":  3,
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				teacherValidator.EXPECT().GetTeacherByID(3).Return(&models.Teacher{ID: 3}, nil)
				updater.EXPECT().UpdateCourse(mock.Anything).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to update course",
		},
		{
			name:     "Error - Course Validation Database Error",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to validate course",
		},
		{
			name:     "Error - Teacher Validation Database Error",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
				"teacher_id":  3,
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				teacherValidator.EXPECT().GetTeacherByID(3).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to validate teacher",
		},
		{
			name:     "Error - Course Not Found During Update",
			courseID: "1",
			requestBody: map[string]any{
				"title":       "Updated Course Title",
				"description": "Updated comprehensive guide",
				"teacher_id":  2,
			},
			mockSetup: func(updater *mocks.MockCourseUpdater, courseValidator *mocks.MockCourseValidator, teacherValidator *mocks.MockTeacherValidator) {
				courseValidator.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Original Title",
					Description: "Original description",
					TeacherID:   &teacherID,
				}, nil)
				teacherValidator.EXPECT().GetTeacherByID(2).Return(&models.Teacher{ID: 2}, nil)
				updater.EXPECT().UpdateCourse(mock.Anything).Return(repository.ErrCourseNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "course not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUpdater := mocks.NewMockCourseUpdater(t)
			mockCourseValidator := mocks.NewMockCourseValidator(t)
			mockTeacherValidator := mocks.NewMockTeacherValidator(t)

			tc.mockSetup(mockUpdater, mockCourseValidator, mockTeacherValidator)

			handler := courses.Update(silentLogger, mockUpdater, mockCourseValidator, mockTeacherValidator)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/courses/"+tc.courseID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tc.courseID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
