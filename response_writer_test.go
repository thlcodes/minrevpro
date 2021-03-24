package minrevpro_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/thlcodes/minrevpro"
)

func TestInformativeResponseWriter(t *testing.T) {
	wantStatusCode := 404
	wantStatusCodeString := strconv.Itoa(wantStatusCode)
	wantStatusText := http.StatusText(wantStatusCode)
	wantBody := []byte("testbody")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(wantStatusCode)
		_, err := w.Write(wantBody)
		if err != nil {
			t.Errorf("Failed to write body: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wc := minrevpro.NewInformativeResponseWriter(w)
		assertNotNil(t, wc, "NewInformativeResponseWriter returned nil").must()
		handler(wc, r)
		gotStatusCode := wc.StatusCode()
		assertEqualInt(t, wantStatusCode, gotStatusCode, "StatusCode mismatch: want %d != got %d", wantStatusCode, gotStatusCode)
		gotStatusCodeString := wc.StatusCodeString()
		assertEqualString(t, wantStatusCodeString, gotStatusCodeString, "StatusCodeString mismatch: want %d != got %s", wantStatusCode, gotStatusCodeString)
		gotStatusText := wc.StatusText()
		assertEqualString(t, wantStatusText, gotStatusText, "StatusText mismatch: want %s != got %s", wantStatusText, gotStatusText)
		gotBytesWritten := wc.BytesWritten()
		assertEqualInt(t, len(wantBody), int(gotBytesWritten), "BytesWritten mismatch: want %d != got %d", len(wantBody), gotBytesWritten)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	assertNil(t, err, "Request failed: %v", err).must()
	assertEqualInt(t, wantStatusCode, resp.StatusCode, "Response.StatusCode mistmatch: want %d != got %d", wantStatusCode, resp.StatusCode)
	gotBody, err := ioutil.ReadAll(resp.Body)
	assertNil(t, err, "Reading body failed: %v", err).must()
	resp.Body.Close()
	assertTrue(t, bytes.Equal(wantBody, gotBody), "Body mismatch: want %s != got %s", wantBody, gotBody)
}
