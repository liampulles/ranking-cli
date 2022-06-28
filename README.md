<div align="center"><img src="podium.png" alt="Flat art of a sports podium."></div>
<div align="center"><small><i><a href="https://www.flaticon.com/free-icons/podium" title="podium icons">Podium icons created by Roundicons - Flaticon</a></i></small></div>
<h1 align="center">
  <b><i>ranking-cli</i></b>
</h1>

<h4 align="center">A Go example CLI modelling a league ranking app.</h4>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#notes">Notes</a> •
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

## Notes

### Architecture

The architecture of the system follows Robert Martin's clean architecture (https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This means the usecase layer has no sight of the adapter layer, and the adapter layer has no sight of the driver layer. The architecture differs slightly though in that I've added an outer `wire` layer for dependency injection, and that I've put the domain logic in `pkg/league`, since it may be useful as a library for other apps.

Using this architecture, one could trivially extend the app (or create a new app) to invoke the ranking code from an HTTP endpoint, or RPC method, etc. without modifying the business logic itself (in `pkg/league`).

I wrote a little piece a while back around using the Clean architecture in Go on my site: https://liampulles.com/2020/09/29/notes-on-applying-the-clean-architecture-in-go.html

### The wire package

As mentioned, the wire package deals with dependency injection. One could use the Google wire tool to do this, but for a simple app like this I think its simple enough to wire it oneself.

### Service pattern

There are a handful of "services" in the system, for example `cli.Engine`. These services consist of an interface and an implementation and may depend on other service interfaces, which are injected by the wire package.

This pattern is very useful for unit testing, since dependent services can be mocked (in this repo, the mocks are generated with mockery) to ensure that the service being tested is the only thing being tested.

Note lines like this:
```go
var _ RowIOGateway = &RowIOGatewayImpl{}
```
This line will fail if *RowIOGatewayImpl does not implement RowIOGateway - this line thus acts as a kind of guard (which is why I use it).

## Makefile

I'm familiar with other sort of local repository management tools (e.g. Gradle), but I like `make` for my personal projects because it is simple to use, powerful, and widely available on Linux-like systems.

Some things you might like to try:

* `make view-cover`: Generate and view test coverage
* `make pre-commit`: Update dependencies, run tests, generate code, lint, etc. - I run this before I submit code generally.
* `make install-tools`: Install any tools needed to work on the project. This should get invoked automatically if you run a make command and a needed tool is not available (I hope).

# CI/CD

For personal projects, I use Travis for CI, and goreleaser for CD (in this context, meaning creating a Github release with binaries).

## License

See [LICENSE](LICENSE)