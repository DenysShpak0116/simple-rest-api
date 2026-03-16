package courses_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-rest-api/internal/handlers/courses"
	"simple-rest-api/internal/handlers/courses/mocks"
	"simple-rest-api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestListCourses(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	teacherID1 := 1
	teacherID2 := 2

	tests := []struct {
		name           string
		mockSetup      func(lister *mocks.MockCoursesLister)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Courses Retrieved",
			mockSetup: func(lister *mocks.MockCoursesLister) {
				lister.EXPECT().GetAllCourses().Return([]*models.Course{
					{
						ID:          1,
						Title:       "Introduction to Go",
						Description: "A comprehensive guide to Go",
						TeacherID:   &teacherID1,
					},
					{
						ID:          2,
						Title:       "Advanced Go",
						Description: "Advanced Go programming patterns",
						TeacherID:   &teacherID2,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "courses retrieved successfully",
		},
		{
			name: "Success - No Courses",
			mockSetup: func(lister *mocks.MockCoursesLister) {
				lister.EXPECT().GetAllCourses().Return([]*models.Course{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "courses retrieved successfully",
		},
		{
			name: "Error - Database Error",
			mockSetup: func(lister *mocks.MockCoursesLister) {
				lister.EXPECT().GetAllCourses().Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to retrieve courses",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLister := mocks.NewMockCoursesLister(t)

			tc.mockSetup(mockLister)

			handler := courses.List(silentLogger, mockLister)

			req := httptest.NewRequest(http.MethodGet, "/courses", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "status code mismatch")
			assert.Contains(t, rr.Body.String(), tc.expectedBody, "response body mismatch")
		})
	}
}
