package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func adjust(c uint32, amount float64) uint8 {
	result := ((float64(c>>8)/255-0.5)*amount + 0.5) * 255

	if result > 255 {
		return 255
	} else if result < 0 {
		return 0
	}

	return uint8(result)
}

func contrast(img image.Image, amount float64) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			o.Set(x, y, color.RGBA{
				adjust(r, amount),
				adjust(g, amount),
				adjust(b, amount),
				uint8(a),
			})
		}
	}

	return o
}

func grayscale(img image.Image) image.Image {
	b := img.Bounds()
        o := image.NewRGBA(b)

        for y := b.Min.Y; y < b.Max.Y; y++ {
                for x := b.Min.X; x < b.Max.X; x++ {
                        r, g, b, a := img.At(x, y).RGBA()
			
			intensity := ((float32(0.2126) * float32(r)) + (float32(0.7152) * float32(g)) + (float32(0.0722) * float32(b))) / 255 
			if intensity > 255 { intensity = 255}
			if intensity < 0 { intensity = 0 }	
			o.Set(x, y, color.RGBA{
                                uint8(intensity),
                                uint8(intensity),
                                uint8(intensity),
                                uint8(a), 
                        })      
                }       
        }       
        
        return o
}

func invert(img image.Image) image.Image {
        b := img.Bounds()
        o := image.NewRGBA(b)

        for y := b.Min.Y; y < b.Max.Y; y++ {
                for x := b.Min.X; x < b.Max.X; x++ {
                        r, g, b, a := img.At(x, y).RGBA()

                        iR := 255 - r
			iG := 255 - g
			iB := 255 - b
                       
                        
                        o.Set(x, y, color.RGBA{
                                uint8(iR),
                                uint8(iG),
                                uint8(iB),
                                uint8(a),
                        })
                }
        }

        return o
}


func clamp(val float32) uint8 {
	if val > 255 { return 255 }
	if val < 0 { return 0 }
	return uint8(val)
}

func motionblur(img image.Image) image.Image {
        b := img.Bounds()
        o := copyimage(img)

        for y := b.Min.Y; y < b.Max.Y; y++ {
                for x := b.Min.X; x < b.Max.X; x++ {
                        r, g, b, _ := o.At(x, y).RGBA()
			
			for j := 1; j < 100; j++ {
				r1, g1, b1, a1 := o.At(x+j, y).RGBA()
                        	iR := float32(0.2 / float32(j)) * float32(r)
				iG := float32(0.2 / float32(j)) * float32(g)
				iB := float32(0.2 / float32(j)) * float32(b)

                        	o.Set(x+j, y, color.RGBA{clamp((float32(r1) + iR) /255),
						   clamp((float32(g1) + iG) /255),
						   clamp((float32(b1) + iB) /255),
						   uint8(a1)})
				}
                        }
        }

        return o
}

func copyimage(img image.Image) *image.RGBA {
	 b := img.Bounds()
        o := image.NewRGBA(b)

        for y := b.Min.Y; y < b.Max.Y; y++ {
                for x := b.Min.X; x < b.Max.X; x++ {
                        r, g, b, a := img.At(x, y).RGBA()
                        
                        o.Set(x, y, color.RGBA{uint8(r),
                                                   uint8(g),
                                                   uint8(b),
                                                   uint8(a)})
                }
        }

        return o

}

func main() {
	amount := flag.Float64("amount", 1.2, "")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Println("Usage: contrast [--amount <num>] in.jpg out.jpg")
		os.Exit(1)
	}

	input, err := os.Open(args[0])
	defer input.Close()
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.Create("contrast-" + args[1])
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}

	output2, err := os.Create("grayscale-" + args[1])
        defer output2.Close()
        if err != nil {
                log.Fatal(err)
        }
	
	output3, err := os.Create("motionblur-" + args[1])
        defer output3.Close()
        if err != nil {
                log.Fatal(err)
        }

        output4, err := os.Create("invert-" + args[1])
        defer output4.Close()
        if err != nil {
                log.Fatal(err)
        }

	img, _, err := image.Decode(input)
	if err != nil {
		log.Fatal(err)
	}

	jpeg.Encode(output, contrast(img, *amount), &jpeg.Options{Quality: 80})
	jpeg.Encode(output2, grayscale(img), &jpeg.Options{Quality: 80})
	jpeg.Encode(output3, motionblur(img), &jpeg.Options{Quality: 80})
	jpeg.Encode(output4, invert(img), &jpeg.Options{Quality: 80})

}
