package thttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var tests = []struct {
	in  string
	out string
}{
	{`"Hello"`, `"Hello"`},
}

type TestProcessor struct{}

func (t TestProcessor) Process(in, out thrift.TProtocol) (bool, thrift.TException) {
	str, err := in.ReadString()
	if err != nil {
		fmt.Println("ReadError:", err)
		return false, nil
	}

	if err = out.WriteString(str); err != nil {
		fmt.Println("WriteError", err)
		return false, nil
	}

	if err = out.Flush(); err != nil {
		fmt.Println("FlushError", err)
		return false, nil
	}

	return true, nil
}

func TestTHandler(t *testing.T) {
	p := TestProcessor{}
	h := NewThriftHandler(p)
	ts := httptest.NewServer(http.HandlerFunc(h))
	defer ts.Close()

	for _, tt := range tests {
		r := strings.NewReader(tt.in)
		res, err := http.Post(ts.URL, "application/json", r)
		if err != nil {
			t.Errorf("Couldn't Post: %q", err)
			continue
		}

		if res.StatusCode != 200 {
			t.Errorf("Status code => %d, want 200", res.StatusCode)
			continue
		}

		ret, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Coudn't read body: %q", err)
			continue
		}

		s := string(ret)
		if tt.out != s {
			t.Errorf("Post %q => body: %q, want %q", tt.in, s, tt.out)
		}
	}
}
