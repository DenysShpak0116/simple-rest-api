package teachers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/teachers"
	"simple-rest-api/internal/handlers/teachers/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTeacher(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    map[string]any
		mockSetup      func(creator *mocks.MockTeacherCreator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Teacher Created",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Smith",
				"department": "Mathematics",
			},
			mockSetup: func(creator *mocks.MockTeacherCreator) {
				creator.EXPECT().CreateTeacher(mock.Anything).Return(1, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":1`,
		},
		{
			name: "Error - Validation Failed (First Name too short)",
			requestBody: map[string]any{
				"first_name": "J",
				"last_name":  "Smith",
				"department": "Mathematics",
			},
			mockSetup: func(creator *mocks.MockTeacherCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name: "Error - Missing Department",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Smith",
			},
			mockSetup: func(creator *mocks.MockTeacherCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name: "Error - Database Error",
			requestBody: map[string]any{
				"first_name": "John",
				"last_name":  "Smith",
				"department": "Mathematics",
			},
			mockSetup: func(creator *mocks.MockTeacherCreator) {
				creator.EXPECT().CreateTeacher(mock.Anything).Return(0, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to create teacher",
		},
		{
			name: "Error - Missing Last Name",
			requestBody: map[string]any{
				"first_name": "John",
				"department": "Mathematics",
			},
			mockSetup: func(creator *mocks.MockTeacherCreator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCreator := mocks.NewMockTeacherCreator(t)

			tc.mockSetup(mockCreator)

			handler := teachers.Create(silentLogger, mockCreator)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/teachers", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
