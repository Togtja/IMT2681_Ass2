package tests

import (
	"RESTGvkGitLab/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkRequest(t *testing.T, method string, URL string, statusCheck int) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := api.SetupHandlers()
	//handler2 := http.HandleFunc(api.StatusHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Result().StatusCode
	if status != statusCheck {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	return rr
}

func onlyGetRequest(t *testing.T, URL string) {
	checkRequest(t, "GET", URL, http.StatusOK)

	checkRequest(t, "POST", URL, http.StatusNotImplemented)
	checkRequest(t, "DELETE", URL, http.StatusNotImplemented)
	unusedRequest(t, URL)
}
func onlyPostGetRequest(t *testing.T, URL string) {
	checkRequest(t, "GET", URL, http.StatusOK)
	checkRequest(t, "POST", URL, http.StatusOK)
	checkRequest(t, "DELETE", URL, http.StatusNotImplemented)
	unusedRequest(t, URL)
}
func onlyDeletePostGetRequest(t *testing.T, URL string) {
	checkRequest(t, "GET", URL, http.StatusOK)
	checkRequest(t, "POST", URL, http.StatusOK)
	checkRequest(t, "DELETE", URL, http.StatusOK)
	unusedRequest(t, URL)
}

//All unused request in this API
func unusedRequest(t *testing.T, URL string) {
	checkRequest(t, "PUT", URL, http.StatusNotImplemented)
	checkRequest(t, "PATCH", URL, http.StatusNotImplemented)
	checkRequest(t, "COPY", URL, http.StatusNotImplemented)
	checkRequest(t, "OPTIONS", URL, http.StatusNotImplemented)
	checkRequest(t, "LINK", URL, http.StatusNotImplemented)
	checkRequest(t, "UNLINK", URL, http.StatusNotImplemented)
	checkRequest(t, "PURGE", URL, http.StatusNotImplemented)
	checkRequest(t, "LOCK", URL, http.StatusNotImplemented)
	checkRequest(t, "UNLOCK", URL, http.StatusNotImplemented)
	checkRequest(t, "PROPFIND", URL, http.StatusNotImplemented)
	checkRequest(t, "VIEW", URL, http.StatusNotImplemented)

}
