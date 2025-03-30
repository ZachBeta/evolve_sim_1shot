# Triangle-Based Organism Representation Tutorial

This tutorial guides you through implementing a triangle-based visualization for organisms in our evolutionary simulator. This approach makes the organisms' direction more visually intuitive and improves the overall clarity of the simulation.

## Overview

Currently, our organisms are represented as circles with a line indicating direction. While functional, this approach has limitations:
- Direction can be hard to see at a glance
- Organism movement patterns are not as visually clear
- Rotating direction indicators don't communicate turning effectively

By changing to a triangle representation, we'll make organisms' movement more intuitive and visually appealing.

## Implementation Steps

### 1. Understanding the Current Implementation

First, examine the current organism drawing code in `pkg/renderer/renderer.go`:

```go
// Current implementation (simplified)
func (r *Renderer) drawOrganisms(screen *ebiten.Image) {
    organisms := r.World.GetOrganisms()

    for _, org := range organisms {
        // Convert world coordinates to screen coordinates
        screenX, screenY := r.worldToScreen(org.Position)

        // Determine color based on chemical preference
        // ...color calculation code...

        // Draw a small circle for the organism
        radius := 3.0
        for y := int(screenY) - int(radius); y <= int(screenY)+int(radius); y++ {
            for x := int(screenX) - int(radius); x <= int(screenX)+int(radius); x++ {
                dx := float64(x) - screenX
                dy := float64(y) - screenY
                if dx*dx+dy*dy <= radius*radius {
                    screen.Set(x, y, color.RGBA{red, green, blue, 255})
                }
            }
        }

        // Draw heading indicator
        headingX := screenX + math.Cos(org.Heading)*8
        headingY := screenY + math.Sin(org.Heading)*8
        ebitenutil.DrawLine(screen, screenX, screenY, headingX, headingY, color.RGBA{255, 255, 255, 200})
    }
}
```

### 2. Implementing Triangle Drawing

Let's implement a new function to draw triangles in our renderer package:

```go
// drawTriangle draws a filled triangle with the specified points and color
func (r *Renderer) drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float64, clr color.Color) {
    // Find the bounding box of the triangle
    minX := math.Min(x1, math.Min(x2, x3))
    maxX := math.Max(x1, math.Max(x2, x3))
    minY := math.Min(y1, math.Min(y2, y3))
    maxY := math.Max(y1, math.Max(y2, y3))

    // Iterate over each pixel in the bounding box
    for y := int(minY); y <= int(maxY); y++ {
        for x := int(minX); x <= int(maxX); x++ {
            // Check if the point is inside the triangle
            if pointInTriangle(float64(x), float64(y), x1, y1, x2, y2, x3, y3) {
                screen.Set(x, y, clr)
            }
        }
    }
}

// pointInTriangle determines if a point is inside a triangle using barycentric coordinates
func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 float64) bool {
    // Calculate area of the full triangle
    area := 0.5 * math.Abs((x2-x1)*(y3-y1) - (x3-x1)*(y2-y1))
    if area < 0.00001 {
        return false // Degenerate triangle
    }

    // Calculate barycentric coordinates
    alpha := 0.5 * math.Abs((x2-x3)*(py-y3) - (y2-y3)*(px-x3)) / area
    beta := 0.5 * math.Abs((x3-x1)*(py-y1) - (y3-y1)*(px-x1)) / area
    gamma := 1.0 - alpha - beta

    // Point is in triangle if all coordinates are between 0 and 1
    return alpha >= 0 && beta >= 0 && gamma >= 0 && alpha <= 1 && beta <= 1 && gamma <= 1
}
```

### 3. Modifying the Organism Drawing Code

Now, update the `drawOrganisms` method to use triangles instead of circles:

```go
func (r *Renderer) drawOrganisms(screen *ebiten.Image) {
    organisms := r.World.GetOrganisms()

    for _, org := range organisms {
        // Convert world coordinates to screen coordinates
        screenX, screenY := r.worldToScreen(org.Position)

        // Determine color based on chemical preference
        // ... existing color calculation code ...

        // Define triangle size (can be adjusted based on organism properties)
        size := 4.0 

        // Calculate triangle vertices
        // The triangle should point in the direction of heading
        // First point: front of the triangle (in heading direction)
        frontX := screenX + math.Cos(org.Heading) * size * 1.5
        frontY := screenY + math.Sin(org.Heading) * size * 1.5
        
        // Calculate the back corners (perpendicular to heading)
        backOffsetX := math.Cos(org.Heading + math.Pi/2) * size
        backOffsetY := math.Sin(org.Heading + math.Pi/2) * size
        
        // Left back corner
        leftX := screenX - math.Cos(org.Heading) * size/2 - backOffsetX
        leftY := screenY - math.Sin(org.Heading) * size/2 - backOffsetY
        
        // Right back corner
        rightX := screenX - math.Cos(org.Heading) * size/2 + backOffsetX
        rightY := screenY - math.Sin(org.Heading) * size/2 + backOffsetY

        // Draw the triangle
        r.drawTriangle(screen, frontX, frontY, leftX, leftY, rightX, rightY, 
                      color.RGBA{red, green, blue, 255})
                      
        // Optional: Add a border for better visibility
        ebitenutil.DrawLine(screen, frontX, frontY, leftX, leftY, color.RGBA{255, 255, 255, 200})
        ebitenutil.DrawLine(screen, leftX, leftY, rightX, rightY, color.RGBA{255, 255, 255, 200})
        ebitenutil.DrawLine(screen, rightX, rightY, frontX, frontY, color.RGBA{255, 255, 255, 200})
        
        // If the organism is selected (for future organism selection feature)
        if r.selectedOrganism != nil && r.selectedOrganism == &org {
            // Draw a larger highlight triangle
            // ... code to draw highlight ...
        }
    }
}
```

### 4. Adding Smooth Rotation Animation

To make the organisms' movement more visually pleasing, we can add smooth rotation animation. This involves:

1. Storing the previous heading for each organism
2. Interpolating between the previous and current heading when drawing

We'll need to modify the Organism struct to track the previous heading:

```go
// In pkg/types/organism.go
type Organism struct {
    Position        Point
    Heading         float64
    PreviousHeading float64  // Add this field
    ChemPreference  float64
    Speed           float64
    SensorAngles    [3]float64
    // Other fields...
}

// In the organism update logic (pkg/organism/movement.go or similar)
func (org *Organism) Update(dt float64) {
    // Store previous heading before updating
    org.PreviousHeading = org.Heading
    
    // ... existing update code ...
    
    // After updating heading, ensure we take the shortest path for rotation
    // This helps animation look better
    for org.Heading - org.PreviousHeading > math.Pi {
        org.PreviousHeading += 2 * math.Pi
    }
    for org.PreviousHeading - org.Heading > math.Pi {
        org.PreviousHeading -= 2 * math.Pi
    }
}
```

Then, in the rendering code, interpolate between the previous and current heading:

```go
// In drawOrganisms function
// Calculate the visual heading with interpolation for smooth rotation
visualHeading := org.PreviousHeading + (org.Heading - org.PreviousHeading) * r.interpolationFactor

// Use visualHeading instead of org.Heading for triangle calculations
frontX := screenX + math.Cos(visualHeading) * size * 1.5
// ... and so on
```

### 5. Optimizing the Implementation

For better performance, we can optimize the triangle drawing:

```go
// Add these fields to Renderer
type Renderer struct {
    // ... existing fields
    triangleImage *ebiten.Image
    triangleOpts  ebiten.DrawImageOptions
}

// Initialize in NewRenderer
func NewRenderer() *Renderer {
    r := &Renderer{
        // ... existing initialization
    }
    
    // Create a triangle image once
    r.triangleImage = ebiten.NewImage(16, 16)
    r.triangleOpts = ebiten.DrawImageOptions{}
    
    // ... rest of initialization
    return r
}

// Draw organisms using image transformation for better performance
func (r *Renderer) drawOrganisms(screen *ebiten.Image) {
    organisms := r.World.GetOrganisms()

    for _, org := range organisms {
        // ... existing coordinate and color calculation
        
        // Clear the triangle image
        r.triangleImage.Clear()
        
        // Draw triangle on the triangle image
        size := 4.0
        // ... calculate triangle points as before
        r.drawTriangle(r.triangleImage, 8+frontX-screenX, 8+frontY-screenY, 
                      8+leftX-screenX, 8+leftY-screenY, 
                      8+rightX-screenY, 8+rightY-screenY, 
                      color.RGBA{red, green, blue, 255})
        
        // Set up transformation
        r.triangleOpts.GeoM.Reset()
        r.triangleOpts.GeoM.Translate(screenX-8, screenY-8)
        
        // Draw the transformed triangle image
        screen.DrawImage(r.triangleImage, &r.triangleOpts)
    }
}
```

## Testing Your Implementation

After implementing these changes, test the simulation to ensure:

1. Organisms are correctly represented as triangles
2. Triangle direction matches the organism's heading
3. Rotation looks smooth during direction changes
4. Colors still properly represent chemical preferences
5. Performance remains good with many organisms

## Further Enhancements

Once basic triangle rendering is working, you might consider these enhancements:

1. **Size Variation**: Vary triangle size based on organism properties
2. **Animation Effects**: Add subtle pulsing or movement effects
3. **Selection Highlighting**: Add a special visual effect for selected organisms
4. **Trail Integration**: When implementing organism trails, ensure they connect to the triangle's position properly

## Conclusion

By replacing the circle+line representation with triangles, we've made the organisms' movement and direction much more intuitive to understand. This visual enhancement significantly improves the simulation's clarity and appeal.

The next logical step would be to implement organism trails, which will further enhance the visualization of movement patterns. 