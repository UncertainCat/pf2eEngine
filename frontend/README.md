# PF2E Combat Simulator Frontend

This React application provides a user interface for the PF2E Combat Simulator.

## Requirements

- Node.js (v14+)
- npm (v6+)

## Setup

### Install Dependencies

```bash
npm install
```

### Development Mode

```bash
npm start
```

This will start the development server on port 3000. You'll need to run the Go backend separately.

### Production Build

```bash
npm run build
```

This creates a production build in the `build` directory, which will be served by the Go backend.

## Project Structure

- `src/` - Source code
    - `App.js` - Main application component
    - `App.css` - Styles
- `public/` - Static files
    - `index.html` - HTML template

## Technologies Used

- React
- Lucide React (for icons)
- Custom CSS for styling (Tailwind-inspired classes)