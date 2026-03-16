package students_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateStudent(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		studentID      string
		requestBody    map[string]any
		mockSetup      func(updater *mocks.MockStudentUpdater)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Student Updated",
			studentID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
				updater.EXPECT().UpdateStudent(mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "student updated successfully",
		},
		{
			name:      "Error - Invalid Student ID",
			studentID: "invalid",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Negative Student ID",
			studentID: "-1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid student id",
		},
		{
			name:      "Error - Validation Failed (First Name too short)",
			studentID: "1",
			requestBody: map[string]any{
				"first_name": "J",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:      "Error - Invalid Email",
			studentID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "invalid-email",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:      "Error - Student Not Found",
			studentID: "99",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
				updater.EXPECT().UpdateStudent(mock.Anything).Return(repository.ErrStudentNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "student not found",
		},
		{
			name:      "Error - Database Update Failed",
			studentID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane@example.com",
			},
			mockSetup: func(updater *mocks.MockStudentUpdater) {
				updater.EXPECT().UpdateStudent(mock.Anything).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to update student",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUpdater := mocks.NewMockStudentUpdater(t)

			tc.mockSetup(mockUpdater)

			handler := students.Update(silentLogger, mockUpdater)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/students/"+tc.studentID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tc.studentID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
