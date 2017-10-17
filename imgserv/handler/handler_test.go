package handler

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestImage_ServeHTTP(t *testing.T) {
	cases := []struct {
		uri  string
		code int
	}{
		{"/generate/png/100/10", http.StatusOK},
		{"/generate/jpg/10/100", http.StatusOK},
		{"/generate/jjj/100/10", http.StatusBadRequest},
		{"/generate/png/99999/99999", http.StatusBadRequest},
		{"/generate/jpg/99999/99999", http.StatusBadRequest},
	}

	for _, test := range cases {
		rr := httptest.NewRecorder()
		h := http.Handler(&Image{})

		h.ServeHTTP(rr, httptest.NewRequest("GET", test.uri, nil))
		if status := rr.Code; status != test.code {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.code)
		}

		if test.code == http.StatusOK {
			if args := strings.Split(test.uri, "/"); len(args) == 5 {
				var err error
				var img image.Image

				switch args[2] {
				case "png":
					img, err = png.Decode(rr.Body)
				case "jpg":
					img, err = jpeg.Decode(rr.Body)
				default:
					t.Errorf("handler received wrong request format for http.StatusOK: %v", test.uri)
				}

				if err != nil {
					t.Errorf("handler returned not an image file: %v for %v", err.Error(), test.uri)
				} else {
					width, _ := strconv.Atoi(args[3])
					height, _ := strconv.Atoi(args[4])

					if img.Bounds().Max.X != width || img.Bounds().Max.Y != height {
						t.Errorf("handler returned wrong sized image file: %v ", test.uri)
					}

					for x := 0; x < img.Bounds().Max.X; x++ {
						for y := 0; y < img.Bounds().Max.Y; y++ {
							r, g, b, a := img.At(x, y).RGBA()
							if !(r+g+b == 0 && a == 65535) {
								t.Errorf("handler returned wrong colored image: %v ", test.uri)
							}
						}
					}
				}

			} else {
				t.Errorf("")
			}
		}
	}
}

func TestStats(t *testing.T) {
	rr := httptest.NewRecorder()
	h := http.HandlerFunc(Stats)
	h.ServeHTTP(rr, httptest.NewRequest("GET", "/stats", nil))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {

		if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
			t.Errorf("content type header does not match: got %v want %v", ctype, "application/json")
		}

		var s stat
		if err := json.Unmarshal(rr.Body.Bytes(), &s); err != nil {
			t.Errorf("handler returned not a JSON document: %v", err.Error())
		}
	}
}
