package mockhttp

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMockHttpServer(t *testing.T) {
	srv := MakeMockHttpServer()
	defer srv.Close()

	resp, err := http.Get(srv.GetBaseUri())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	got, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(got))

	srv.AssignHandler(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oh no!"))
	})
	resp, err = http.Get(srv.GetBaseUri())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	got, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "oh no!", string(got))

	srv.AssignDefaultHandler()
	resp, err = http.Get(srv.GetBaseUri())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	got, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(got))
}
