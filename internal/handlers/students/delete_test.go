package students_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/students"
	"simple-rest-api/internal/handlers/students/mocks"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestDeleteStudent(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		studentID      string
		mockSetup      func(deleter *mocks.MockStudentDeleter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Student Deleted",
			studentID: "1",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
				deleter.EXPECT().DeleteStudent(1).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:      "Error - Invalid Student ID",
			studentID: "invalid",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Negative Student ID",
			studentID: "-1",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Zero Student ID",
			studentID: "0",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Student Not Found",
			studentID: "99",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
				deleter.EXPECT().DeleteStudent(99).Return(repository.ErrStudentNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "student not found",
		},
		{
			name:      "Error - Database Error",
			studentID: "1",
			mockSetup: func(deleter *mocks.MockStudentDeleter) {
				deleter.EXPECT().DeleteStudent(1).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to delete student",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockDeleter := mocks.NewMockStudentDeleter(t)

			tc.mockSetup(mockDeleter)

			handler := students.Delete(silentLogger, mockDeleter)

			req := httptest.NewRequest(http.MethodDelete, "/students/"+tc.studentID, nil)
			req.SetPathValue("id", tc.studentID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			if tc.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
			}
		})
	}
}
