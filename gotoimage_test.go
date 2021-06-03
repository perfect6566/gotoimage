package gotoimage

import (
	"bytes"
	"image/png"
	"math"
	"os"
	"testing"
)

func TestRender(t *testing.T) {


	got, err := Render("data/example")
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	imgGot, err := png.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatalf("cannot read rendered image: %v", err)
	}

	want, err := os.ReadFile("data/want.png")
	if err != nil {
		t.Fatalf("cannot read gold test file: %v", err)
	}

	imgWant, err := png.Decode(bytes.NewReader(want))
	if err != nil {
		t.Fatalf("cannot read gold image: %v", err)
	}

	if math.Abs(float64(imgGot.Bounds().Dx()-imgWant.Bounds().Dx())) > 5 ||
		math.Abs(float64(imgGot.Bounds().Dy()-imgWant.Bounds().Dy())) > 5 {

		err := os.WriteFile("testdata/got.png", got, os.ModePerm)
		if err != nil {
			t.Errorf("failed to write image: %v", err)
		}

		t.Fatalf("image size does not match: got %+v, want %+v", imgGot.Bounds(), imgWant.Bounds())
	}
}

