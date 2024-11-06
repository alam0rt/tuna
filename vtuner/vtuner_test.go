package vtuner_test

import (
	"fmt"
	"testing"

	"github.com/alam0rt/tuna/vtuner"
)

func TestFoo(t *testing.T) {
	page := vtuner.NewPage([]vtuner.Item{
		&vtuner.Display{
			Display: "foo",
		},
		&vtuner.Previous{
			Url: "http://localhost:8080/previous",
		},
	}, false)

	_, err := page.Render()
	if err != nil {
		t.Fatal(err)
	}

	out, err := page.Render()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("output: %s", out)
	fmt.Print(string(out))
	t.Fatal(string(out))
}
