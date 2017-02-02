package sift

import (
	"testing"
	"fmt"
	"time"
)

func TestSiftMatch(t *testing.T) {
	c := NewSiftClient(nil)
	var b bool
	timeCostWrapper("test1",func(){
		b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
		fmt.Printf("test1 %v\n", b)
	})
	timeCostWrapper("test2",func(){
		b = c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test2 %v\n", b)
	})
	timeCostWrapper("test3",func(){
		b = c.Match(&ImageEntity{"img/jianpan_a.jpeg", "003"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test3 %v\n", b)
	})
	timeCostWrapper("test4",func(){
		b = c.Match(&ImageEntity{"img/shubiao_a.jpeg", "005"}, &ImageEntity{"img/jianpan_b.jpeg", "004"})
		fmt.Printf("test4 %v\n", b)
	})

}

func TestEvicted(t *testing.T) {
	c :=NewSiftClient(&SiftClientOption{
		CacheSize: 200,
	})
	var b bool
	timeCostWrapper("test1",func(){
		for i:=0;i<10000;i++{
			b = c.Match(&ImageEntity{"img/beaver.png", 100+i}, &ImageEntity{"img/beaver_xform.png", 200+i})
		}
		//b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
	})
	timeCostWrapper("test2",func(){
		b = c.Match(&ImageEntity{"img/beaver.png", "001"}, &ImageEntity{"img/beaver_xform.png", "002"})
	})
	for {
		time.Sleep(1000*time.Millisecond)
	}

}
func timeCostWrapper(s string, foo func()) {
	now := time.Now()
	foo()
	delta := time.Now().Sub(now).Nanoseconds()
	fmt.Printf("%s %dms\n",s,delta/1e6)
}