### v2.0.0-beta3 -> v2.0.0-beta4

#### Features

 - [`8df8999`](https://github.com/deis/workflow-manager-api/commit/8df89991d9f4ef16e2c762a79d2962e373dafef1) _scripts/deploy.sh,Makefile: add auto-deployment to the staging deis cluster
 - [`f205e13`](https://github.com/deis/workflow-manager-api/commit/f205e13813fe0649d1f8de0234b31d71c3f69941) server.go,handlers/clusters_age_handler.go: design the cluster age REST API interface

#### Fixes

 - [`f5daff4`](https://github.com/deis/workflow-manager-api/commit/f5daff46fa9688b4a519e83cd8c92109496ebd25) data: fix the cluster age filter

#### Documentation

 - [`a2f321d`](https://github.com/deis/workflow-manager-api/commit/a2f321d0170fd8a6eb954debd99b2ed275d5db93) CHANGELOG.md: update for v2.0.0-beta3

#### Maintenance

 - [`e5757f4`](https://github.com/deis/workflow-manager-api/commit/e5757f426e913f9905f51e83aefc525ed9153774) glide: remove obsolete deps

### v2.0.0-beta3

#### Features

 - [`9deebbf`](https://github.com/deis/workflow-manager-api/commit/9deebbf37b77a9d46bd532b5e76ca5a35170a59d) server.go,handlers: create multi-get-endpoint to get all latest releases
 - [`05ac9c4`](https://github.com/deis/workflow-manager-api/commit/05ac9c42e43b89dbb391e82b5f7cb6435ab61d46) handlers,server.go: add endpoint to get latest release for specified component/train
 - [`777688e`](https://github.com/deis/workflow-manager-api/commit/777688e61a471721200ca96adb7fdc6717221f5b) data: versions handlers + data scaffolding
 - [`11d28fb`](https://github.com/deis/workflow-manager-api/commit/11d28fbf18fa55d9c0b199cdc5f2b0c7a504d248) Dockerfile: added root certs to docker image
 - [`3352d7a`](https://github.com/deis/workflow-manager-api/commit/3352d7a539947ce5af2e667093ffb98349f0b807) deploy: add script to build and deploy images from master
 - [`f0ac5c0`](https://github.com/deis/workflow-manager-api/commit/f0ac5c07b1896fd0e7cc12f43686349ed5777929) CI: welcome travis

#### Fixes

 - [`a93301f`](https://github.com/deis/workflow-manager-api/commit/a93301f90041b4bd7dbccc426b841b0f60d296ea) rootfs: copy only the binary into the image
 - [`447ffd7`](https://github.com/deis/workflow-manager-api/commit/447ffd7836ad5969da96831f95a3a5aec9900021) data: updateClusterDBRecord now has valid SQL
 - [`dd06933`](https://github.com/deis/workflow-manager-api/commit/dd06933ce8789527e3868e4d7a8a48b053e077de) data: removed uniqueness constraints in clusters_checkins table
 - [`85444af`](https://github.com/deis/workflow-manager-api/commit/85444af1d0f447ac65bb47b6827867de37bdf073) handlers/handlers.go: add format string for error printing
 - [`fcdd316`](https://github.com/deis/workflow-manager-api/commit/fcdd316d2d23e0b57e6def5567f5dec69f9df051) data: re-using a single db connection
 - [`e6dde58`](https://github.com/deis/workflow-manager-api/commit/e6dde58df6bef17812e671f63831a6ce93325085) handlers: rationalizing error logging
 - [`32e1413`](https://github.com/deis/workflow-manager-api/commit/32e1413614cfffb8cef6205a03bfc63faa44fbb7) timestamp: time.Time type
 - [`c96115d`](https://github.com/deis/workflow-manager-api/commit/c96115d26cd6ff93ed845462e1c821773679fc1f) debugging: we don't want to obfuscate data errors!
 - [`5dfd183`](https://github.com/deis/workflow-manager-api/commit/5dfd183f8fd40b9e13b6c2af0ea60fc6e90cdec7) _scripts/deploy.sh: make deploy script executable
 - [`1dd4751`](https://github.com/deis/workflow-manager-api/commit/1dd4751e42a452a9df061980e11afb82f6484bdf) data: add tests for the data package
 - [`2af5a7f`](https://github.com/deis/workflow-manager-api/commit/2af5a7fb2540c38b7329ca1334692f2e4a9222b3) handlers: add tests for handlers util functions
 - [`d057ff9`](https://github.com/deis/workflow-manager-api/commit/d057ff935d37e3359445bfce2f87d32c0d6bc171) server_test.go: finish server tests
 - [`bfa7016`](https://github.com/deis/workflow-manager-api/commit/bfa7016b13220825083a8282bc652b0b4eeed587) handlers_test.go: add unit tests for handlers
 - [`8e6cfae`](https://github.com/deis/workflow-manager-api/commit/8e6cfaec4b7f11ef90d41915688059f893703d63) ci: travis needs glide deps which make bootstrap provides
 - [`eba71fb`](https://github.com/deis/workflow-manager-api/commit/eba71fb76b5ff263102e58ed715dad5db33d76f7) data: replace non-existent ParseJSONComponent func with internal implementation

#### Documentation

 - [`becd067`](https://github.com/deis/workflow-manager-api/commit/becd0679165d1e60f51ea3d68e683a3b6fca8bdb) data: basic overview of data implementation

#### Maintenance

 - [`c2f2975`](https://github.com/deis/workflow-manager-api/commit/c2f297586bfc968d3d581198a7081014ae0de87b) glide: manage dependencies with glide
