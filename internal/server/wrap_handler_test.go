package server

import (
	"errors"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRedirectError_Error(t *testing.T) {
	assert.Equal(t, "redirect to ", Redirect("").Error())
	assert.Equal(t, "redirect to test-url", Redirect("test-url").Error())
}

func TestRedirectError_To(t *testing.T) {
	assert.Equal(t, "", Redirect("").To())
	assert.Equal(t, "test-url", Redirect("test-url").To())
}

func TestStatusError_Code(t *testing.T) {
	assert.Equal(t, 0, StatusError(0).Code())
	assert.Equal(t, 200, StatusError(200).Code())
}

func TestStatusError_Error(t *testing.T) {
	assert.Equal(t, "0 ", StatusError(0).Error())
	assert.Equal(t, "200 OK", StatusError(200).Error())
	assert.Equal(t, "999 ", StatusError(999).Error())
}

type mockNext struct {
	app    AppVars
	w      http.ResponseWriter
	r      *http.Request
	Err    error
	Called int
}

func (m *mockNext) GetHandler() Handler {
	return func(app AppVars, w http.ResponseWriter, r *http.Request) error {
		m.app = app
		m.w = w
		m.r = r
		m.Called = m.Called + 1
		return m.Err
	}
}

var logger = logging.New(os.Stdout, "test ")

func Test_wrapHandler_successful_request(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{
		CurrentUserDetails: mockUserDetails,
		FirmDetails:        mockFirmDetails,
	}

	errorTemplate := &mockTemplates{}
	envVars := EnvironmentVars{}
	nextHandlerFunc := wrapHandler(logger, mockClient, errorTemplate, envVars)
	next := mockNext{}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	assert.Nil(t, next.Err)
	assert.Equal(t, w, next.w)
	assert.Equal(t, r, next.r)
	assert.Equal(t, 1, next.Called)
	assert.Equal(t, "test-url", next.app.Path)
	assert.Equal(t, mockClient.CurrentUserDetails, next.app.User)
	assert.Equal(t, mockClient.FirmDetails, next.app.Firm)
	assert.Equal(t, 200, w.Result().StatusCode)
}

func Test_wrapHandler_error_creating_AppVars(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{error: errors.New("some API error")}

	errorTemplate := &mockTemplates{}
	envVars := EnvironmentVars{}
	nextHandlerFunc := wrapHandler(logger, mockClient, errorTemplate, envVars)
	next := mockNext{}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	assert.Equal(t, 0, next.Called)
	assert.Equal(t, 1, errorTemplate.count)
	assert.Equal(t, ErrorVars{Code: 500, Error: "some API error"}, errorTemplate.lastVars)
	assert.Equal(t, 500, w.Result().StatusCode)
}

func Test_wrapHandler_404_error_cannot_render_template(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{
		CurrentUserDetails: mockUserDetails,
		FirmDetails:        mockFirmDetails,
	}

	errorTemplate := &mockTemplates{error: errors.New("some template error")}
	envVars := EnvironmentVars{}
	nextHandlerFunc := wrapHandler(logger, mockClient, errorTemplate, envVars)
	next := mockNext{Err: StatusError(404)}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	assert.Equal(t, 1, next.Called)
	assert.Equal(t, w, next.w)
	assert.Equal(t, r, next.r)
	assert.Equal(t, 1, errorTemplate.count)
	assert.IsType(t, ErrorVars{}, errorTemplate.lastVars)
	assert.Equal(t, 404, errorTemplate.lastVars.(ErrorVars).Code)
	assert.Equal(t, "404 Not Found", errorTemplate.lastVars.(ErrorVars).Error)
	assert.Equal(t, 404, w.Result().StatusCode)
}

func Test_wrapHandler_redirects_if_unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{error: sirius.ErrUnauthorized}

	errorTemplate := &mockTemplates{}
	envVars := EnvironmentVars{SiriusURL: "sirius-url", Prefix: "prefix/"}
	nextHandlerFunc := wrapHandler(logger, mockClient, errorTemplate, envVars)
	next := mockNext{}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	assert.Equal(t, 0, errorTemplate.count)
	assert.Equal(t, 302, w.Result().StatusCode)
	location, err := w.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "sirius-url/auth", location.String())
}

func Test_wrapHandler_follows_local_redirect(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)

	mockClient := mockApiClient{
		CurrentUserDetails: mockUserDetails,
		FirmDetails:        mockFirmDetails,
	}

	errorTemplate := &mockTemplates{}
	envVars := EnvironmentVars{Prefix: "prefix/"}
	nextHandlerFunc := wrapHandler(logger, mockClient, errorTemplate, envVars)
	next := mockNext{Err: Redirect("redirect-to-here")}
	httpHandler := nextHandlerFunc(next.GetHandler())
	httpHandler.ServeHTTP(w, r)

	assert.Equal(t, 0, errorTemplate.count)
	assert.Equal(t, 302, w.Result().StatusCode)
	location, err := w.Result().Location()
	assert.Nil(t, err)
	assert.Equal(t, "prefix/redirect-to-here", location.String())
}
