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

func TestCreateStudent(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    map[string]any
		mockSetup      func(creator *mocks.MockStudentCreator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Student Created",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
				creator.EXPECT().CreateStudent(mock.Anything).Return(1, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":1`,
		},
		{
			name: "Error - Validation Failed (First Name too short)",
			requestBody: map[string]any{
				"first_name": "J",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name: "Error - Invalid Email",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "invalid-email",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name: "Error - Duplicate Email",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
				creator.EXPECT().CreateStudent(mock.Anything).Return(0, repository.ErrDuplicateEmail)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   "email already exists",
		},
		{
			name: "Error - Database Error",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
				creator.EXPECT().CreateStudent(mock.Anything).Return(0, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create student",
		},
		{
			name: "Error - Missing Last Name",
			requestBody: map[string]any{
				"first_name": "John",
				"email":      "john@example.com",
			},
			mockSetup: func(creator *mocks.MockStudentCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCreator := mocks.NewMockStudentCreator(t)

			tc.mockSetup(mockCreator)

			handler := students.Create(silentLogger, mockCreator)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
