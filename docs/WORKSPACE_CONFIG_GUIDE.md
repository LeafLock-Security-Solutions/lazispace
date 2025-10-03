# LaziSpace Workspace Configuration Guide

## Overview

LaziSpace uses YAML or JSON configuration files to define your development workspace. This guide explains all configuration options, validation rules, and provides comprehensive examples to help you create effective workspace configurations.

## File Format

LaziSpace supports both YAML and JSON formats:
- **YAML files**: `.yml` or `.yaml` extensions (recommended for readability)
- **JSON files**: `.json` extension (also supported)

The configuration structure is identical in both formats. Examples in this guide use YAML.

## What is LaziSpace

LaziSpace is a lightweight, cross-platform CLI tool for managing and automating local development environments. It solves the problem of manually opening multiple terminal windows, navigating to different directories, and starting various services every time you begin development work.

### The Problem LaziSpace Solves

As a developer working on modern applications (especially microservices), you typically need to:
- Open multiple terminal windows or tabs
- Navigate each to the correct directory
- Start different services (backend, frontend, database, etc.)
- Set up environment variables for each service
- Remember the exact commands to run
- Arrange terminal windows for efficient monitoring

This repetitive setup takes time and breaks your focus before you even start coding.

### How LaziSpace Works

LaziSpace automates this entire workflow:

1. **Define Once**: Create a configuration file describing your workspace layout
   - Which terminals (panes) you need
   - Where each should run (directories)
   - What commands to execute in each
   - How to arrange them visually (splits and tabs)

2. **Register**: Register your workspace configuration as a system command
   ```bash
   lspace register myproject
   ```

3. **Launch Anywhere**: Run your workspace from anywhere on your system
   ```bash
   myproject
   ```
   
   LaziSpace instantly:
   - Opens all configured terminal panes
   - Navigates each to the correct directory
   - Executes startup commands automatically
   - Arranges everything in your defined layout
   - Sets up all environment variables

### System Command Feature

When you register a workspace with `lspace register myproject`, LaziSpace creates a system command called `myproject` that launches your entire development environment. 

A **system command** is a command available globally in your terminal, like `ls`, `git`, or `docker`. After registration, your workspace name becomes such a command available from anywhere on your system.

**Example behavior:**
```bash
# Before registration - command doesn't exist
$ myproject
bash: myproject: command not found

# Register the workspace
$ lspace register myproject

# After registration - launch from anywhere
$ cd ~
$ myproject
# LaziSpace launches your complete development environment:
# - All services running
# - All terminals in correct directories
# - All commands executed
# - Everything organized in your layout

# Works from any location
$ cd /some/other/project
$ myproject
# Your workspace starts in the correct directories with all services running
```

## Workspace Components

### What is a Pane

A **pane** is an individual terminal window within your workspace. Each pane:
- Runs in a specific directory
- Can execute startup commands automatically
- Has its own shell and environment variables
- Represents one terminal session (like a tab in your terminal app)

Think of a pane as one terminal window where you might run a specific service, command, or task.

### What is a Split

A **split** is a container that divides screen space between multiple child elements (panes or other splits). Splits determine the visual layout:
- **Horizontal split**: Divides space left-to-right (side-by-side panes)
- **Vertical split**: Divides space top-to-bottom (stacked panes)
- Splits can be nested to create complex layouts
- Each split must contain at least 2 children

Example layout with splits:
```
┌─────────────┬─────────────┐
│    Pane 1   │    Pane 2   │  ← Horizontal split
├─────────────┼─────────────┤
│    Pane 3   │    Pane 4   │
└─────────────┴─────────────┘
```

## Basic Structure

Every LaziSpace configuration file contains these main sections:

```yaml
name: workspace-name                # Required: becomes system command name
title: "Display Name"               # Optional: human-readable title
rootDir: /path/to/project          # Optional: base directory (see path rules)
defaultShell: /bin/bash            # Optional: default shell for all panes
globalEnv:                         # Optional: environment variables for all panes
  NODE_ENV: development
commandHistory:                    # Optional: global command shortcuts
  - npm start
  - npm test

layout:                            # Required: workspace layout definition
  tabs:
    - title: "Development"         # Optional: tab title (auto-generated if empty)
      active: true                 # Optional: mark as active tab
      children:                    # Required: at least one child per tab
        - pane:                    # Either 'pane' or 'split'
            title: "Server"        # Optional: pane title (auto-generated if empty)
            path: ./backend        # Optional: working directory (see path rules)
            command:               # Optional: startup command
              command: "npm start"
              type: "run"          # Must be: "run" or "fill"
            shell: /bin/zsh        # Optional: pane-specific shell
            env:                   # Optional: pane-specific environment variables
              PORT: 3000
            commandHistory:        # Optional: pane-specific command shortcuts
              - npm run dev
              - npm test
```

## Field Validation Rules

### Workspace Name (System Command Creation)

The `name` field becomes a **system command** after registration. This command will be available globally in your terminal, so it must follow strict naming rules to avoid conflicts.

**Requirements:**
- Required field
- Must start with a letter (a-z, A-Z)
- Can contain: letters, numbers, underscores, single hyphens
- Cannot contain: spaces, special characters, consecutive hyphens
- Cannot start or end with hyphens
- Cannot conflict with existing system commands

**Valid Examples:**
```yaml
name: my-project        # Creates command: my-project
name: backend           # Creates command: backend
name: api-server-v2     # Creates command: api-server-v2
name: user_service      # Creates command: user_service
```

**Invalid Examples:**
```yaml
name: 123project        # Invalid: starts with number
name: my--project       # Invalid: consecutive hyphens
name: -myproject        # Invalid: starts with hyphen
name: myproject-        # Invalid: ends with hyphen
name: my project        # Invalid: contains space
name: ls                # Invalid: conflicts with system command
name: git               # Invalid: conflicts with system command
```

**System Command Conflict Check:**
LaziSpace validates that your chosen name doesn't conflict with existing commands like `ls`, `git`, `docker`, `npm`, etc.

### Workspace Title

Optional human-readable display name for the workspace.

**Accepted Values:**
- Any string up to 30 characters
- Empty string (uses workspace name)
- Can be omitted entirely

**Examples:**
```yaml
title: "My Development Environment"    # Valid
title: "API Server Development"        # Valid
title: ""                             # Valid: will use workspace name
# title field can be omitted entirely
```

### Root Directory (rootDir)

Base directory for the workspace. Required in specific scenarios based on pane path usage.

**Path Resolution Rules:**
1. **Empty pane paths**: If any pane has empty `path`, `rootDir` must be provided
2. **Relative pane paths**: If any pane has relative `path`, `rootDir` must be provided
3. **All absolute pane paths**: `rootDir` is optional

**Requirements when provided:**
- Must be absolute path (starts with `/` on Unix, `C:\` on Windows)
- Directory must exist and be accessible

**Examples:**

Valid configurations:
```yaml
# Case 1: rootDir required (pane has empty path)
rootDir: /home/user/myproject
layout:
  tabs:
    - children:
        - pane:
            path: ""              # Empty path requires rootDir

# Case 2: rootDir required (pane has relative path)
rootDir: /home/user/myproject
layout:
  tabs:
    - children:
        - pane:
            path: ./backend       # Relative path requires rootDir

# Case 3: rootDir optional (all absolute paths)
layout:
  tabs:
    - children:
        - pane:
            path: /home/user/myproject/backend    # Absolute path, no rootDir needed
```

### Shell Configuration

LaziSpace uses a three-tier shell resolution system with clear precedence rules.

**Shell Resolution Precedence (highest to lowest):**
1. **Pane `shell`** (highest priority) - overrides everything
2. **Workspace `defaultShell`** (fallback) - used when pane shell is empty
3. **System default shell** (automatic fallback) - used when both above are empty

**Requirements when specified:**
- Must be absolute path (e.g., `/bin/bash`, `/usr/bin/zsh`)
- Executable must exist on the system
- File must not be a directory

**System Default Shell Detection:**
LaziSpace automatically detects these shells in this order:
1. `/bin/bash`
2. `/usr/bin/bash`
3. `/bin/zsh`
4. `/usr/bin/zsh`
5. `/bin/sh`
6. `/usr/bin/sh`

**Accepted Values:**
- Absolute path to shell executable (e.g., `/bin/bash`)
- Empty string or omitted (uses fallback)

**Examples:**

```yaml
# Example 1: defaultShell covers all panes
defaultShell: /bin/bash           # All panes use bash unless overridden
layout:
  tabs:
    - children:
        - pane:
            title: "Server"       # Uses /bin/bash (from defaultShell)
        - pane:
            title: "Client"
            shell: /bin/zsh       # Uses /bin/zsh (overrides defaultShell)

# Example 2: Mixed shell specification
defaultShell: /bin/bash
layout:
  tabs:
    - children:
        - pane:                   # Uses /bin/bash (from defaultShell)
            title: "Backend"
        - pane:
            shell: /usr/bin/fish  # Uses /usr/bin/fish (overrides defaultShell)
            title: "Database"

# Example 3: No shells specified (uses system default)
layout:
  tabs:
    - children:
        - pane:
            title: "Server"       # Uses system default (/bin/bash if available)
```

### Environment Variables

Environment variables can be defined at workspace level (global) and pane level, with clear inheritance and precedence rules.

**Scope and Inheritance:**
- **globalEnv**: Available to ALL panes in the workspace
- **Pane env**: Available only to that specific pane
- **Inheritance**: Pane env inherits from globalEnv
- **Precedence**: Pane env values override globalEnv values with same key

**Variable Name Requirements:**
- Must start with letter (a-z, A-Z) or underscore (_)
- Can contain letters, numbers, underscores
- Cannot be empty string
- Case-sensitive

**Variable Value Requirements:**
- Any string value (no length restrictions)
- Can be empty string

**Examples:**

```yaml
# Global environment variables (inherited by all panes)
globalEnv:
  NODE_ENV: development      # Available to all panes
  LOG_LEVEL: debug          # Available to all panes
  API_URL: https://api.example.com

layout:
  tabs:
    - children:
        - pane:
            title: "API Server"
            env:
              PORT: 3000         # Pane-specific variable
              NODE_ENV: production # OVERRIDES global NODE_ENV for this pane only
            # Final env for this pane:
            # PORT=3000, NODE_ENV=production, LOG_LEVEL=debug, API_URL=https://api.example.com
            
        - pane:
            title: "Frontend"
            env:
              PORT: 3001         # Different port
              REACT_APP_API_URL: http://localhost:3000
            # Final env for this pane:
            # PORT=3001, NODE_ENV=development, LOG_LEVEL=debug, API_URL=https://api.example.com, REACT_APP_API_URL=http://localhost:3000
            
        - pane:
            title: "Database"
            # No pane env specified
            # Final env for this pane (inherits only global):
            # NODE_ENV=development, LOG_LEVEL=debug, API_URL=https://api.example.com
```

### Command History

Command history provides predefined command shortcuts for quick access. These are NOT command histories in the traditional shell sense, but rather quick command templates.

**Purpose:**
- Store frequently used commands for quick access
- Cycle through common commands using keyboard shortcuts
- Fill commands into terminal (with option to edit before running)
- Quick reminders of useful commands for the project

**Scope and Inheritance:**
- **Workspace commandHistory**: Available globally, inherited by all panes
- **Pane commandHistory**: Extends the workspace commandHistory for that pane
- **Inheritance**: Pane inherits global commandHistory and can add more
- **Merge behavior**: Pane commandHistory is appended to workspace commandHistory

**Requirements:**
- Maximum 1000 entries per list (workspace or pane)
- Cannot contain empty strings or whitespace-only commands
- Each entry must be a valid string

**Examples:**

```yaml
# Global command shortcuts (inherited by all panes)
commandHistory:
  - git status
  - git pull origin main
  - docker-compose up -d
  - docker-compose down

layout:
  tabs:
    - children:
        - pane:
            title: "API Server"
            path: ./backend
            commandHistory:           # Extends global commandHistory
              - npm run dev
              - npm run test
              - npm run lint
              - nodemon server.js
            # This pane has access to 8 commands total:
            # Global: git status, git pull origin main, docker-compose up -d, docker-compose down
            # Pane:   npm run dev, npm run test, npm run lint, nodemon server.js
            
        - pane:
            title: "Frontend"
            path: ./frontend
            commandHistory:           # Extends global commandHistory
              - npm start
              - npm run build
              - npm run test:coverage
            # This pane has access to 7 commands total:
            # Global: git status, git pull origin main, docker-compose up -d, docker-compose down
            # Pane:   npm start, npm run build, npm run test:coverage
            
        - pane:
            title: "Database"
            path: ./database
            # No commandHistory specified
            # This pane has access to 4 commands (global only):
            # git status, git pull origin main, docker-compose up -d, docker-compose down
```

**Usage Pattern:**
Users will cycle through these commands using keyboard shortcuts (to be defined later) and can:
1. Fill the command into the terminal
2. Edit the command if needed
3. Execute directly or manually

## Layout Structure

The layout defines the visual organization of your workspace using tabs, splits, and panes.

### Tabs

Workspaces must have at least one tab. Tabs organize your workspace into logical sections.

**Requirements:**
- Minimum 1 tab required
- Each tab must have at least 1 child
- Maximum 1 tab can be marked as `active: true`
- Tab titles are optional (auto-generated if empty)
- Maximum title length: 20 characters

**Active Tab Rules:**
- **No active tabs specified**: First tab automatically becomes active
- **One active tab**: Works correctly
- **Multiple active tabs**: Validation error

**Accepted Values for `active`:**
- `true`: Mark this tab as active
- `false` or omitted: Tab is not active

**Examples:**

```yaml
layout:
  tabs:
    - title: "Development"
      active: true                # Only one tab can be active
      children: [...]
    - title: "Testing"
      children: [...]             # active defaults to false
    - title: "Monitoring"
      active: false               # Explicitly set to false
      children: [...]

# Auto-generated titles and auto-active
layout:
  tabs:
    - children: [...]             # Title becomes "Tab 1", automatically active
    - children: [...]             # Title becomes "Tab 2"
```

### Layout Nodes

Each tab contains children that are either panes or splits. Each child must be exactly one type.

**Union Constraint:**
Each layout node must define exactly one of:
- `pane`: Terminal configuration
- `split`: Container for multiple children

**Invalid Configurations:**
```yaml
# WRONG - has both pane and split
- pane: { title: "Server" }
  split: { direction: horizontal, children: [...] }

# WRONG - has neither
- title: "Something"

# CORRECT - exactly one
- pane: { title: "Server" }

# CORRECT - exactly one  
- split: { direction: horizontal, children: [...] }
```

### Panes

Panes represent individual terminal windows with specific configurations.

**Pane Fields:**
```yaml
pane:
  title: "Server"                 # Optional: max 30 chars, auto-generated if empty
  path: ./backend                 # Optional: working directory (see path rules)
  command:                        # Optional: startup command
    command: "npm start"          # Command string (any valid shell command)
    type: "run"                   # Must be: "run" or "fill"
  shell: /bin/zsh                # Optional: absolute path to shell
  env:                           # Optional: environment variables (key-value pairs)
    PORT: 3000
  commandHistory:                # Optional: command shortcuts (max 1000 entries)
    - npm run dev
    - npm test
```

**Path Resolution for Panes:**

1. **Empty path (`""` or omitted)**: Uses workspace `rootDir`
   ```yaml
   rootDir: /home/user/project
   pane:
     path: ""                     # Resolves to: /home/user/project
   ```

2. **Relative path**: Joined with workspace `rootDir`
   ```yaml
   rootDir: /home/user/project
   pane:
     path: ./backend              # Resolves to: /home/user/project/backend
     # path: ../other            # Resolves to: /home/user/other
   ```

3. **Absolute path**: Used directly, `rootDir` ignored
   ```yaml
   pane:
     path: /opt/services/api      # Used exactly as specified
   ```

**Command Configuration:**

**`command.command` Accepted Values:**
- Any valid shell command string
- Empty string (no startup command)
- Multi-line commands using YAML string syntax

**`command.type` Accepted Values (Required if command specified):**
- `"run"`: Execute command immediately when pane starts
- `"fill"`: Type command into pane but don't execute (user can edit/run manually)

**Examples:**
```yaml
pane:
  command:
    command: "npm start"
    type: "run"                   # Executes immediately

pane:
  command:
    command: "npm run dev"
    type: "fill"                  # Types command, waits for user to press Enter

pane:
  command:
    command: |                    # Multi-line command
      export NODE_ENV=development
      npm start
    type: "run"

pane:
  command:
    command: ""                   # No startup command
    type: "run"
```

### Splits

Splits divide space between multiple children and can be nested to create complex layouts.

**Split Fields:**
```yaml
split:
  direction: "horizontal"         # Required: "horizontal" or "vertical"
  size: 50                       # Optional: percentage (0-100)
  children:                      # Required: minimum 2 children
    - pane: { ... }
    - pane: { ... }
```

**`direction` Accepted Values (Required):**
- `"horizontal"`: Divides space left-to-right (side-by-side layout)
- `"vertical"`: Divides space top-to-bottom (stacked layout)

**`size` Accepted Values (Optional):**
- Integer from 0 to 100 (percentage of parent space)
- Omitted: Equal distribution among siblings

**Split Requirements:**
- `direction` is required
- Minimum 2 children required
- No maximum children limit
- Maximum nesting depth: 10 levels (prevents infinite recursion)
- Children can be panes or nested splits

**Layout Examples:**

Horizontal split (side-by-side):
```yaml
split:
  direction: "horizontal"
  children:
    - pane: { title: "Left" }
    - pane: { title: "Right" }
```

Vertical split (stacked):
```yaml
split:
  direction: "vertical"  
  children:
    - pane: { title: "Top" }
    - pane: { title: "Bottom" }
```

Nested splits:
```yaml
split:
  direction: "horizontal"
  children:
    - pane: { title: "Left" }
    - split:
        direction: "vertical"
        children:
          - pane: { title: "Top Right" }
          - pane: { title: "Bottom Right" }
```

## Complete Examples

### Simple Single Service

```yaml
name: simple-app
title: "Simple Application"
rootDir: /home/user/simple-app
defaultShell: /bin/bash

layout:
  tabs:
    - title: "Development"
      active: true
      children:
        - pane:
            title: "Server"
            command:
              command: "npm start"
              type: "run"
```

### Complex Microservices Setup

```yaml
name: microservices
title: "Microservices Development Environment"
rootDir: /home/user/microservices
defaultShell: /bin/zsh

# Global environment (inherited by all panes)
globalEnv:
  NODE_ENV: development
  LOG_LEVEL: debug
  DATABASE_URL: postgresql://localhost:5432/devdb

# Global command shortcuts (inherited by all panes)
commandHistory:
  - docker-compose up -d
  - docker-compose down
  - git status
  - git pull origin main

layout:
  tabs:
    - title: "Services"
      active: true
      children:
        - split:
            direction: "horizontal"
            children:
              - split:
                  direction: "vertical"
                  children:
                    - pane:
                        title: "API Gateway"
                        path: ./gateway
                        command:
                          command: "npm run dev"
                          type: "run"
                        env:
                          PORT: 8080          # Overrides any global PORT if set
                        commandHistory:       # Extends global commandHistory
                          - npm run dev
                          - npm run test
                          - npm run lint
                        # Total commands: 7 (4 global + 3 pane)
                    - pane:
                        title: "User Service"
                        path: ./services/user
                        command:
                          command: "npm run dev"
                          type: "run"
                        env:
                          PORT: 3001
                          NODE_ENV: production # Overrides global NODE_ENV
                        # No pane commandHistory - has 4 global commands only
              - split:
                  direction: "vertical"
                  children:
                    - pane:
                        title: "Order Service"
                        path: ./services/order
                        command:
                          command: "npm run dev"
                          type: "run"
                        env:
                          PORT: 3002
                    - pane:
                        title: "Frontend"
                        path: ./frontend
                        command:
                          command: "npm start"
                          type: "run"
                        env:
                          REACT_APP_API_URL: http://localhost:8080
                        commandHistory:
                          - npm start
                          - npm run build
                          - npm run test:coverage
                        # Total commands: 7 (4 global + 3 pane)
    - title: "Database"
      children:
        - split:
            direction: "horizontal"
            children:
              - pane:
                  title: "Database Logs"
                  command:
                    command: "docker logs -f postgres-container"
                    type: "run"
              - pane:
                  title: "Database Client"
                  command:
                    command: "psql postgresql://localhost:5432/devdb"
                    type: "fill"          # User can edit connection string
    - title: "Monitoring"
      children:
        - split:
            direction: "vertical"
            children:
              - pane:
                  title: "System Monitor"
                  command:
                    command: "htop"
                    type: "run"
              - pane:
                  title: "Application Logs"
                  command:
                    command: "tail -f logs/application.log"
                    type: "run"
```

## Common Configuration Patterns

### Pattern 1: Service + Database + Logs

```yaml
split:
  direction: "horizontal"
  children:
    - pane:
        title: "Application"
        command: { command: "npm start", type: "run" }
    - split:
        direction: "vertical"
        children:
          - pane:
              title: "Database"
              command: { command: "docker-compose up db", type: "run" }
          - pane:
              title: "Logs"
              command: { command: "tail -f logs/app.log", type: "run" }
```

### Pattern 2: Multiple Services Side by Side

```yaml
split:
  direction: "horizontal"
  children:
    - pane:
        title: "Service A"
        path: ./service-a
        command: { command: "npm run dev", type: "run" }
    - pane:
        title: "Service B"  
        path: ./service-b
        command: { command: "npm run dev", type: "run" }
    - pane:
        title: "Service C"
        path: ./service-c
        command: { command: "npm run dev", type: "run" }
```

## Troubleshooting Common Issues

### Issue: "rootDir is required"
**Cause**: Pane has empty or relative path but no rootDir specified  
**Solution**: Add rootDir or use absolute paths in all panes

```yaml
# Fix by adding rootDir
rootDir: /home/user/project
layout:
  tabs:
    - children:
        - pane:
            path: ./backend     # Now valid with rootDir

# Or fix by using absolute path
layout:
  tabs:
    - children:
        - pane:
            path: /home/user/project/backend    # Absolute path, no rootDir needed
```

### Issue: "only one tab can be marked as active"
**Cause**: Multiple tabs have `active: true`  
**Solution**: Remove `active: true` from all but one tab

```yaml
# Fix: only one active tab
layout:
  tabs:
    - title: "Development"
      active: true            # Only this one
      children: [...]
    - title: "Testing"
      # Remove: active: true
      children: [...]
```

### Issue: "must have either 'pane' or 'split' defined"
**Cause**: Layout node has neither pane nor split, or has both  
**Solution**: Define exactly one

```yaml
# Wrong - has both
- pane: { title: "Server" }
  split: { direction: "horizontal", children: [...] }

# Right - exactly one
- pane: { title: "Server" }

# Or exactly one
- split: { direction: "horizontal", children: [...] }
```

### Issue: "startup command type must be 'run' or 'fill'"
**Cause**: Invalid value for `command.type`  
**Solution**: Use only accepted values

```yaml
# Wrong
command:
  command: "npm start"
  type: "execute"         # Invalid

# Right
command:
  command: "npm start"  
  type: "run"             # Valid: "run" or "fill" only
```

### Issue: "split direction must be 'horizontal' or 'vertical'"
**Cause**: Invalid value for `split.direction`  
**Solution**: Use only accepted values

```yaml
# Wrong
split:
  direction: "sideways"   # Invalid

# Right  
split:
  direction: "horizontal" # Valid: "horizontal" or "vertical" only
```

## Validation Command

Before registering your workspace, validate the configuration:

```bash
lspace validate workspace.yml
```

This will show all validation errors with specific locations and suggestions for fixes.