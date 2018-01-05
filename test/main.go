package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gitchander/logfw"
	"github.com/gitchander/logl"
)

func main() {
	example1()
}

func example1() {

	c := logfw.Config{
		FileName:   "logs/test.log",
		MaxSize:    3 * logfw.Kilobyte,
		MaxBackups: 15,
	}
	w, err := logfw.NewWriter(c)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	logger := logl.New(
		logl.Config{
			Handler: &logl.StreamHandler{
				Output: w,
				Format: logl.JsonFormat(),
				//				Format: &logl.TextFormat{
				//					HasLevel:      true,
				//					Date:          true,
				//					Time:          true,
				//					Microseconds:  true,
				//					ShieldSpecial: true,
				//				},
			},
			Level: logl.LevelDebug,
		},
	)

	r := newRandTime()
	for i := 0; i < 1000; i++ {
		//level := logl.Level(randRange(r, int(logl.LevelCritical), int(logl.LevelDebug)+1))
		level := logl.LevelWarning
		logger.Messagef(level, "id %d, text: %s", i, randText(r))
	}
}

func example2() {

	c := logfw.Config{
		FileName:   "logs/test.log",
		MaxSize:    3 * logfw.Kilobyte,
		MaxBackups: 14,
	}
	w, err := logfw.NewWriter(c)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("+")
		w.Close()
	}()

	logger := log.New(w, "master ", log.Ldate|log.Ltime)

	logger.Println("message ok")
	logger.Fatal("message bad")
}

func newRandTime() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randRange(r *rand.Rand, min, max int) int {
	return min + r.Intn(max-min)
}

func randText(r *rand.Rand) string {
	special := []byte(".,;:-!@#$%^&*")
	n := randRange(r, 10, 50)
	bs := make([]byte, n)
	for i := range bs {
		k := r.Intn(100)
		if k < 60 {
			bs[i] = byte(randRange(r, int('a'), int('z')+1))
		} else if k < 80 {
			bs[i] = ' '
		} else {
			bs[i] = special[r.Intn(len(special))]
		}
	}
	return string(bs)
}
