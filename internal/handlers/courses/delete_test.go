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
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestDeleteCourse(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		courseID       string
		mockSetup      func(deleter *mocks.MockCourseDeleter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success - Course Deleted",
			courseID: "1",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
				deleter.EXPECT().DeleteCourse(1).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:     "Error - Invalid Course ID",
			courseID: "invalid",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Negative Course ID",
			courseID: "-1",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Zero Course ID",
			courseID: "0",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid course id",
		},
		{
			name:     "Error - Course Not Found",
			courseID: "99",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
				deleter.EXPECT().DeleteCourse(99).Return(repository.ErrCourseNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "course not found",
		},
		{
			name:     "Error - Database Error",
			courseID: "1",
			mockSetup: func(deleter *mocks.MockCourseDeleter) {
				deleter.EXPECT().DeleteCourse(1).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to delete course",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockCourseDeleter(t)

			tc.mockSetup(mockDeleter)

			handler := courses.Delete(silentLogger, mockDeleter)

			req := httptest.NewRequest(http.MethodDelete, "/courses/"+tc.courseID, nil)
			req.SetPathValue("id", tc.courseID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			if tc.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
			}
		})
	}
}
