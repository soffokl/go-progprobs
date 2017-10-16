package handler

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type stat struct {
	Count     uint `json:"num_images"`
	AvgWidth  uint `json:"average_width_px"`
	AvgHeight uint `json:"average_height_px"`
	lock      sync.Mutex
}

var (
	stats stat
)

// AddItem increase counter if generated images recalculate average value of height and width
func (s *stat) AddItem(h, w uint) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.AvgHeight = (s.Count*s.AvgHeight + h) / (s.Count + 1)
	s.AvgWidth = (s.Count*s.AvgWidth + w) / (s.Count + 1)
	s.Count++
}

// JSON returns JSON formated stats about requested images
func (s *stat) JSON() (out []byte, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return json.Marshal(s)
}

// Image returns a black image/png that's of the specified width and height
func Image(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(strings.TrimPrefix(r.RequestURI, "/generate/"), "/")
	if len(args) == 3 {
		height, _ := strconv.Atoi(args[1])
		width, _ := strconv.Atoi(args[2])

		if height <= 0 || width <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		img := image.NewRGBA(image.Rect(0, 0, height, width))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

		switch args[0] {
		case "png":
			if err := png.Encode(w, img); err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "jpg":
			if err := jpeg.Encode(w, img, nil); err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stats.AddItem(uint(height), uint(width))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Stats returns application/json with the number of images generated since the server was started and the average width and height of those images.
func Stats(w http.ResponseWriter, r *http.Request) {
	if out, err := stats.JSON(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
