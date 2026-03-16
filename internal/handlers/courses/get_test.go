package courses_test

import (
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
)

func TestGetCourse(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	teacherID := 5

	tests := []struct {
		name           string
		courseID       string
		mockSetup      func(getter *mocks.MockCourseGetter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success - Course Retrieved",
			courseID: "1",
			mockSetup: func(getter *mocks.MockCourseGetter) {
				getter.EXPECT().GetCourseByID(1).Return(&models.Course{
					ID:          1,
					Title:       "Introduction to Go",
					Description: "A comprehensive guide to Go",
					TeacherID:   &teacherID,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":1`,
		},
		{
			name:     "Error - Invalid Course ID",
			courseID: "invalid",
			mockSetup: func(getter *mocks.MockCourseGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Negative Course ID",
			courseID: "-5",
			mockSetup: func(getter *mocks.MockCourseGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Zero Course ID",
			courseID: "0",
			mockSetup: func(getter *mocks.MockCourseGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Course Not Found",
			courseID: "99",
			mockSetup: func(getter *mocks.MockCourseGetter) {
				getter.EXPECT().GetCourseByID(99).Return(nil, repository.ErrCourseNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "course not found",
		},
		{
			name:     "Error - Database Error",
			courseID: "1",
			mockSetup: func(getter *mocks.MockCourseGetter) {
				getter.EXPECT().GetCourseByID(1).Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve course",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := mocks.NewMockCourseGetter(t)

			tc.mockSetup(mockGetter)

			handler := courses.Get(silentLogger, mockGetter)

			req := httptest.NewRequest(http.MethodGet, "/courses/"+tc.courseID, nil)
			req.SetPathValue("id", tc.courseID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
