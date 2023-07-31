package image

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var paths = []string{"yellow.jpg", "fog.jpg", "lazers.jpg", "manhattan.jpg"}
var dimx = []int{4080, 4080, 4080, 4080}
var dimy = []int{3072, 3072, 3072, 3072}

func TestOpen(t *testing.T) {

	for n, path := range paths {
		img, err := Open(fmt.Sprintf("%v/%v", "../samples/", path))
		if err != nil {
			t.FailNow()
		}

		bounds := img.Bounds()

		if bounds.Dx() != dimx[n] {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n], bounds.Dx())
		}

		if bounds.Dy() != dimy[n] {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n], bounds.Dy())
		}
	}
}

func __TestOpenCr2(t *testing.T) {

	path := "../samples/sample.CR2"
	img, err := OpenCr2(path)
	if err != nil {
		t.Fail()
	}

	bounds := img.Bounds()

	if bounds.Dx() != 5184 {
		t.Fatalf("%v expected, %v found", 5184, bounds.Dx())
	}

	if bounds.Dy() != 3456 {
		t.Fatalf("%v expected, %v found", 3456, bounds.Dy())
	}
}

func __TestConvertCr2ToJpegAndSave(t *testing.T) {

	path := "../samples/sample.CR2"

	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fail()
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "CR2") {
			jpeg := fmt.Sprintf("%v%v", strings.TrimSuffix(entry.Name(), ".CR2"), ".JPG")
			t.Logf("This file will be convered %v to %v", entry.Name(), jpeg)

			fqp := fmt.Sprintf("%v/%v", path, entry.Name())
			img, err := OpenCr2(fqp)
			if err != nil {
				t.Fail()
			}

			fqp2 := fmt.Sprintf("%v/jpg/%v", path, jpeg)
			err = Save(img, fqp2)
			if err != nil {
				t.Errorf(err.Error())
			}

		}
	}

}

func TestResize(t *testing.T) {

	for n, path := range paths {

		img, err := Open(fmt.Sprintf("%v/%v", "../samples", path))
		if err != nil {
			t.Fail()
		}

		bounds := img.Bounds()

		if bounds.Dx() != dimx[n] {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n], bounds.Dx())
		}

		if bounds.Dy() != dimy[n] {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n], bounds.Dy())
		}

		// resize abs
		img2, err := ResizeAbs(img, bounds.Dx()/2, bounds.Dy()/2)
		if err != nil {
			t.Fail()
		}

		bounds = img2.Bounds()

		if bounds.Dx() != dimx[n]/2 {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n]/2, bounds.Dx())
		}

		if bounds.Dy() != dimy[n]/2 {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n]/2, bounds.Dy())
		}

		img3, err := ResizeScale(img, 0.1)
		if err != nil {
			t.Fail()
		}

		bounds = img3.Bounds()

		if bounds.Dx() != dimx[n]/10 {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n]/10, bounds.Dx())
		}

		if bounds.Dy() != dimy[n]/10 {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n]/10, bounds.Dy())
		}
	}

}

func TestRotate(t *testing.T) {

	for n, path := range paths {

		img, err := Open(fmt.Sprintf("%v/%v", "../samples", path))

		if err != nil {
			t.FailNow()
		}

		img2, err := Rotate(img, 90)
		if err != nil {
			t.FailNow()
		}

		bounds := img2.Bounds()

		if bounds.Dy() != dimx[n] {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n], bounds.Dy())
		}

		if bounds.Dx() != dimy[n] {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n], bounds.Dx())
		}

		img3, err := Rotate(img, 180)
		if err != nil {
			t.FailNow()
		}
		bounds = img3.Bounds()

		if bounds.Dx() != dimx[n] {
			t.Fatalf("Dimx: %v expected: %v, found: %v", path, dimx[n], bounds.Dx())
		}

		if bounds.Dy() != dimy[n] {
			t.Fatalf("Dimy: %v expected: %v, found: %v", path, dimy[n], bounds.Dy())
		}

		if err != nil {
			t.Fail()
		}
	}

}

func TestMetadata(t *testing.T) {

	for _, path := range paths {
		meta, err := Metadata(fmt.Sprintf("%v/%v", "../samples", path))
		if err != nil {
			t.FailNow()
		}
		_ = meta
	}

}
