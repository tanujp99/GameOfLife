# Conway's Game of Life Simulator

A beautiful and interactive implementation of Conway's Game of Life built with Go and HTMX.

## Features

- **Interactive Grid**: Click any cell to toggle it alive/dead
- **Step-by-Step**: Advance the simulation one generation at a time
- **Auto-Run**: Start automatic simulation with configurable speed
- **Random Generation**: Generate random patterns
- **Clear Grid**: Reset the entire grid
- **Responsive Design**: Works on desktop and mobile devices
- **Real-time Updates**: Uses HTMX for seamless updates without page refreshes

## Conway's Game of Life Rules

The Game of Life follows these simple rules:

1. **Survival**: A live cell with 2 or 3 neighbors survives to the next generation
2. **Death**: A live cell with fewer than 2 neighbors dies (underpopulation)
3. **Death**: A live cell with more than 3 neighbors dies (overpopulation)
4. **Birth**: A dead cell with exactly 3 neighbors becomes alive (reproduction)

## Prerequisites

- Go 1.21 or later
- A modern web browser

## Installation

1. **Install Go** (if not already installed):
   - Download from [golang.org](https://golang.org/dl/)
   - Follow the installation instructions for your operating system

2. **Clone or download this project** to your local machine

3. **Navigate to the project directory**:
   ```bash
   cd gameoflife/go
   ```

## Running the Application

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Start the server**:
   ```bash
   go run main.go
   ```

3. **Open your browser** and navigate to:
   ```
   http://localhost:8080
   ```

## How to Use

### Basic Controls

- **Step**: Click to advance the simulation by one generation
- **Start Auto**: Begin automatic simulation (runs every 200ms)
- **Stop Auto**: Pause the automatic simulation
- **Random**: Generate a random pattern (30% cell density)
- **Clear**: Clear all cells from the grid

### Interactive Features

- **Click any cell** to toggle it between alive and dead states
- **Watch patterns evolve** according to Conway's rules
- **Experiment with different starting configurations**

### Grid Information

- **Grid Size**: 50x30 cells (1500 total cells)
- **Wrapping Edges**: The grid wraps around at the edges (toroidal surface)
- **Real-time Updates**: All changes are reflected immediately without page refreshes

## Technical Details

### Backend (Go)
- **Framework**: Standard Go HTTP server with custom router implementation
- **Game Logic**: Pure Go implementation of Conway's rules
- **API Endpoints**: RESTful endpoints for all game operations
- **State Management**: In-memory game state with thread-safe operations

### Frontend (HTMX)
- **Framework**: HTMX for dynamic updates without JavaScript complexity
- **Styling**: Modern CSS with responsive design
- **Interactivity**: Click-to-toggle cells with immediate feedback
- **Auto-run**: JavaScript-powered automatic simulation

### Architecture
- **Server-Side Rendering**: Go templates for initial page load
- **HTMX Integration**: Dynamic content updates via HTMX requests
- **RESTful API**: Clean separation between frontend and backend
- **Custom Router**: Lightweight router implementation without external dependencies
- **Responsive Design**: Mobile-friendly interface

## File Structure

```
gameoflife/
├── main.go              # Main Go server and game logic
├── router.go            # Custom router implementation
├── go.mod               # Go module dependencies
├── templates/
│   └── index.html       # Main HTML template with HTMX
├── static/
│   └── style.css        # CSS styling
└── README.md           # This file
```

## Customization

### Changing Grid Size
Edit the `Width` and `Height` values in `main.go`:
```go
gameState = GameState{
    Width:   50,  // Change this
    Height:  30,  // Change this
    Running: false,
    Step:    0,
}
```

### Adjusting Auto-Run Speed
Modify the interval in `templates/index.html`:
```javascript
setInterval(function() {
    htmx.ajax('POST', '/api/step', {target: '#grid', swap: 'innerHTML'});
}, 200); // Change this value (milliseconds)
```

### Changing Random Density
Modify the probability in `main.go`:
```go
gameState.Grid[y][x].Alive = rand.Float32() < 0.3 // Change 0.3 to desired density
```

## Browser Compatibility

This application works with all modern browsers that support:
- ES6 JavaScript
- CSS Grid
- HTMX (loaded from CDN)

## License

This project is open source and available under the MIT License.

## Acknowledgments

- **John Conway** for inventing the Game of Life
- **HTMX** for making dynamic web applications simple
- **Go Team** for the excellent programming language

## Contributing

Feel free to submit issues, feature requests, or pull requests to improve this implementation! 