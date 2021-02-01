![backmeup logo](https://raw.githubusercontent.com/d-Rickyy-b/backmeup/master/docs/backmeup_logo_transparent.png)
# backmeup - a lightweight backup utility for the CLI
[![build](https://github.com/d-Rickyy-b/backmeup/workflows/build/badge.svg)](https://github.com/d-Rickyy-b/backmeup/actions?query=workflow%3Abuild)

When managing several servers, you often find yourself in a need of making backups. I searched for tools online that could make it into a painless experience but never found an easy-to-use, lightweight, portable, CLI tool which is easy to configure and does not need a remote server for backups.
That's why I created **backmeup**.

### Key features
- Easy to use
- Define multiple backups in a single config file
- **Portable** - you can copy the **single executable** with a configuration file on all your machines
- **Lightweight** - the executables are < 10 mb
- Exclude files and paths with .gitignore-like syntax
- Group together multiple source paths into one backup
- Config files written in yaml
- Usable from the CLI
- Multi-platform support

### Limitations
The goal of backmeup is not to replace professional backup tools. It cannot do incremental or differential backups. It doesn't know about any sort of backup strategies.
Also, it **doesn't** do any kind of deduplication! Apart from that, currently, there is no way to schedule your backups. 
To do so, you'd need to make use of external job schedulers, such as [cron](https://en.wikipedia.org/wiki/Cron).

The sole purpose of backmeup is to simplify basic backup functionality previously done via (e.g.) a shell script or manually via tar commands.
It does that by providing a simple config file and a single executable. 
If you want your backups to be synced to a remote machine, you must take care of that yourself or use other tools.

# Installation

## Download precompiled executable
Under the releases page you'll always find the most recent builds of backmeup as executables for linux, windows and macos (amd64).
These executables are built with `go build -ldflags="-s -w" .` and shrinked with [upx](https://github.com/upx/upx/).

In case that's not suitable for you, you can always build the executable from source.
 
## Compile from source
This tool is written in Go (golang), so you need to have Go installed on your system first. You can find [installation instructions here](https://golang.org/doc/install).

Then you need to clone the git repository.
Afterwards select the operating system you want to compile for, and the corresponding architecture (Go allows for easy cross compilation) from [this](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) list.
Set the environment variables `GOOS` and `GOARCH` with those values and eventually use `go build .` to generate an executable.
```bash
$ git clone https://github.com/d-Rickyy-b/backmeup
$ cd backmeup
$ GOOS="linux" && GOARCH="amd64"
$ go build -ldflags="-s -w" .
```
After that you have an executable named `backmeup` in your current working directory.

## Making it persistent
While backmeup is meant to be portable and can be executed from everywhere, you might want to permanently install it on your system.
On linux you can simply move the file into the `/bin/` directory.
```bash
$ mv backmeup /bin/backmeup
```

On Windows, you can create a new directory in the `Program Files` directory and move the binary there.
```
> mkdir C:\Program Files\backmeup\
> mv backmeup.exe C:\Program Files\backmeup\backmeup.exe
```

# Usage
All you need for this tool to work is a `config.yml` file. Simply pass the path to that file via the `-c`/`--config` switch. 
```
usage: backmeup [-h|--help] -c|--config "<value>" [-v|--verbose]

                The lightweight backup tool for the CLI

Arguments:

  -h  --help     Print help information
  -c  --config   Path to the config.yml file
  -v  --verbose  Enable verbose logging. Default: false
```

Starting your backups is as simple as this
```
$ backmeup -c config.yml
```

# How to create a config?
Configuring your backups is easy. Just create a `config.yml` file that contains the information about the sources and destination paths for your backups.

```yaml
backup_unit_name:
  sources:
    - C:\Users\admin\Documents\
    - 'C:\Users\admin\Dropbox'
    - "C:\\Users\\admin\\AppData\\Roaming\\.minecraft"
  destination: 'C:\backups'
  excludes:
    - "*.zip"
    - "*.rar"
  archive_type: "tar.gz"
  add_subfolder: false
  enabled: true
``` 
The minimal configuration consists only of the unit's name (obviously), at least one source and the destination path:

```yaml
backup_unit_name:
  sources:
    - C:\Users\admin\Documents\
  destination: 'C:\backups'
``` 

The name of your backup is the key at root level. Starting from there you can configure the following parameters.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| sources | list[strings] | Yes | | All paths to the directories you want to include in your backup |
| destination | string | Yes | | The destination directory, where the backup of this unit will be stored at |
| excludes | list[strings] | No | `[]` | .gitignore like filters for excluding files or dirs from the backup |
| archive_type | string | No | `tar.gz` | The type of archive to be used (`tar.gz` or `zip` are valid options) |
| add_subfolder | boolean | No | `false` | Creates a new subfolder in `<destination>` for this unit if set to true |
| enabled | boolean | No | `true` | Switch to disable each unit individually |
| use_absolute_paths | boolean | No | `true` | Uses absolute file paths in the archive (see [#11](https://github.com/d-Rickyy-b/backmeup/issues/11)) |

Be careful when using quotes in paths. For most strings you don't even need to use quotes at all. When using double quotes (`"`), you must escape backslashes (`\`) when you want to use them as literal characters (such as in Windows paths). 
Check [this handy article](https://www.yaml.info/learn/quote.html) for learning more about quotes in yaml.
