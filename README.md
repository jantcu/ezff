# ezff

A minimal Go wrapper for FFmpeg to simplify common video trimming tasks via the command line. Makes it easier to remember and execute basic cuts without typing full FFmpeg commands.

## Features

- **trim**: Trim both the beginning and end of a video.
- **trim-start**: Trim the beginning of a video.
- **trim-end**: Trim the end of a video.
- **trim-mid**: Remove a middle section of a video.
- Automatically checks for FFmpeg installation.
- Generates output filenames if not specified (appends `_trim` and increments if needed).

## Prerequisites

- [FFmpeg](https://ffmpeg.org/) installed and in your PATH.
- [Go](https://go.dev/) (version 1.21 or later) for building.

## Installation

1. Clone the repository: `git clone git@github.com:jantcu/ezff.git`
2. Build the binary: `cd ezff && go build`
3. Move the binary to a directory in your PATH (e.g., `/usr/local/bin`): `sudo mv ezff /usr/local/bin/`

Or instead of step 3, use `go install` if you structure it as a module (see [Go docs](https://go.dev/doc/install) for details).

## Usage

```
ezff <command> [args]
```

Run `ezff` without arguments for a quick usage summary.

### Commands

#### trim
Trims specified seconds from both the beginning and end of the video. Uses stream copying for speed and calculates the effective duration using `ffprobe`.

**Syntax:** `ezff trim <input> <trim_start> <trim_end> [--output <output>]`

- `<input>`: Path to input video file.
- `<trim_start>`: Float seconds to trim from the start.
- `<trim_end>`: Float seconds to trim from the end.
- `--output` (optional): Output file path. Defaults to `<input>_trim.mp4` (or increments if already exists).

**Example:** (removes first 3 seconds and last 5 seconds)

```
ezff trim myvideo.mp4 3 5
```

Equivalent to: `ffmpeg -ss 3 -i myvideo.mp4 -t <duration-8> -c copy myvideo_trim.mp4`

#### trim-start
Trims the specified number of seconds from the start of the video. Uses stream copying for speed.

**Syntax:** `ezff trim-start <input> <seconds> [--output <output>]`

- `<input>`: Path to input video file.
- `<seconds>`: Float seconds to trim from the start.
- `--output` (optional): Same as above.

**Example:** (removes 2 seconds from start)

```
ezff trim-start myvideo.mp4 2 --output=mynewvideo.mp4
```

Equivalent to: `ffmpeg -ss 2 -i myvideo.mp4 -c copy mynewvideo.mp4`

#### trim-end
Trims the specified number of seconds from the end of the video. Calculates duration using `ffprobe`.

**Syntax:** `ezff trim-end <input> <seconds> [--output <output>]`

- `<input>`: Path to input video file.
- `<seconds>`: Float seconds to trim from the end.
- `--output` (optional): Same as above.

**Example:** (removes 3 seconds from end)

```
ezff trim-end myvid.mp4 3
```

Equivalent to: `ffmpeg -i myvid.mp4 -c copy -t <duration-3> myvid_trim.mp4`

#### trim-mid
Removes a section from the middle of the video (keeps start and end parts).

**Syntax:** `ezff trim-mid <input> <start_cut> <end_cut> [--output <output>]`

- `<input>`: Path to input video file.
- `<start_cut>`: Start time (seconds) of section to remove.
- `<end_cut>`: End time (seconds) of section to remove (must be > start_cut).
- `--output` (optional): Same as above.

**Example:** (removes section starting at 5 second mark and ending at 10 second mark)

```
ezff trim-mid input.mp4 5 10 --output=new.mp4
```

Equivalent to: `ffmpeg -i input.mp4 -vf "select='not(between(t,5,10))',setpts=N/FRAME_RATE/TB" -af "aselect='not(between(t,5,10))',asetpts=N/SR/TB" new.mp4`

## Notes

- All commands require FFmpeg (and ffprobe for `trim` and `trim-end`). If not found, a warning is printed, and commands will fail.
- Output files are placed in the same directory as the input by default.
- Supports common video formats (MP4, etc.) via FFmpeg.
- Error handling: Invalid args or missing files will print to stderr and exit with code 1.

## License

MIT License. See [LICENSE](LICENSE) for details.

## Contributing

Feel free to open [issues](https://github.com/jantcu/ezff/issues) or [PRs](https://github.com/jantcu/ezff/pulls) for bugs, new commands, or improvements. Keep it minimal!
