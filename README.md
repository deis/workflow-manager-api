# Deis Workflow Manager Service API

The Workflow Manager Service API is responsible for interfacing with Deis Workflow data. It is a golang https implementation that speaks JSON. The API is the source of authority for the following:

* a list of all Deis clusters whose Workflow Manager clients have chosen to share anonymous usage data
* metadata relevant to each unique Deis cluster: which components are installed, and which versions
* various statistics on each cluster, e.g.:
  * count of clusters "currently" running
  * count of clusters "recently" disappeared
  * count of clusters "recently" added
  * average age of clusters
* "latest" stable version for each Deis cluster component

Additionally, the API is the official interface for accepting Workflow Manager data, e.g.:

* Deis cluster component version CRUD operations
  * i.e., centrally store latest stable version information for all Deis cluster components
* Deis unique cluster anonymous registration
  * i.e., receive Deis Workflow Manager client usage statistics

# Usage

To test:
```
$ go test
PASS
ok  	github.com/deis/workflow-manager-api	0.012s
```
To build:
```
$ go build -o workflow-manager-api *.go
```
To run:
```
$ ./workflow-manager-api
```

# Status

The API is currently for demonstration only. Included in the codebase are a private key and self-signed certificate; these should be considered disposable at this time and for testing/demo only.

## License

Copyright 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
