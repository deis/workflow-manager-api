As we try to grow and scale the Deis business, we are constantly in the dark as to how many clusters are actually running. “Deis Workflow Manager” is aimed at helping customers who have deployed workflow in exchange for some anonymous usage statistics.

# Goals for the project

- Respect customers and their data
  - No identifying information should be delivered. Including hostnames, ip addresses, customer application names, or git repositories.
- Gather anonymous usage data so that we can know
  - How many “clusters” are running
  - What versions of the components are running
- Ability to “push” information about updates to customers visiting the manager UI
- Serve as a quick and easy admin UI, for basic workflow information (versions, basic status: there, not there)

# As a Workflow Manager User I should

- Be able to view a simple web UI that I can reach through the kubernetes proxy
  - I should not have to authenticate to this UI
- Be able to view a list of services, replication controllers and pods that are installed or running in the Deis namespace
  - So that I can quickly, at a glance, know if all the components are installed
- Be able to see the versions of installed workflow components, and be notified if they are out of date
  - Annotations, container images, (whatever makes sense here)
- See a large banner if deis workflow components are out of date
  - Poll a hosted releases API off-cluster for release information at least every 12 hours
- Should be able to disable, via pod environment variables, to:
  - Disable update checks (default: on)
  - Enable anonymous data collection (default: OFF)
- **Be able to opt-in to sending anonymous & anonymized usage data:**
  - `k describe node <node>`
    - total cpu / memory / disk
    - docker versions, os image, kubelet versions
    - creation timestamp, etc.
  - Enumeration of deis workflow components and versions

# As Deis, the company I should be able to

- Uniquely, but anonymously, identify a customer cluster
- Query by installed version
- Get a count of workflow instances currently running (which have been seen in the last 24 hours)
- Get a count of workflow instances that have disappeared (week, month, quarter)
- Get a count of workflow instances that have appeared (week, month, quarter)
- Get the average age of workflow instances (time since cluster was first created)
- Uniquely and anonymously identify an instance of workflow
- Exclude staff, ci and test clusters from the stats

# Architectural Overview

The Workflow Manager systems can be broken down into 5 distinct pieces:

- **A deis cluster “peer module”*** that has visibility into other deis components (router, controller, logger, etc). This module will:
  - Deploy onto a customer’s kubernetes cluster, as part of deis workflow installation
  - Collect and organize deis cluster behavior, activity and  “state”, and deliver that organized data to consumers (CLI, UI)
    - e.g., a simple dictionary of deis components and their currently running versions
    - e.g., a list of components whose versions are out-of-date
  - Have access to the public internet, where it will securely send anonymous deis cluster usage data
    - e.g., deis components and their versions running in this cluster
- **A deis cluster admin web UI*** that exposes the organized data from the above module
  - runs on customer cluster
- **A publically accessible service API** that accepts anonymized deis cluster data, and acts as a deis cluster’s source of authority for the deis project.
  - hosted by Deis team
  - e.g., returns an authoritative “latest” version for all deis components
- **A private key-value data store** to keep persistent data. The service API in the above bullet point will be its consumer.
  - run/hosted by Deis team
  - e.g., stores authoritative deis version data
  - e.g., stores metadata about deis clusters in the wild that have consented to sharing usage data
- **A deis workflow manager web UI** whose primary purpose is presenting anonymous deis cluster usage data. A lightweight analytics dashboard.
  - run/hosted by Deis team

# Deis Workflow Manager Module

This is what’s referred to as a deis cluster “_peer module_” above. My thinking here is that this will be, in shorthand, a “go package” that will slot in conveniently alongside other deis cluster components (router, controller, etc), with access to usage and other metadata about the cluster (I assume via etcd), and whose primary purpose is to organize that metadata into consumable, organized collections. We should assume that data organization is both on-demand (a response to a request) and periodic/cached (let’s call these “reports”).

This module will listen on HTTP, and its API endpoints will be organized in a RESTFul-ish way.

# Deis Workflow Manager Service API

This should be a very simple RESTful HTTP service that sits on top of a persistent key-value data store. It will be primarily a RESTful convenience layer that maps to an underlying, predictable data schema. I don’t see any data schemas being enforced at this layer, just this kind of thing:

- `data = { versions: { thing1: [1, 2, 3], thing2: [1, 2] } }`
- HTTP route = `/versions/:thing`
- HTTP response = <a JSON array>


# Deis Cluster admin web UI

We want to put something here that will be generic enough to accommodate future “admin web UI” functionality (an alternative to the CLI). In MVP fashion, though, we want to start with a simple read-only, unauthenticated UI that is essentially a part of the deis cluster, and is accessible over HTTP via the kubernetes proxy.

# Deis Workflow manager web UI

We want to build something here that is dead-simple, but is future-friendly to being extended into a kind of general purpose deis analytics dashboard.
