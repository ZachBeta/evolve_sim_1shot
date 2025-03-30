package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/renderer"
	"github.com/zachbeta/evolve_sim/pkg/simulation"
	"github.com/zachbeta/evolve_sim/pkg/world"
)

func main() {
	fmt.Printf("Evolutionary Simulator v%s\n", config.Version)
	fmt.Println("A simulation of single-cell organisms responding to chemical gradients")

	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to configuration file")
	headless := flag.Bool("headless", false, "Run in headless mode (no UI)")
	exportStats := flag.Bool("exportStats", false, "Export statistics to CSV and JSON")
	duration := flag.Float64("duration", 60.0, "Simulation duration in seconds (headless mode only)")
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to file")
	flag.Parse()

	// Start CPU profiling if requested
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Load configuration
	cfg, err := config.LoadFromFile(*configPath)
	if err != nil {
		// If the config file doesn't exist, try to create a default one
		if os.IsNotExist(err) {
			defaultCfg := config.DefaultConfig()
			if err := config.SaveToFile(defaultCfg, *configPath); err != nil {
				log.Fatalf("Failed to create default config: %v", err)
			}
			fmt.Printf("Created default configuration file at: %s\n", *configPath)
			cfg = defaultCfg
		} else {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	// Initialize the world
	world := world.NewWorld(cfg)

	// Initialize the simulator
	simulator := simulation.NewSimulator(world, cfg)

	// Initialize the renderer if not in headless mode
	if !*headless {
		gameRenderer := renderer.NewRenderer(world, simulator, cfg)

		// Set up Ebiten game
		ebiten.SetWindowSize(cfg.Render.WindowWidth, cfg.Render.WindowHeight)
		ebiten.SetWindowTitle("Evolution Simulator")
		ebiten.SetMaxTPS(cfg.Render.FrameRate)

		// Start the game
		if err := ebiten.RunGame(gameRenderer); err != nil {
			log.Fatalf("Failed to run game: %v", err)
		}
	} else {
		// Headless mode for batch processing or testing
		fmt.Println("Running in headless mode")
		runHeadless(simulator, *duration, *exportStats)
	}
}

// runHeadless executes the simulation without visualization
func runHeadless(simulator *simulation.Simulator, duration float64, exportStats bool) {
	// Calculate the number of steps needed
	// This assumes timestep is 1/60 (default)
	steps := int(duration / simulator.TimeStep)

	// Stats collection
	var stats []simulation.SimulationStats
	startTime := time.Now()

	// Progress reporting
	reportInterval := steps / 10
	if reportInterval < 1 {
		reportInterval = 1
	}

	// Run the simulation
	for i := 0; i < steps; i++ {
		simulator.Step()

		// Collect stats every 60 steps (approximately once per second)
		if i%60 == 0 {
			stat := simulator.CollectStats()
			stat.RealTimeElapsed = time.Since(startTime)
			stats = append(stats, stat)
		}

		// Report progress
		if i%reportInterval == 0 {
			progress := float64(i) / float64(steps) * 100
			fmt.Printf("Simulation progress: %.1f%% (time: %.2fs)\n", progress, simulator.Time)
		}
	}

	fmt.Printf("Simulation completed in %.2f seconds (simulation time: %.2fs)\n",
		time.Since(startTime).Seconds(), simulator.Time)

	// Export statistics if requested
	if exportStats && len(stats) > 0 {
		timestamp := time.Now().Format("20060102-150405")
		csvPath := fmt.Sprintf("stats_%s.csv", timestamp)
		jsonPath := fmt.Sprintf("stats_%s.json", timestamp)

		if err := simulation.ExportStatsCSV(stats, csvPath); err != nil {
			fmt.Printf("Failed to export CSV: %v\n", err)
		} else {
			fmt.Printf("Exported statistics to %s\n", csvPath)
		}

		if err := simulation.ExportStatsJSON(stats, jsonPath); err != nil {
			fmt.Printf("Failed to export JSON: %v\n", err)
		} else {
			fmt.Printf("Exported statistics to %s\n", jsonPath)
		}
	}
}
