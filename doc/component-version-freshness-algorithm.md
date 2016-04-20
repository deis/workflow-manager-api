# Version Freshness Strategy

This document will outline how we determine if a deis cluster component is up-to-date with the latest release for a particular release "train" (e.g., "beta", "stable").

## Release Description

A release is attached to a particular deis component (e.g., "deis-router"). A release has the following properties:

- train
- version string
- release timestamp
- description and generic (not according to a strict schema) metadata

## Ideal Version Freshness Algorithm (TBD)

An ideal implementation would follow Semantic Versioning 2.0.0, which is outlined here:

- http://semver.org

The realities of the deis v2 release process (we want to get to a stable, supportable 2.0 platform release ASAP) preclude a strict adherence to sever. In the meantime, however, we do want to have a best-effort strategy to compare releases for version freshness.

## Best-effort Version Freshness Algorithm

Using a release's "train" value as an isolating property, we can simply compare among common same-train releases the release timestamp to determine the most recent release. Some examples:

- "deis-router" component, "beta" train, "2.0.0-rc1" version, "2016-03-30T23:54:39Z" release timestamp
- "deis-router" component, "beta" train, "2.0.0-rc2" version, "2016-03-31T23:54:39Z" release timestamp
  - The 2nd "2.0.0-rc2" version will be judged as the most recent release

A summary to the above approach could be stated like this: it's a best-effort approach that relies on a strictly comparable data type (the release timestamp) to do the programmatic version recency comparison, and that uses a simple human-readable version string to suggest version precedence to the user community (and that can be incorporated into a CI/CD-driven auto-incrementation strategy); most importantly, it enables us to refine the version string over time based on feedback from the community without having to implement specific, complicated business logic (e.g, a full semver type comparison library) that might need to be reimplemented in the future.

## Gotchas

We'll rely on our CI/CD processes to tag releases with a release string and release timestamp. Therefore, within those processes (both human and programmatic) we'll need to do our best to ensure that we rationally increment the version string along with the release timestamp. For example, we won't have the luxury of fixing/cleaning up a prior release, retaining the prior release string while updating the release timestamp to reflect the fix time. We could anticipate such scenarios, and include multiple timestamp fields to accommodate, but this adds complexity to the best-effort strategy that is more appropriately suited for the eventual, "ideal" strategy (e.g., semver).
