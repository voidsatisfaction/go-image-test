package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/xor-gate/goexif2/exif"
)

func main() {
	// Check files are well oriented
	// searchDir := "./test-img"
	// err := filepath.Walk(searchDir, func(p string, f os.FileInfo, err error) error {
	// 	if f.Name() == searchDir[2:] {
	// 		return nil
	// 	}
	// 	fmt.Printf("file: %+v\n\n", f.Name())
	// 	fileName := f.Name()
	// 	filePath := path.Join(searchDir, fileName)
	// 	file, err := os.Open(filePath)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	x, err := exif.Decode(file)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	orientation, _ := x.Get(exif.Orientation)
	// 	fmt.Printf("fileOrientation: %+v\n\n", orientation)
	// 	return nil
	// })

	// Get Orientation
	fileNames := []string{
		"up.jpg", "up-mirrored.jpg", "down.jpg", "down-mirrored.jpg",
		"left-mirrored.jpg", "left.jpg", "right-mirrored.jpg", "right.jpg",
	}
	for _, fileName := range fileNames {
		filePath := path.Join("./test-img", fileName)
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		orientation, err := checkFileOrientation(f)
		if err != nil {
			fmt.Errorf("Failed Orientation get")
		}

		// Fix orientation
		fixOrientation(f, orientation)

		// Delete Exif

		// Save
	}
}

func checkFileOrientation(f *os.File) (int, error) {
	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return -1, err
	}

	orientation, err := tag.Int(0)
	if err != nil {
		return -1, err
	}
	return orientation, nil
}

func fixOrientation(f *os.File, orientation int) {
	transFunctions := map[int]func(image.Image) *image.NRGBA{
		1: imaging.Clone,
		2: imaging.FlipH,
		3: imaging.Rotate180,
		4: imaging.FlipV,
		5: imaging.Transpose,
		6: imaging.Rotate270,
		7: imaging.Transverse,
		8: imaging.Rotate90,
	}
	if orientation == -1 {
		return
	}
	f.Seek(0, 0)
	srcImage, err := imaging.Open(f.Name())
	if err != nil {
		fmt.Errorf("image open error!")
	}

	dstImage := transFunctions[orientation](srcImage)
	if err != nil {
		fmt.Errorf("Cannot make output file")
	}
	dst := fmt.Sprintf("test-img/example_%d.jpg", orientation)
	fmt.Println(dst)
	err = imaging.Save(dstImage, dst)
	if err != nil {
		fmt.Errorf("file save Error occured!")
	}
}

// // アップロード前処理:Exif除去,Orientation調整
// func (m *UploadImage) Format() error {
// 	// check EXIF
// 	e, err := exif.Decode(m.File)
// 	if err != nil {
// 		// EXIF情報がないときは抜ける
// 		return nil
// 	}
// 	tag, err := e.Get(exif.Orientation)
// 	if err != nil {
// 		// Orientationがないときは抜ける
// 		return nil
// 	}
// 	orientation, err := tag.Int(0)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert orientation: %v", err)
// 	}
//
// 	// rotate
// 	m.File.Seek(0, 0)
// 	srcImage, _, err := image.Decode(m.File)
// 	if err != nil {
// 		return fmt.Errorf("failed to decode srcImage: %v", err)
// 	}
// 	rect := srcImage.Bounds()
//
// 	if orientation >= 5 && orientation <= 8 {
// 		rect = image.Rect(0, 0, rect.Size().Y, rect.Size().X)
// 	}

// 	dstImage := image.NewNRGBA(rect)
// 	affine, ok := affines[orientation]
// 	if !ok {
// 		affine = graphics.I // Assign default affine func
// 	}

// 	err = affine.TransformCenter(dstImage, srcImage, interp.Bilinear)
// 	if err != nil {
// 		return fmt.Errorf("failed to affine map: %v", err)
// 	}
// 	dstFilename := m.File.Name() + "_out"
// 	dstFile, err := os.Create(dstFilename)
// 	if err != nil {
// 		return fmt.Errorf("failed to open dstFile: %v", err)
// 	}
// 	err = jpeg.Encode(dstFile, dstImage, &jpeg.Options{Quality: 100})
// 	if err != nil {
// 		return fmt.Errorf("failed to encode dstImage: %v", err)
// 	}

// 	m.File = dstFile
// 	m.File.Seek(0, 0)
// 	err = m.CheckDimension() // 画像の縦横の再計算
// 	if err != nil {
// 		return fmt.Errorf("failed to recheck dimension: %v", err)
// 	}
// 	return nil
// }
