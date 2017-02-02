package sift

/*
#cgo CFLAGS : -Iinclude -I/usr/local/include/opencv
#cgo LDFLAGS: -Llib/ubuntu -L/usr/local/lib  -lopensift -lopencv_calib3d -lopencv_contrib -lopencv_core -lopencv_features2d -lopencv_flann -lopencv_gpu -lopencv_highgui -lopencv_imgproc -lopencv_legacy -lopencv_ml -lopencv_nonfree -lopencv_objdetect -lopencv_ocl -lopencv_photo -lopencv_stitching -lopencv_superres -lopencv_ts -lopencv_video -lopencv_videostab -ltbb -lXext -lX11 -lICE -lSM -lGL -lGLU -lrt -lpthread -lm -ldl -lgtk-x11-2.0 -lgdk-x11-2.0 -latk-1.0 -lgio-2.0 -lpangoft2-1.0 -lpangocairo-1.0 -lgdk_pixbuf-2.0 -lcairo -lpango-1.0 -lfontconfig -lgobject-2.0 -lglib-2.0 -lfreetype
#include "sift.h"
#include "imgfeatures.h"
#include "kdtree.h"
#include "utils.h"
#include "xform.h"

#include <cv.h>
#include <cxcore.h>
#include <highgui.h>

#include <stdio.h>

#define KDTREE_BBF_MAX_NN_CHKS 200

#define NN_SQ_DIST_RATIO_THR 0.49

int compare_features(struct feature* f0, int count0,struct feature* f1,int count1){
	struct kd_node* kd_root;
	double d0, d1;
	struct feature** nbrs;
	int num_matches = 0;
	size_t i;
	kd_root = kdtree_build(f1, count1);

	for(i = 0; i < count0; ++i) {
		int k;
		struct feature* feat;
		feat = f0 + i;
		k = kdtree_bbf_knn( kd_root, feat, 2, &nbrs, KDTREE_BBF_MAX_NN_CHKS );
		if( k == 2 ) {
			d0 = descr_dist_sq( feat, nbrs[0] );
			d1 = descr_dist_sq( feat, nbrs[1] );
		if( d0 < d1 * NN_SQ_DIST_RATIO_THR ) ++num_matches;
		}
		free( nbrs );
	}
	kdtree_release( kd_root );
	return num_matches;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
	"sync"
)

/*
the impl
*/
func MatchImpl(path1, path2 string) (int, int, int) {
	var features, features2 *C.struct_feature
	img := C.cvLoadImage(C.CString(path1), C.int(1))
	img2 := C.cvLoadImage(C.CString(path2), C.int(1))
	n := C.sift_features(img, &features)
	n2 := C.sift_features(img2, &features2)
	fmt.Printf("nums %d\nnums %d\n", n, n2)
	res := C.compare_features(features, n, features2, n2)
	fmt.Printf("res %d\n", res)
	// release memory
	C.free(unsafe.Pointer(features))
	C.free(unsafe.Pointer(features2))
	C.cvReleaseImage(&img)
	C.cvReleaseImage(&img2)
	return int(n), int(n2), int(res)

}

func Match(ie1, ie2 *ImageEntity) bool {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("panic recover %s\n", e)
		}
	}()

	return DEFAULT_SIFT_CLIENT.Match(ie1, ie2)
}

const (
	DEFAULT_KDTREE_BBF_MAX_NN_CHKS = 200
	DEFAULT_NN_SQ_DIST_RATIO_THR = 0.49
)

//pai nao dai
var DEFAULT_THRESHOLDFUNC thresholdFunc = func(f1 int, f2 int, m int) bool {
	return float32(m) / float32(f2) >= float32(0.2) || float32(m) / float32(f1) >= float32(0.2) || m > 80
}

type thresholdFunc func(int, int, int) bool



/*
feature_t
C struct feature wrapper
*/
type feature_t struct {
	Feat *C.struct_feature
}
type img_t struct {
	Img *C.struct__IplImage //don't know why
}

/**
C functions wrapper
 */
func cvLoadImg(path string) *img_t {
	return &img_t{
		Img: C.cvLoadImage(C.CString(path), C.int(1)),
	}
}

func freeImg(img *img_t) {
	C.cvReleaseImage(&(img.Img))
}

func sift_features(img *img_t) *featureData {
	var features *C.struct_feature
	n := C.sift_features(img.Img, &features)
	return &featureData{
		Feature:&feature_t{
			Feat:features,
		},
		Count:int(n),
	}
}

func free_featureData(featData *featureData) {
	C.free(unsafe.Pointer(featData.Feature.Feat))
}

func loadFeatureData(ie *ImageEntity) *featureData {
	img := cvLoadImg(ie.FilePath)
	if img.Img == nil {
		panic(fmt.Sprintf("can't load image path %s", ie.FilePath))
	}
	featData := sift_features(img)
	featData.ImageEntity = *ie
	freeImg(img)
	return featData
}

func compareFeatureData(feature1, feature2 *featureData) int {
	return int(C.compare_features(feature1.Feature.Feat, C.int(feature1.Count), feature2.Feature.Feat, C.int(feature2.Count)))
}
/*
imgSignature
 */
type imgSignature struct {
	FilePath string
	Md5      string
}

func (is *imgSignature)Match(img *imgSignature) bool {
	return false
}

/**
featureData
 */
type featureData struct {
	Feature *feature_t
	Count   int
	ImageEntity
}

/**

 */
type ImageEntity struct {
	FilePath string
	UniqueID interface{} //used for cache
}

/**
-------interface----------
 */
type Client interface {
	Load(*ImageEntity) error
	Match(*ImageEntity, *ImageEntity) bool
	/*cache interface*/
}

/*
siftClient
*/
type siftClient struct {
	kdtree_bbf_max_nn_chks int
	nn_sq_dist_ratio_thr   float32
	tf                     thresholdFunc
	lru                    Lru //interface
	lruMutex               *sync.RWMutex
}

func (c *siftClient)Load(ie *ImageEntity) error {
	return nil
}

func (c *siftClient)Match(ie1 *ImageEntity, ie2 *ImageEntity) bool {

	var features1, features2 *featureData
	var ok1, ok2 bool
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	features1, ok1 = c.getCache(ie1)
	features2, ok2 = c.getCache(ie2)
	if !ok1 {
		wg.Add(1)
		go func(){
			features1 = loadFeatureData(ie1)
			c.setCache(ie1, features1)
			wg.Done()
		}() // this a time cost operation,could be



	}

	if !ok2 {
		wg.Add(1)
		go func(){
			features2 = loadFeatureData(ie2)
			c.setCache(ie2, features2)
			wg.Done()
		}()

	}

	wg.Wait()
	res := compareFeatureData(features1, features2)
	if c.tf != nil {
		return c.tf(features1.Count, features2.Count, res)
	}else {
		return DEFAULT_THRESHOLDFUNC(features1.Count, features2.Count, res)
	}
}

func (c *siftClient)getCache(ie *ImageEntity) (*featureData, bool) {
	c.lruMutex.RLock()
	defer c.lruMutex.RUnlock()
	if feature, ok := c.lru.Get(ie.UniqueID); ok {
		return feature.(*featureData), ok
	}
	return nil, false
}

func (c *siftClient)setCache(ie *ImageEntity, feature *featureData) {
	c.lruMutex.Lock()
	c.lru.Add(ie.UniqueID, feature)
	c.lruMutex.Unlock()
}

/*
SiftClientOption
*/
var DEFAULT_SIFT_OPTION SiftClientOption = SiftClientOption{
	Kdtree_bbf_max_nn_chks: DEFAULT_KDTREE_BBF_MAX_NN_CHKS,
	Nn_sq_dist_ratio_thr:DEFAULT_NN_SQ_DIST_RATIO_THR,
	Tf:DEFAULT_THRESHOLDFUNC,
	CacheSize: 1024, //1024
}

var DEFAULT_SIFT_CLIENT Client = NewSiftClient(&DEFAULT_SIFT_OPTION)

var DEFAULT_ON_EVICTED func(interface{},interface{}) = func (key,value interface{}) {
	fmt.Println("on evicted")
	featData := value.(*featureData)
	free_featureData(featData)
}

type SiftClientOption struct {
	Kdtree_bbf_max_nn_chks int
	Nn_sq_dist_ratio_thr   float32
	Tf                     thresholdFunc
	CacheSize              int //TODO
}


/*

 */

func NewSiftClient(option *SiftClientOption) Client {
	if option == nil {
		option = &DEFAULT_SIFT_OPTION
	}
	var c *siftClient = &siftClient{}
	if option.Kdtree_bbf_max_nn_chks != 0 {
		c.kdtree_bbf_max_nn_chks = option.Kdtree_bbf_max_nn_chks

	}else {
		c.kdtree_bbf_max_nn_chks = DEFAULT_KDTREE_BBF_MAX_NN_CHKS
	}

	if option.Nn_sq_dist_ratio_thr <= 1e-7 {
		c.nn_sq_dist_ratio_thr = DEFAULT_NN_SQ_DIST_RATIO_THR
	}else {
		c.nn_sq_dist_ratio_thr = option.Nn_sq_dist_ratio_thr
	}

	c.lru = NewClassicLru(option.CacheSize)

	c.lruMutex = new(sync.RWMutex)

	c.lru.SetOnEvicted(DEFAULT_ON_EVICTED)

	c.tf = option.Tf

	return c
}
