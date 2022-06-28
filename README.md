<div align="center"><img src="podium.png" alt="Flat art of a sports podium."></div>
<div align="center"><small><i><a href="https://www.flaticon.com/free-icons/podium" title="podium icons">Podium icons created by Roundicons - Flaticon</a></i></small></div>
<h1 align="center">
  <b><i>ranking-cli</i></b>
</h1>

<h4 align="center">A Go example CLI modelling a league ranking app.</h4>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#license">License</a>
</p>

<p align="center">
  <a href="https://github.com/liampulles/ranking-cli/releases">
    <img src="https://img.shields.io/github/v/release/liampulles/ranking-cli">
  </a>
  <a href="https://app.travis-ci.com/github/liampulles/ranking-cli">
    <img src="https://app.travis-ci.com/liampulles/ranking-cli.svg?branch=main" alt="[Build Status]">
  </a>
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/liampulles/ranking-cli">
  <a href="https://goreportcard.com/report/github.com/liampulles/ranking-cli">
    <img src="https://goreportcard.com/badge/github.com/liampulles/ranking-cli" alt="[Go Report Card]">
  </a>
  <a href="https://codecov.io/gh/liampulles/ranking-cli" > 
    <img src="https://codecov.io/gh/liampulles/ranking-cli/branch/main/graph/badge.svg?token=RU6ycM2b3J"/> 
  </a>
  <a href="https://github.com/liampulles/ranking-cli/blob/master/LICENSE.md">
    <img src="https://img.shields.io/github/license/liampulles/ranking-cli.svg" alt="[License]">
  </a>
</p>

## Installation

### From release

Simply download a release from the releases page and run! You may wish to put it on your PATH for it to be easily executable (e.g. in `/usr/local/bin`).

### Via `go install`

Run:

```shell
go install github.com/liampulles/ranking-cli/cmd/sportrank@latest
```

To get the latest version on master, otherwise change `latest` to a specific version.

### From source

To install from source (requires `make`) clone the repo and in the root folder run

```shell
make install
```

This will build the app and place it in your configured GOBIN directory.

## Usage

Given a file `input.txt` with contents:

```
Lions 3, Snakes 3
Tarantulas 1, FC Awesome 0
Lions 1, FC Awesome 1
Tarantulas 3, Snakes 1
Lions 4, Grouches 0
```

Running the following:

```shell
sportrank -i input.txt
```

Will produce output:

```
1. Tarantulas, 6 pts
2. Lions, 5 pts
3. FC Awesome, 1 pt
3. Snakes, 1 pt
5. Grouches, 0 pts
```

You can also choose to take input from STDIN and/or output to a file, e.g.:

```
cat example.txt | sportrank -o output.txt
```

## License

See [LICENSE](LICENSE)