# LaziSpace File Locations

This document tracks all files and directories created by LaziSpace for proper cleanup during uninstallation.

## Configuration Files

### Production (XDG Standard)
- **Linux**: `~/.config/lazispace/`
- **macOS**: `~/Library/Application Support/lazispace/`
- **Windows**: `%APPDATA%\lazispace\`

Contains:
- `registry.json` - Registered workspaces
- User preferences (future)

### Development
- `./dev-data/config/` - Development registry and config

### Test
- `./test-data/config/` - Test-specific data (auto-cleaned)

## Data Files

### Production (XDG Standard)
- **Linux**: `~/.local/share/lazispace/`
- **macOS**: `~/Library/Application Support/lazispace/`
- **Windows**: `%LOCALAPPDATA%\lazispace\`

Contains:
- Workspace metadata
- Cache files
- Temporary data

### Development
- `./dev-data/data/` - Development data

### Test
- `./test-data/data/` - Test data

## Log Files

### Production (XDG State)
- **Linux**: `~/.local/state/lazispace/logs/`
- **macOS**: `~/Library/Logs/lazispace/`
- **Windows**: `%LOCALAPPDATA%\lazispace\logs\`

Contains:
- `lazispace.log` - Current log file
- `lazispace-YYYY-MM-DD.log.gz` - Rotated logs (kept per retention policy)

### Development
- `./dev-data/logs/` - Development logs

### Test
- `./test-data/logs/` - Test logs

## Binary/Executable

Location depends on installation method:
- **Manual install**: User-specified (e.g., `/usr/local/bin/lspace`)
- **go install**: `$GOPATH/bin/lspace` or `~/go/bin/lspace`
- **Homebrew** (future): `/usr/local/bin/lspace`

## Uninstall Checklist

To completely remove LaziSpace:

1. Remove binary: `rm $(which lspace)`
2. Remove config: `rm -rf ~/.config/lazispace` (Linux/macOS)
3. Remove data: `rm -rf ~/.local/share/lazispace` (Linux)
4. Remove logs: `rm -rf ~/.local/state/lazispace` (Linux)
5. Unregister workspaces (removes system commands)

Future: `lspace uninstall --purge` command will automate this.
