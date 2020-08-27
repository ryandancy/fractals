package cli

import (
	"fmt"
	"github.com/copperium/fractals/fractal"
	"github.com/integrii/flaggy"
	"image/png"
	"os"
	"runtime"
)

const subcommandName = "gen"

var (
	Subcommand  *flaggy.Subcommand
	fractalType string
	left        = 0.0
	bottom      = 0.0
	width       = 1.0
	height      = 1.0
	threshold   = 1000.0
	paramStr    = "0+0i"
	iters       = 100
	imgWidth    = 1000
	colors      = "blue-to-yellow"
	boldMode    = false
)

func init() {
	Subcommand = flaggy.NewSubcommand(subcommandName)
	Subcommand.Description = "generate a fractal from the command line"
	Subcommand.AddPositionalValue(&fractalType, "fractal", 1, true, "Type of fractal: mandelbrot or julia")
	Subcommand.Float64(&left, "l", "left", "Left bound of image in fractal")
	Subcommand.Float64(&bottom, "b", "bottom", "Bottom bound of image in fractal")
	Subcommand.Float64(&width, "W", "width", "Width of image in fractal")
	Subcommand.Float64(&height, "H", "height", "Height of image in fractal")
	Subcommand.Float64(&threshold, "t", "threshold", "Fractal computation threshold")
	Subcommand.String(&paramStr, "p", "param", "Complex parameter for Julia fractal")
	Subcommand.Int(&iters, "i", "iters", "Number of fractal iterations")
	Subcommand.Int(&imgWidth, "w", "image-width", "Width of image (pixels)")
	Subcommand.String(&colors, "c", "colors", "Color scheme: blue-to-yellow, red-to-green, or greyscale")
	Subcommand.Bool(&boldMode, "bm", "bold-mode", "Enable bold mode!")
}

func Exec() {
	var frac fractal.Fractal
	switch fractalType {
	case "mandelbrot":
		frac = &fractal.Mandelbrot{Threshold: threshold}
	case "julia":
		var realParam, imagParam float64
		_, err := fmt.Sscanf(paramStr, "%f%fi", &realParam, &imagParam)
		if err != nil {
			flaggy.ShowHelpAndExit("Invalid complex format")
		}
		frac = &fractal.Julia{Threshold: threshold, Param: complex(realParam, imagParam)}
	default:
		flaggy.ShowHelpAndExit("Unknown fractal type: options are mandelbrot, julia")
	}

	bounds := fractal.Rect{
		BottomLeft: &fractal.Point{X: left, Y: bottom},
		TopRight:   &fractal.Point{X: left + width, Y: bottom + height},
	}
	pixelSize := (bounds.TopRight.X - bounds.BottomLeft.X) / float64(imgWidth)

	var colorModel fractal.ColorModel
	switch colors {
	case "greyscale":
		colorModel = &fractal.GreyscaleColorModel{Threshold: iters}
	case "blue-to-yellow":
		colorModel = &fractal.HueColorModel{Threshold: iters, HueRange: fractal.BlueToYellow, BoldMode: boldMode}
	case "red-to-green":
		colorModel = &fractal.HueColorModel{Threshold: iters, HueRange: fractal.RedToGreen, BoldMode: boldMode}
	default:
		flaggy.ShowHelpAndExit("Unknown colors: options are greyscale, blue-to-yellow, red-to-green")
	}

	viz := fractal.Image{
		Fractal:       frac,
		Model:         colorModel,
		FractalBounds: bounds,
		Iters:         iters,
		PixelSize:     pixelSize,
	}
	cached := viz.ToCachedImage(runtime.NumCPU())

	err := png.Encode(os.Stdout, cached)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
	}
}
