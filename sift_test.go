package sift

import (
	"testing"
	"fmt"
	"time"
)

func TestSiftMatch(t *testing.T) {
	c := NewSiftClient(nil)
	var b bool
	timeCostWrapper("test1", func() {
		b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
		fmt.Printf("test1 %v\n", b)
	})
	timeCostWrapper("test2", func() {
		b = c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test2 %v\n", b)
	})
	timeCostWrapper("test3", func() {
		b = c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test3 %v\n", b)
	})
	timeCostWrapper("test4", func() {
		b = c.Match(&ImageEntity{"img/shubiao_a.jpeg", "005"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test4 %v\n", b)
	})

	timeCostWrapper("macth same file same id", func() {
		c.Match(&ImageEntity{"img/jianpan_a.jpeg", "103"}, &ImageEntity{"img/jianpan_a.jpeg", "103"})
	})

	timeCostWrapper("macth same file diff id", func() {
		c.Match(&ImageEntity{"img/jianpan_a.jpeg", "113"}, &ImageEntity{"img/jianpan_a.jpeg", "123"})
	})


}

func TestEvicted(t *testing.T) {
	c := NewSiftClient(&SiftClientOption{
		CacheSize: 20,
	})
	var b bool
	timeCostWrapper("test1", func() {
		for i := 0; i < 10000; i++ {
			b = c.Match(&ImageEntity{"img/beaver.png", 100 + i}, &ImageEntity{"img/beaver_xform.png", 200 + i})
		}
		//b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
	})
	timeCostWrapper("test2", func() {
		b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
	})
	for {
		time.Sleep(1000 * time.Millisecond)
	}

}

func TestLoad(t *testing.T) {
	c := NewSiftClient(&SiftClientOption{
		CacheSize: 200,
	})
	timeCostWrapper("img/beaver.png laod", func() {
		c.Load(&ImageEntity{"img/beaver.png", "001"})
	})
	timeCostWrapper("img/beaver_xform.png load", func() {
		c.Load(&ImageEntity{"img/beaver_xform.png", "002"})
	})
	timeCostWrapper("panic load", func() {
		c.Load(&ImageEntity{"beaver_xform.png", "002"})
	})
	timeCostWrapper("img/jianpan_a.jpeg load", func() {
		c.Load(&ImageEntity{"img/jianpan_a.jpeg", "003"})
	})
	timeCostWrapper("cache macth", func() {
		c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
	})

	timeCostWrapper("half cache macth", func() {
		c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
	})

	timeCostWrapper("full cache macth", func() {
		c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
	})

	timeCostWrapper("null cache macth", func() {
		c.Match(&ImageEntity{"img/jianpan_a.jpeg", "103"}, &ImageEntity{"img/jianpan_b.jpeg", "104"})
	})


}

func TestGetlplImageSize(t *testing.T) {
	timeCostWrapper("beaver", func() {
		img := cvLoadImg("img/beaver.png")
		fmt.Printf("img/beaver.png %d\n", getImageSize(img))
	})

	timeCostWrapper("beaver_xform", func() {
		img := cvLoadImg("img/beaver_xform.png")
		fmt.Printf("img/beaver_xform.png %d\n", getImageSize(img))
	})

	timeCostWrapper("jianpan_a", func() {
		img := cvLoadImg("img/jianpan_a.jpeg")
		fmt.Printf("img/jianpan_a.jpeg %d\n", getImageSize(img))
	})
}

func timeCostWrapper(s string, foo func()) {
	now := time.Now()
	foo()
	delta := time.Now().Sub(now).Nanoseconds()
	fmt.Printf("%s %dms\n", s, delta / 1e6)
}