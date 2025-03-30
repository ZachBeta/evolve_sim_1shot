package renderer

import (
	"image/color"
	"math"
)

// ColorScheme defines a color gradient to use for visualizations
type ColorScheme struct {
	Name        string
	Description string
	ColorStops  []ColorStop
}

// ColorStop defines a color at a specific position in the gradient
type ColorStop struct {
	Position float64 // 0.0 to 1.0
	Color    color.RGBA
}

// HSL represents a color in HSL color space
type HSL struct {
	H float64 // Hue (0-360)
	S float64 // Saturation (0-1)
	L float64 // Lightness (0-1)
}

// Predefined color schemes for chemical concentration visualization
var (
	// Viridis - perceptually uniform, increases linearly in lightness, colorblind friendly
	ViridisScheme = ColorScheme{
		Name:        "Viridis",
		Description: "Perceptually uniform, colorblind friendly",
		ColorStops: []ColorStop{
			{0.0, color.RGBA{68, 1, 84, 255}},    // Dark purple
			{0.25, color.RGBA{59, 82, 139, 255}}, // Blue/purple
			{0.5, color.RGBA{33, 144, 141, 255}}, // Teal
			{0.75, color.RGBA{93, 201, 99, 255}}, // Green
			{1.0, color.RGBA{253, 231, 37, 255}}, // Yellow
		},
	}

	// Magma - like viridis but with more dramatic contrast
	MagmaScheme = ColorScheme{
		Name:        "Magma",
		Description: "Higher contrast with dark-to-bright transition",
		ColorStops: []ColorStop{
			{0.0, color.RGBA{0, 0, 4, 255}},       // Almost black
			{0.25, color.RGBA{80, 18, 123, 255}},  // Deep purple
			{0.5, color.RGBA{182, 54, 121, 255}},  // Magenta
			{0.75, color.RGBA{251, 136, 97, 255}}, // Orange
			{1.0, color.RGBA{252, 253, 191, 255}}, // Pale yellow
		},
	}

	// Plasma - vibrant with dramatic hue transitions
	PlasmaScheme = ColorScheme{
		Name:        "Plasma",
		Description: "Vibrant with dramatic hue transitions",
		ColorStops: []ColorStop{
			{0.0, color.RGBA{13, 8, 135, 255}},    // Deep blue
			{0.25, color.RGBA{126, 3, 168, 255}},  // Purple
			{0.5, color.RGBA{204, 71, 120, 255}},  // Pink
			{0.75, color.RGBA{248, 149, 64, 255}}, // Orange
			{1.0, color.RGBA{240, 249, 33, 255}},  // Yellow
		},
	}

	// Classic - the original blue-green-red scheme
	ClassicScheme = ColorScheme{
		Name:        "Classic",
		Description: "Traditional blue-green-red gradient",
		ColorStops: []ColorStop{
			{0.0, color.RGBA{0, 0, 255, 255}}, // Blue
			{0.5, color.RGBA{0, 255, 0, 255}}, // Green
			{1.0, color.RGBA{255, 0, 0, 255}}, // Red
		},
	}
)

// GetColorFromScheme returns an interpolated color from the scheme at the given position (0-1)
func GetColorFromScheme(scheme ColorScheme, position float64) color.RGBA {
	// Clamp position to 0-1 range
	position = math.Max(0, math.Min(1, position))

	// If we have only one color stop or position is at the min/max, return that color
	if len(scheme.ColorStops) == 1 || position <= scheme.ColorStops[0].Position {
		return scheme.ColorStops[0].Color
	}

	if position >= scheme.ColorStops[len(scheme.ColorStops)-1].Position {
		return scheme.ColorStops[len(scheme.ColorStops)-1].Color
	}

	// Find the two color stops that this position falls between
	var leftStop, rightStop ColorStop
	for i := 0; i < len(scheme.ColorStops)-1; i++ {
		if scheme.ColorStops[i].Position <= position && position <= scheme.ColorStops[i+1].Position {
			leftStop = scheme.ColorStops[i]
			rightStop = scheme.ColorStops[i+1]
			break
		}
	}

	// Calculate the relative position between these two stops
	relativePos := (position - leftStop.Position) / (rightStop.Position - leftStop.Position)

	// Convert RGB colors to HSL for better interpolation
	leftHSL := RGBToHSL(leftStop.Color)
	rightHSL := RGBToHSL(rightStop.Color)

	// Interpolate in HSL space
	h := interpolateHue(leftHSL.H, rightHSL.H, relativePos)
	s := leftHSL.S + relativePos*(rightHSL.S-leftHSL.S)
	l := leftHSL.L + relativePos*(rightHSL.L-leftHSL.L)

	// Convert back to RGB
	return HSLToRGB(HSL{H: h, S: s, L: l})
}

// RGBToHSL converts an RGB color to HSL
func RGBToHSL(rgb color.RGBA) HSL {
	r := float64(rgb.R) / 255
	g := float64(rgb.G) / 255
	b := float64(rgb.B) / 255

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	l := (max + min) / 2

	var h, s float64

	if max == min {
		// Achromatic
		h = 0
		s = 0
	} else {
		d := max - min

		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}

		h *= 60
	}

	return HSL{H: h, S: s, L: l}
}

// HSLToRGB converts an HSL color to RGB
func HSLToRGB(hsl HSL) color.RGBA {
	var r, g, b float64

	if hsl.S == 0 {
		// Achromatic
		r = hsl.L
		g = hsl.L
		b = hsl.L
	} else {
		var q float64
		if hsl.L < 0.5 {
			q = hsl.L * (1 + hsl.S)
		} else {
			q = hsl.L + hsl.S - hsl.L*hsl.S
		}

		p := 2*hsl.L - q

		r = hueToRGB(p, q, hsl.H/360+1/3)
		g = hueToRGB(p, q, hsl.H/360)
		b = hueToRGB(p, q, hsl.H/360-1/3)
	}

	return color.RGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: 255,
	}
}

// hueToRGB is a helper function for HSLToRGB
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1/6 {
		return p + (q-p)*6*t
	}
	if t < 1/2 {
		return q
	}
	if t < 2/3 {
		return p + (q-p)*(2/3-t)*6
	}
	return p
}

// interpolateHue interpolates between two hue values, taking the shortest path
func interpolateHue(h1, h2, t float64) float64 {
	// Ensure h1 and h2 are in range [0, 360]
	h1 = math.Mod(h1, 360)
	if h1 < 0 {
		h1 += 360
	}

	h2 = math.Mod(h2, 360)
	if h2 < 0 {
		h2 += 360
	}

	// Find the shortest path
	diff := h2 - h1
	if diff > 180 {
		h1 += 360
	} else if diff < -180 {
		h2 += 360
	}

	return h1 + t*(h2-h1)
}
