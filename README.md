# PF2E Engine

PF2E Engine is a Go-based combat simulation system inspired by the Pathfinder 2E ruleset. It allows for initiative-based turn-taking combat between entities, with support for degrees of success, damage rolls, and structured phases of gameplay.

## Package Structure

### `main`
- **Purpose**: Entry point of the application. It initializes the combatants, manages the game loop, and orchestrates the phases of gameplay.
- **Key Responsibilities**:
    - Initializes entities and game state.
    - Runs the setup, combat, and resolution phases.

### `entity`
- **Purpose**: Defines and manages combat entities, including their attributes and behaviors.
- **Key Responsibilities**:
    - `Entity` struct: Encapsulates attributes like HP, AC, attack bonus, and damage bonus.
    - Methods for checking life status, taking damage, and rolling initiative.

### `attributes`
- **Purpose**: Adds a framework for managing additional attributes like Strength, Dexterity, etc.
- **Key Responsibilities**:
    - `Attributes` struct: Stores core ability scores.
    - Methods for applying modifiers dynamically to attributes.

### `combat`
- **Purpose**: Implements core combat mechanics, including attacks and initiative.
- **Key Responsibilities**:
    - `attack.go`: Handles dice rolls, calculates degrees of success, and executes attacks.
    - `initiative.go`: Rolls and sorts initiative for entities.

### `game`
- **Purpose**: Manages the overall game state, including turn order and game phases.
- **Key Responsibilities**:
    - `manager.go`: Maintains the current game state, manages turn order, and checks combat status.
    - `phases.go`: Defines structured phases (setup, combat, resolution) to organize gameplay flow.

### `dice`
- **Purpose**: Provides reusable utilities for dice rolling.
- **Key Responsibilities**:
    - Rolling single dice (`Roll`).
    - Rolling multiple dice (`RollMultiple`).
    - Rolling dice with modifiers (`RollWithModifier`).

## How It Works
1. **Setup Phase**:
    - Rolls initiative for all entities and determines turn order.
2. **Combat Phase**:
    - Entities take turns attacking their opponents based on initiative.
    - Attacks are resolved using dice rolls to determine degrees of success and damage.
3. **Resolution Phase**:
    - The game announces the winner or ends in a draw if no entities remain alive.

## Usage
1. Clone the repository.
2. Run the application using `go run main.go`.
3. Observe the combat simulation in the console output.

## Future Enhancements
- Add support for conditions and status effects.
- Expand `Attributes` to influence combat rolls and abilities.
- Include additional entity actions like defending or healing.
- Implement a graphical user interface (GUI) for a more interactive experience.