# Deis Component Versions Publishing

The Workflow Manager Service API is the authoritative source of version history for all Deis component releases. It is responsible for receiving and for processing new version data. This document outlines that process in detail, as well as provides strategic guidance for future "publishers" who will be responsible for delivering new version data to the API.

# How to publish a new version to the API

The API's "publish version" business logic responds to this `HTTP POST` URL route endpoint:

* POST `/{apiVersion}/versions/{train}/{component}/{version}`

It accepts as request body data with media type "application/json" a JSON representation of the `ComponentVersion` type defined here:

* https://github.com/deis/workflow-manager/blob/master/types/types.go

Here is an example against an instance of the API in the real world:

```
curl -H "Content-Type: application/json" -X POST -d \
'{"component": {"name": "deis-builder"}, "version": {"train": "beta", "version": "2.0.0-beta2", "released": "2016-04-16T23:54:39Z07:00"}}' https://versions.deis.com/v2/versions/beta/deis-builder/2.0.0-beta2
```

Note that the release time at the JSON data point above `version.released` does not include a time zone. This is because the `release_timestamp` column is defined as type `timestamp without time zone` in the postgres data `versions` table. Including the time zone, e.g., "`2016-04-16T23:54:39Z07:00`" won't do any damage: it will be ignored when it's stored as a record in the database table. Also note the general-purpose `data` JSON object: this maps to a `json` column type in the `versions` table. Include keys and values related to the release here, e.g., bugfix info, high level descriptions, links to issues, etc. The below JSON-ish representation will help to describe the data that we'll be including with every release version publishing event:

```
{
  "component": { // component object, required
    "name": "deis-builder" // component name, required
  },
  "version": { // version object, required
    "train": "beta", // version train for this release, required
    "version": "2.0.0-beta2", // version string for this release, required
    "released": "2016-04-16T23:54:39", // release time for this release, required
    "data": { // data object, required
      "notes": "here's beta2!" // example of data property, all properties are optional
    }
  }
}
```

In summary, all properties in the JSON representation of a `ComponentVersion` above are required, with the exception of the `version.data` payload, which may be `{}` (i.e., an empty object). Additionally, the HTTPS destination URL carries properties that map to properties in the request body. For maximum app visibility (in terms of HTTP log history), the API requires this redundant data. In the URL "`https://versions.deis.com/v2/versions/beta/deis-builder/2.0.0-beta2`" above, the final three dynamic segments map to the following data properties in the request body:

* "`beta`" maps to the `version.train` key's value in the request body
* "`deis-builder`" maps to the `component.name` key's value in the request body
* "`2.0.0-beta2`" maps to the `version.version` key's value in the request body

The dynamic segment strings in the URL route must match their corresponding values as explained above, otherwise the API will complain with a `400` response.

# For near-future strategic consideration: Automating all this!

There are five bits of data that any CI process will need to be able to programmatically determine in order to automate release version publishing. These are listed below, along with some possible CI-accessible sources of authority that can help to accomplish data determination:

* The component name
  * This should be inferred from the repository name. The name passed to the versions API can be matched to the name injected into the component's chart.
* The version train
  * This can be inferred from git metadata: branch and tag, possibly. ":master" branch transitions can map to a "stable" or "production" train (definitive name TBD); we need to standardize our git branch metadata decoration processes in order to successfully map non-master branch transitions to component trains.
* The version string
  * We should be guided by the existing release tag processes, and then adapt non-master branch transitions (i.e., non-"stable" train versions) to them
* The release timestamp
  * This can map generally to the time of the git transition; caveat emptor: we need to be cognizant of golang-interpretable timestamp strings, e.g., the "RFC3339" string format as in the example above. See https://golang.org/pkg/time/#pkg-constants for more information.
* General purpose release data
  * TBD. Once there is a solid release notes strategy in github/dockerhub/quay.io, we should have CI consume that information and inject it into the "`data`" payload as specified above.
