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

## License

Copyright 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
