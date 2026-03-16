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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCourse(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    map[string]any
		mockSetup      func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Course Created",
			requestBody: map[string]any{
				"title":       "Introduction to Go",
				"description": "A comprehensive guide to Go",
				"teacher_id":  1,
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
				validator.EXPECT().GetTeacherByID(1).Return(&models.Teacher{ID: 1}, nil)
				creator.EXPECT().CreateCourse(mock.Anything).Return(42, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":42`,
		},
		{
			name: "Error - Validation Failed (Title too short)",
			requestBody: map[string]any{
				"title":       "A",
				"description": "A comprehensive guide to Go",
				"teacher_id":  1,
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name: "Error - Teacher Not Found",
			requestBody: map[string]any{
				"title":       "Introduction to Go",
				"description": "A comprehensive guide to Go",
				"teacher_id":  99,
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
				validator.EXPECT().GetTeacherByID(99).Return(nil, errors.New("teacher not found"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "teacher not found",
		},
		{
			name: "Error - Database Creation Failed",
			requestBody: map[string]any{
				"title":       "Introduction to Go",
				"description": "A comprehensive guide to Go",
				"teacher_id":  1,
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
				validator.EXPECT().GetTeacherByID(1).Return(&models.Teacher{ID: 1}, nil)
				creator.EXPECT().CreateCourse(mock.Anything).Return(0, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create course",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCourseCreator(t)
			mockValidator := mocks.NewMockTeacherValidator(t)

			tc.mockSetup(mockCreator, mockValidator)

			handler := courses.Create(silentLogger, mockCreator, mockValidator)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
