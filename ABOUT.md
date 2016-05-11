As we try to grow and scale the Deis business, we are constantly in the dark as to how many clusters are actually running. “Deis Workflow Manager” is aimed at helping customers who have deployed workflow in exchange for some anonymous usage statistics.

# Goals for the project

- Respect customers and their data
  - No identifying information should be delivered. Including hostnames, ip addresses, customer application names, or git repositories.
- Gather anonymous usage data so that we can know
  - How many “clusters” are running
  - What versions of the components are running
- Ability to “push” information about updates to customers visiting the manager UI
- Serve as a quick and easy admin UI, for basic workflow information (versions, basic status: there, not there)

# Operators

Operators of a Deis cluster should be able to do the following using Workflow Manager:

- View a simple web UI, reachable through the Deis router or similar proxy, without authentication
- View a list of services, replication controllers and pods that are installed or running in the Deis namespace
  - The list is intended to provide a quick, at-a-glance view of all components so that operators can determine if all necessary components are installed
- See the versions of installed workflow components, and be notified if they are out of date
  - Annotations, container images, (whatever makes sense here)
- See a large banner if deis workflow components are out of date
  - Poll a hosted releases API off-cluster for release information at least every 12 hours
- Toggle the following behavior via pod environment variables:
  - Update checks (default: on)
  - Anonymous data collection (default: off)
- **Opt-in to sending anonymous and anonymized usage data:**
  - `k describe node <node>`
    - total cpu / memory / disk
    - docker versions, os image, kubelet versions
    - creation timestamp, etc.
  - Enumeration of deis workflow components and versions

# Selected Deis Employees

Administrators, product managers and executives at Deis should be able to do the following using the Workflow Manager API server:

- Uniquely, but anonymously, identify a customer cluster
- Query by installed version
- Get a count of workflow instances currently running (which have been seen in the last 24 hours)
- Get a count of workflow instances that have disappeared (week, month, quarter)
- Get a count of workflow instances that have appeared (week, month, quarter)
- Get the average age of workflow instances (time since cluster was first created)
- Uniquely and anonymously identify an instance of workflow
- Select who else can see the above information (i.e. exclude staff, ci and test clusters)

# Architectural Overview

The Workflow Manager systems can be broken down into 5 distinct pieces:

- **A deis cluster “peer module”*** that has visibility into other deis components (router, controller, logger, etc). This module will:
  - Deployed onto a customer’s kubernetes cluster, as part of deis workflow installation
  - Collect and organize deis cluster behavior, activity and  “state”, and deliver that organized data to consumers (CLI, UI)
    - e.g., a simple dictionary of deis components and their currently running versions
    - e.g., a list of components whose versions are out-of-date
  - Have access to the public internet, where it will securely send anonymous deis cluster usage data
    - e.g., deis components and their versions running in this cluster
- **A deis cluster admin web UI*** that exposes the organized data from the above module
  - Runs on customer cluster
- **A publically accessible service API** that accepts anonymized deis cluster data, and acts as a deis cluster’s source of authority for the deis project.
  - Hosted by Deis team
  - e.g., returns an authoritative “latest” version for all deis components
- **A private key-value data store** to keep persistent data. The service API in the above bullet point will be its consumer.
  - run/hosted by Deis team
  - e.g., stores authoritative deis version data
  - e.g., stores metadata about deis clusters in the wild that have consented to sharing usage data
- **A deis workflow manager web UI** whose primary purpose is presenting anonymous deis cluster usage data. A lightweight analytics dashboard.
  - Run/hosted by Deis team

# Deis Workflow Manager Module

This is what’s referred to as a deis cluster “_peer module_” above. My thinking here is that this will be, in shorthand, a “go package” that will slot in conveniently alongside other deis cluster components (router, controller, etc), with access to usage and other metadata about the cluster (I assume via etcd), and whose primary purpose is to organize that metadata into consumable, organized collections. We should assume that data organization is both on-demand (a response to a request) and periodic/cached (let’s call these “reports”).

This module will listen on HTTP, and its API endpoints will be organized in a RESTFul way.

# Deis Workflow Manager Service API

This is a RESTful HTTP service that sits on top of a persistent key-value data store. It will be primarily a RESTful convenience layer that maps to an underlying, predictable data schema. The following are some examples of the schema and API:

- `data = { versions: { thing1: [1, 2, 3], thing2: [1, 2] } }`
- HTTP route = `/versions/:thing`
- HTTP response = <a JSON array>


# Deis Cluster admin web UI

The Workflow Manager will provide a generic “admin web UI” (an alternative to the CLI). The first version (MVP) will provide a simple, read-only, unauthenticated UI as part of the default Deis cluster install. It is accessible over HTTP via the Deis router or similar proxy.

# Deis Workflow Manager API web UI

The Workflow Manager API will provide a web UI that provides read-only access to the data that the Workflow Manager API stores. This dashboard will hold visual representations (such as graphs) of the analytics data stored in the data.
