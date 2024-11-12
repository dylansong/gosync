# GoSync - File Synchronization Tool

GoSync is a simple yet powerful file synchronization tool written in Go that allows you to synchronize files from a source directory to multiple target directories simultaneously.

## Features
- Synchronize files from one source to multiple destinations
- YAML-based configuration
- Support for multiple sync configurations in a single file
- Easy to use command-line interface

## Installation
```bash
curl -L https://github.com/dylansong/gosync/releases/download/v1.0.0/v1.0.0.tar.gz | tar -xz -C /usr/local/bin
```

## Usage

### 1. Create a Configuration File
Create a `config.yaml` file with your sync configurations:

```yaml
sync_configs:
  - name: "project_sync"
    source_dir: "/path/to/source"
    target_dirs:
      - "/path/to/target1"
      - "/path/to/target2"
```

### 2. Run GoSync
You can run GoSync in several ways:

```bash
# Use default config.yaml in current directory
gosync

# Specify a custom config file
gosync -config /path/to/config.yaml

# Generate a sample configuration
gosync -example > config.yaml
```

### Command Line Options
- `-config`: Specify the path to your configuration file (default: `config.yaml`)
- `-example`: Display a sample configuration file
- `-help`: Show help information

## Configuration File Structure
- `name`: A unique identifier for each sync configuration
- `source_dir`: The directory containing files to be synchronized
- `target_dirs`: A list of target directories where files will be copied to

The tool will maintain the same directory structure in target directories as in the source directory.

## Development

### Build from source
```bash
go build -o gosync main.go
```

### Create release archive
```bash
mkdir -p archive && tar -czf archive/v1.0.0.tar.gz gosync
```

### Create GitHub release
```bash
hub release create \
    -m "Release v1.0.0" \
    -m "Release description" \
    -a ./archive/v1.0.0.tar.gz \
    v1.0.0
```

## archive
```bash
mkdir -p archive && tar -czf archive/v1.0.0.tar.gz gosync
```

//install on local
```bash
./deploy-local.sh
```
