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
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestGetTeacher(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		teacherID      string
		mockSetup      func(getter *mocks.MockTeacherGetter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success - Teacher Retrieved",
			teacherID: "1",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
				getter.EXPECT().GetTeacherByID(1).Return(&models.Teacher{
					ID:         1,
					FirstName:  "John",
					LastName:   "Smith",
					Department: "Mathematics",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":1`,
		},
		{
			name:      "Error - Invalid Teacher ID",
			teacherID: "invalid",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Negative Teacher ID",
			teacherID: "-5",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Zero Teacher ID",
			teacherID: "0",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid teacher id",
		},
		{
			name:      "Error - Teacher Not Found",
			teacherID: "99",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
				getter.EXPECT().GetTeacherByID(99).Return(nil, repository.ErrTeacherNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "teacher not found",
		},
		{
			name:      "Error - Database Error",
			teacherID: "1",
			mockSetup: func(getter *mocks.MockTeacherGetter) {
				getter.EXPECT().GetTeacherByID(1).Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve teacher",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := mocks.NewMockTeacherGetter(t)

			tc.mockSetup(mockGetter)

			handler := teachers.Get(silentLogger, mockGetter)

			req := httptest.NewRequest(http.MethodGet, "/teachers/"+tc.teacherID, nil)
			req.SetPathValue("id", tc.teacherID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
