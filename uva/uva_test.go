package uva_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/uva-cli/uva"
)

// three-problem fixture — indices match the uHunt wire format
// [0]=pid [1]=num [2]=title [3]=dacu ... [7]=sube ... [14]=tle ... [16]=wa ... [18]=ac [19]=rtl [20]=status
var fixtureProblems = `[
[36,100,"The 3n+1 problem",108913,0,0,0,6949,0,0,139127,0,100069,387,78813,5209,361456,6555,255926,3000,1,0],
[37,101,"The Blocks Problem",18108,0,0,0,933,0,0,16197,0,25728,24,13080,200,29446,8539,33408,3000,1,0],
[38,102,"Ecological Bin Packing",7654,0,0,0,500,0,0,5000,0,2000,10,1000,100,5000,200,10000,2000,1,0]
]`

var fixtureProblem100 = `[36,100,"The 3n+1 problem",108913,0,0,0,6949,0,0,139127,0,100069,387,78813,5209,361456,6555,255926,3000,1,0]`

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/p" {
			_, _ = w.Write([]byte(fixtureProblems))
			return
		}
		if r.URL.Path == "/api/p/num/100" {
			_, _ = w.Write([]byte(fixtureProblem100))
			return
		}
		http.NotFound(w, r)
	}))
}

func newClient(srv *httptest.Server) *uva.Client {
	cfg := uva.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0
	return uva.NewClient(cfg)
}

func TestList(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	c := newClient(srv)
	problems, err := c.List(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 3 {
		t.Fatalf("got %d problems, want 3", len(problems))
	}
	if problems[0].Num != 100 {
		t.Errorf("first Num = %d, want 100", problems[0].Num)
	}
	if problems[0].Title != "The 3n+1 problem" {
		t.Errorf("first Title = %q, want %q", problems[0].Title, "The 3n+1 problem")
	}
}

func TestListLimit(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	c := newClient(srv)
	problems, err := c.List(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 2 {
		t.Fatalf("got %d problems with limit=2, want 2", len(problems))
	}
}

func TestSearch(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	c := newClient(srv)
	problems, err := c.Search(context.Background(), "blocks", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 {
		t.Fatalf("got %d results for 'blocks', want 1", len(problems))
	}
	if problems[0].Num != 101 {
		t.Errorf("search result Num = %d, want 101", problems[0].Num)
	}
}

func TestGetByNum(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	c := newClient(srv)
	p, err := c.GetByNum(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("GetByNum returned nil, want a Problem")
	}
	if p.Num != 100 {
		t.Errorf("Num = %d, want 100", p.Num)
	}
}
