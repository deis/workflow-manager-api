
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Deis Workflow is no longer maintained.<br />Please [read the announcement](https://deis.com/blog/2017/deis-workflow-final-release/) for more detail. |
|---:|---|
| 09/07/2017 | Deis Workflow [v2.18][] final release before entering maintenance mode |
| 03/01/2018 | End of Workflow maintenance: critical patches no longer merged |
| | [Hephy](https://github.com/teamhephy/workflow) is a fork of Workflow that is actively developed and accepts code contributions. |

# Deis Workflow Manager Service API

The Workflow Manager Service API is responsible for interfacing with Deis Workflow data. It is a golang https implementation that speaks JSON. The API is the source of authority for the following:

* metadata about unique Deis clusters that have checked in: which components are installed at the time of check-in, and which versions
* released version history for each Deis cluster component

Additionally, the API is the official interface for accepting Workflow Manager data, e.g.:

* Deis cluster component version CRUD operations
  * i.e., centrally store latest stable version information for all Deis cluster components
* Deis unique cluster anonymous registration
  * i.e., receive Deis Workflow Manager client usage statistics

# Usage

To download dependencies:
```
$ make bootstrap
```
To test:
```
$ make test
```
_Note_: if you prefer to run tests with `go test`, you'll need
[`glide`](https://github.com/Masterminds/glide) on your `PATH`. Once you have it,
run `go test -tags testonly $(glide nv)` to run tests.
To build:
```
$ IMAGE_PREFIX=$MY_DOCKERHUB_ACCOUNT make build docker-build docker-push
```
(All of the above operations assume a local Docker environment.)

# Status

A working, minimal API is currently live at https://versions.deis.com.
