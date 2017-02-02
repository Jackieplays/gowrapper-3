package main

import (
	"github.com/JodeZer/opensift_gowrapper"
)

func main() {
	sift.Match(&sift.ImageEntity{"img/test.jpeg", "003"}, &sift.ImageEntity{"img/beaver_xform.png", "002"})
	sift.Match(&sift.ImageEntity{"../img/test.jpeg", "004"}, &sift.ImageEntity{"img/beaver_xform.png", "002"})
}