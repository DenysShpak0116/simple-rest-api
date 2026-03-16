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

	"github.com/stretchr/testify/assert"
)

func TestListTeachers(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		mockSetup      func(lister *mocks.MockTeachersLister)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Teachers Retrieved",
			mockSetup: func(lister *mocks.MockTeachersLister) {
				lister.EXPECT().GetAllTeachers().Return([]*models.Teacher{
					{
						ID:         1,
						FirstName:  "John",
						LastName:   "Smith",
						Department: "Mathematics",
					},
					{
						ID:         2,
						FirstName:  "Jane",
						LastName:   "Doe",
						Department: "Physics",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "teachers retrieved successfully",
		},
		{
			name: "Success - No Teachers",
			mockSetup: func(lister *mocks.MockTeachersLister) {
				lister.EXPECT().GetAllTeachers().Return([]*models.Teacher{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "teachers retrieved successfully",
		},
		{
			name: "Error - Database Error",
			mockSetup: func(lister *mocks.MockTeachersLister) {
				lister.EXPECT().GetAllTeachers().Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve teachers",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLister := mocks.NewMockTeachersLister(t)

			tc.mockSetup(mockLister)

			handler := teachers.List(silentLogger, mockLister)

			req := httptest.NewRequest(http.MethodGet, "/teachers", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
