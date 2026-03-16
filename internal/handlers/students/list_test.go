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

	"github.com/stretchr/testify/assert"
)

func TestListStudents(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		mockSetup      func(lister *mocks.MockStudentsLister)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Students Retrieved",
			mockSetup: func(lister *mocks.MockStudentsLister) {
				lister.EXPECT().GetAllStudents().Return([]*models.Student{
					{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john@example.com",
					},
					{
						ID:        2,
						FirstName: "Jane",
						LastName:  "Smith",
						Email:     "jane@example.com",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "students retrieved successfully",
		},
		{
			name: "Success - No Students",
			mockSetup: func(lister *mocks.MockStudentsLister) {
				lister.EXPECT().GetAllStudents().Return([]*models.Student{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "students retrieved successfully",
		},
		{
			name: "Error - Database Error",
			mockSetup: func(lister *mocks.MockStudentsLister) {
				lister.EXPECT().GetAllStudents().Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve students",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLister := mocks.NewMockStudentsLister(t)

			tc.mockSetup(mockLister)

			handler := students.List(silentLogger, mockLister)

			req := httptest.NewRequest(http.MethodGet, "/students", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
