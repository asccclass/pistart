package main

import (
   "os"
   "fmt"
   "image"
   "strings"
   "image/color"
   "image/jpeg"
   "image/png"
   _ "image/gif"
   _ "image/jpeg"
   _ "image/png"
   "golang.org/x/image/bmp"
)

// ConvertToRedBlackWhite converts image to red, black, white only
func ConvertToRedBlackWhite(img image.Image)(*image.RGBA) {
   bounds := img.Bounds()
   dst := image.NewRGBA(bounds)
   for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
      for x := bounds.Min.X; x < bounds.Max.X; x++ {
         r, g, b, _ := img.At(x, y).RGBA()
         // normalize to 0-255
         r8 := uint8(r >> 8)
         g8 := uint8(g >> 8)
         b8 := uint8(b >> 8)
         // classify color
         c := classifyColor(r8, g8, b8)
         dst.Set(x, y, c)
      }
   }
   return dst
}

// classifyColor returns color.Red, color.Black, or color.White
func classifyColor(r, g, b uint8) color.Color {
   // Red: R dominant and G,B low
   if r > 150 && g < 100 && b < 100 {
      return color.RGBA{255, 0, 0, 255}
   }
   // Black: all low
   if r < 80 && g < 80 && b < 80 {
      return color.Black
   }
   // Otherwise: white
   return color.White
}

// ResizeImage resizes an image to targetWidth x targetHeight using nearest neighbor
func ResizeImage(img image.Image, width, height int) *image.RGBA {
   dst := image.NewRGBA(image.Rect(0, 0, width, height))
   srcBounds := img.Bounds()
   for y := 0; y < height; y++ {
      for x := 0; x < width; x++ {
         sx := (x * srcBounds.Dx()) / width
         sy := (y * srcBounds.Dy()) / height
         dst.Set(x, y, img.At(sx+srcBounds.Min.X, sy+srcBounds.Min.Y))
      }
   }
   return dst
}

// IsBMP checks if a file is BMP format
func IsBMP(filename string) (bool) {
   f, err := os.Open(filename)
   if err != nil {
      fmt.Println(err.Error())
      return false
   }
   defer f.Close()
   buf := make([]byte, 2)
   _, err = f.Read(buf)
   if err != nil {
      fmt.Println(err.Error())
      return false
   }
   return buf[0] == 'B' && buf[1] == 'M'
}

// LoadImage loads and decodes an image file (supports bmp, jpeg, png, gif)
func LoadImage(filename string) (image.Image, error) {
   if filename == "" {
      return nil, fmt.Errorf("No input image")
   }
   f, err := os.Open(filename)
   if err != nil {
      return nil, err
   }
   defer f.Close()

   var img image.Image
   if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
      img, err = jpeg.Decode(f)
   } else if strings.HasSuffix(filename, ".png") {
      img, err = png.Decode(f)
   } else if IsBMP(filename) {
      img, _, err = image.Decode(f)
      return nil, err
   } else {
      return img, fmt.Errorf("不支援的輸入格式")
   }
   return img, err
}

func toWaveBmp(input, output string) {
   img, err := LoadImage(input)
   if err != nil {
      fmt.Println("Error loading image:", err.Error())
      return
   }
   img = ConvertToRedBlackWhite(img)                 // Convert to red/black/white
   img = ResizeImage(img, 800, 528) // Resize
   outFile, err := os.Create(output)                 // Save result outFile)io.Writer
   if err != nil {
      fmt.Println("Error creating output:", err)
      return
   }
   defer outFile.Close()

   // Output as BMP
   if err := bmp.Encode(outFile, img); err != nil {
      fmt.Println("Error saving BMP:", err.Error())
      return
   }
}
