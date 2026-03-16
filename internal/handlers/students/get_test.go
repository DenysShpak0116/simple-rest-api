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
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetStudent(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		studentID      string
		mockSetup      func(getter *mocks.MockStudentGetter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Student Retrieved",
			studentID: "1",
			mockSetup: func(getter *mocks.MockStudentGetter) {
				getter.EXPECT().GetStudentByID(1).Return(&models.Student{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":1`,
		},
		{
			name:      "Error - Invalid Student ID",
			studentID: "invalid",
			mockSetup: func(getter *mocks.MockStudentGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Negative Student ID",
			studentID: "-5",
			mockSetup: func(getter *mocks.MockStudentGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Zero Student ID",
			studentID: "0",
			mockSetup: func(getter *mocks.MockStudentGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Student Not Found",
			studentID: "99",
			mockSetup: func(getter *mocks.MockStudentGetter) {
				getter.EXPECT().GetStudentByID(99).Return(nil, repository.ErrStudentNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "student not found",
		},
		{
			name:      "Error - Database Error",
			studentID: "1",
			mockSetup: func(getter *mocks.MockStudentGetter) {
				getter.EXPECT().GetStudentByID(1).Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve student",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := mocks.NewMockStudentGetter(t)

			tc.mockSetup(mockGetter)

			handler := students.Get(silentLogger, mockGetter)

			req := httptest.NewRequest(http.MethodGet, "/students/"+tc.studentID, nil)
			req.SetPathValue("id", tc.studentID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
