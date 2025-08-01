package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	Alive bool `json:"alive"`
}

type GameState struct {
	Grid     [][]Cell `json:"grid"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Running  bool      `json:"running"`
	Step     int       `json:"step"`
}

var (
	gameState GameState
	indexTemplate *template.Template
	// Pre-allocate backup grid to avoid memory allocation on each step
	backupGrid [][]Cell
)

func main() {
	// Precompile template once at startup for better performance
	var err error
	indexTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
	gameState = GameState{
		Width:   50,
		Height:  30,
		Running: false,
		Step:    0,
	}
	
	gameState.Grid = make([][]Cell, gameState.Height)
	backupGrid = make([][]Cell, gameState.Height)
	for i := range gameState.Grid {
		gameState.Grid[i] = make([]Cell, gameState.Width)
		backupGrid[i] = make([]Cell, gameState.Width)
	}

	r := NewRouter()

	r.HandleFunc("GET", "/", homeHandler)
	r.HandleFunc("GET", "/api/state", getStateHandler)
	r.HandleFunc("POST", "/api/step", stepHandler)
	r.HandleFunc("POST", "/api/toggle", toggleCellHandler)
	r.HandleFunc("POST", "/api/clear", clearHandler)
	r.HandleFunc("POST", "/api/random", randomHandler)
	r.HandleFunc("POST", "/api/toggle-auto", toggleAutoHandler)
	r.HandleFunc("POST", "/api/auto-step", autoStepHandler)
	r.HandleFunc("GET", "/api/step-count", stepCountHandler)
	r.HandleFunc("GET", "/api/status", statusHandler)

	// Local Config
	// fmt.Println("Server starting on http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", r))

	// Remote Server Config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, gameState)
}

func getStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func stepHandler(w http.ResponseWriter, r *http.Request) {
	gameState.Grid = nextGeneration(gameState.Grid)
	gameState.Step++
	gameState.Running = false
	
	w.Header().Set("Content-Type", "text/html")
	renderGrid(w, gameState.Grid)
}

func toggleCellHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	x, err1 := strconv.Atoi(r.FormValue("x"))
	y, err2 := strconv.Atoi(r.FormValue("y"))
	
	if err1 != nil || err2 != nil || x < 0 || x >= gameState.Width || y < 0 || y >= gameState.Height {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}
	
	gameState.Grid[y][x].Alive = !gameState.Grid[y][x].Alive
	
	w.Header().Set("Content-Type", "text/html")
	renderCell(w, gameState.Grid[y][x], x, y)
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	// Clear existing grid instead of allocating new memory
	for y := 0; y < gameState.Height; y++ {
		for x := 0; x < gameState.Width; x++ {
			gameState.Grid[y][x].Alive = false
		}
	}
	gameState.Step = 0
	gameState.Running = false
	
	w.Header().Set("Content-Type", "text/html")
	renderGrid(w, gameState.Grid)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	
	for y := 0; y < gameState.Height; y++ {
		for x := 0; x < gameState.Width; x++ {
			gameState.Grid[y][x].Alive = rand.Float32() < 0.3
		}
	}
	gameState.Step = 0
	gameState.Running = false
	
	w.Header().Set("Content-Type", "text/html")
	renderGrid(w, gameState.Grid)
}

func toggleAutoHandler(w http.ResponseWriter, r *http.Request) {
	gameState.Running = !gameState.Running
	
	var buttonText string
	if gameState.Running {
		buttonText = "Auto Stop"
	} else {
		buttonText = "Auto Start"
	}
	
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<button class="btn btn-success" id="auto-toggle" hx-post="/api/toggle-auto" hx-target="#auto-controls" hx-swap="innerHTML">%s</button>`, buttonText)
}

func autoStepHandler(w http.ResponseWriter, r *http.Request) {
	if !gameState.Running {
		w.Header().Set("Content-Type", "text/html")
		renderGrid(w, gameState.Grid)
		return
	}
	
	gameState.Grid = nextGeneration(gameState.Grid)
	gameState.Step++
	
	w.Header().Set("Content-Type", "text/html")
	renderGrid(w, gameState.Grid)
}

func stepCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%d", gameState.Step)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	var statusText string
	if gameState.Running {
		statusText = "Status: Running"
	} else {
		statusText = "Status: Stopped"
	}
	
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<span class="status-text">%s</span>`, statusText)
}

func nextGeneration(grid [][]Cell) [][]Cell {
	height := len(grid)
	width := len(grid[0])
	
	// Reuse pre-allocated backup grid instead of creating new one
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			neighbors := countNeighbors(grid, x, y)
			alive := grid[y][x].Alive
			
			if alive && (neighbors == 2 || neighbors == 3) {
				backupGrid[y][x].Alive = true
			} else if !alive && neighbors == 3 {
				backupGrid[y][x].Alive = true
			} else {
				backupGrid[y][x].Alive = false
			}
		}
	}
	
	// Swap grids
	gameState.Grid, backupGrid = backupGrid, gameState.Grid
	return gameState.Grid
}

func countNeighbors(grid [][]Cell, x, y int) int {
	height := len(grid)
	width := len(grid[0])
	count := 0
	
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			
			ny := y + dy
			nx := x + dx
			
			if ny < 0 {
				ny = height - 1
			} else if ny >= height {
				ny = 0
			}
			
			if nx < 0 {
				nx = width - 1
			} else if nx >= width {
				nx = 0
			}
			
			if grid[ny][nx].Alive {
				count++
			}
		}
	}
	
	return count
}

func renderGrid(w http.ResponseWriter, grid [][]Cell) {
	var builder strings.Builder
	builder.Grow(len(grid) * len(grid[0]) * 120) // Pre-allocate capacity
	
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			class := "cell"
			if grid[y][x].Alive {
				class = "cell alive"
			}
			builder.WriteString(fmt.Sprintf(`<div class="%s" hx-post="/api/toggle" hx-vals='{"x":%d,"y":%d}' hx-target="this" hx-swap="outerHTML"></div>`, class, x, y))
		}
	}
	
	w.Write([]byte(builder.String()))
}

func renderCell(w http.ResponseWriter, cell Cell, x, y int) {
	class := "cell"
	if cell.Alive {
		class += " alive"
	}
	
	fmt.Fprintf(w, `<div class="%s" hx-post="/api/toggle" hx-vals='{"x":%d,"y":%d}' hx-target="this" hx-swap="outerHTML"></div>`, class, x, y)
}