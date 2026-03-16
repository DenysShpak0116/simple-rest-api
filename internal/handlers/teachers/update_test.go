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
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateTeacher(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		teacherID      string
		requestBody    map[string]any
		mockSetup      func(updater *mocks.MockTeacherUpdater)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Teacher Updated",
			teacherID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
				updater.EXPECT().UpdateTeacher(mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "teacher updated successfully",
		},
		{
			name:      "Error - Invalid Teacher ID",
			teacherID: "invalid",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Negative Teacher ID",
			teacherID: "-1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Validation Failed (First Name too short)",
			teacherID: "1",
			requestBody: map[string]any{
				"first_name": "J",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:      "Error - Missing Department",
			teacherID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
		{
			name:      "Error - Teacher Not Found",
			teacherID: "99",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
				updater.EXPECT().UpdateTeacher(mock.Anything).Return(repository.ErrTeacherNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "teacher not found",
		},
		{
			name:      "Error - Database Update Failed",
			teacherID: "1",
			requestBody: map[string]any{
				"first_name": "Jane",
				"last_name":  "Doe",
				"department": "Physics",
			},
			mockSetup: func(updater *mocks.MockTeacherUpdater) {
				updater.EXPECT().UpdateTeacher(mock.Anything).Return(errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to update teacher",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUpdater := mocks.NewMockTeacherUpdater(t)

			tc.mockSetup(mockUpdater)

			handler := teachers.Update(silentLogger, mockUpdater)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/teachers/"+tc.teacherID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tc.teacherID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
