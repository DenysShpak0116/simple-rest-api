package teachers_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/teachers"
	"simple-rest-api/internal/handlers/teachers/mocks"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestDeleteTeacher(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		teacherID      string
		mockSetup      func(deleter *mocks.MockTeacherDeleter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Teacher Deleted",
			teacherID: "1",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
				deleter.EXPECT().DeleteTeacher(1).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:      "Error - Invalid Teacher ID",
			teacherID: "invalid",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Negative Teacher ID",
			teacherID: "-1",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Zero Teacher ID",
			teacherID: "0",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Teacher Not Found",
			teacherID: "99",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
				deleter.EXPECT().DeleteTeacher(99).Return(repository.ErrTeacherNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "teacher not found",
		},
		{
			name:      "Error - Database Error",
			teacherID: "1",
			mockSetup: func(deleter *mocks.MockTeacherDeleter) {
				deleter.EXPECT().DeleteTeacher(1).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to delete teacher",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockTeacherDeleter(t)

			tc.mockSetup(mockDeleter)

			handler := teachers.Delete(silentLogger, mockDeleter)

			req := httptest.NewRequest(http.MethodDelete, "/teachers/"+tc.teacherID, nil)
			req.SetPathValue("id", tc.teacherID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			if tc.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
			}
		})
	}
}
