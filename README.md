# FontsourceDownloader

This tool allows you to download all fonts from [Fountsource](https://fontsource.org) to your desired location using [jsdelivr](https://www.jsdelivr.com/) to download them.

## Installation

```bash
go install github.com/zergo0/fontsourcedownloader/@latest
```

## Usage

```bash
fontsourcedownloader -out /path/to/output
```

Possible arguments:

- `-out` - output path (ex: `/path/to/output`)
- `-formats` - comma separated list of formats to download (ex: `woff2,woff`)
- `-weights` - comma separated list of weights to download (ex: `400,700`)
- `-styles` - comma separated list of styles to download (ex: `normal,italic`)
- `-subsets` - comma separated list of subsets to download (ex: `latin,latin-ext`)

## Development

The project uses `make` to make your life easier. If you're not familiar with Makefiles you can take a look at [this quickstart guide](https://makefiletutorial.com).

Whenever you need help regarding the available actions, just use the following command.

```bash
make help
```

### Setup

To get your setup up and running the only thing you have to do is

```bash
make all
```

This will initialize a git repo, download the dependencies in the latest versions and install all needed tools.
If needed code generation will be triggered in this target as well.

### Test & lint

Run linting

```bash
make lint
```

Run tests

```bash
make test
```

Made with [![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/SchwarzIT/go-template)
