# Testing Locally

  Sometimes we wish to test a local version of Workflow Manager API, rather than hitting our staging or prod environments, potentially also using a certain version of the [Workflow Manager client](https://github.com/deis/workflow-manager).  

  Here we provide the dependencies and steps needed to do so.

## Dependencies

  - [psql](https://www.postgresql.org/docs/9.2/static/app-psql.html)
  - a running [Deis Workflow](https://github.com/deis/workflow/blob/master/src/installing-workflow/index.md) cluster...
  - installed on a [Kubernetes](http://kubernetes.io/) cluster

## Steps

  1. TODO: outline postgres setup steps here


  2. Once the instance status is `Available`, we can seed `dbname` with test data:

    ```console
    psql \
      -f /path/to/test_data.sql \
      --host <postgres endpoint> \
      --port 5432 \
      --username dbuser \
      --dbname dbname
    ```

  3. On the installed Deis Workflow cluster we will launch a local version of Workflow Manager API. Here we refer to the routable ip for reaching the controller as `ROUTABLE_IP`, which can be the internal router service IP, externally accessible load balancer IP or node IP if using node port:

    ```console
    export
    export DEIS_CONTROLLER_URL=http://deis.${ROUTABLE_IP}.nip.io
    deis auth:register $DEIS_CONTROLLER_URL
    deis apps:create --no-remote wfm-api
    deis config:set \
      WORKFLOW_MANAGER_API_DBUSER=“${DBUSER}" \
      WORKFLOW_MANAGER_API_DBPASS=“${DBPASS}" \
      WORKFLOW_MANAGER_API_PORT=8081 \
      -a wfm-api
    deis pull quay.io/deisci/workflow-manager-api:canary -a wfm-api
    # optionally, specify org/workflow-manager-api:tag
    # to test a different wfm-api version, provided the image is
    # publicly accessible
    ```

    Let's verify that our wfm-api app is healthy. The following should return the current cluster count, depending on `test_data.sql` provided:

    `curl http://wfm-api.${ROUTABLE_IP}.nip.io/v3/clusters/count`

  4. Update the existing `deis-workflow-manager` pod to point to our local wfm-api app:

    ```console
    kubectl edit rc deis-workflow-manager —namespace=deis
    # update VERSIONS_API_URL to point to local wfm-api endpoint:
    # http://wfm-api.${ROUTABLE_IP}.nip.io
    # optionally, can also specify a different workflow-manager
    # image, provided it is publicly accessible
    kubectl scale rc deis-workflow-manager —namespace=deis —replicas=0
    kubectl scale rc deis-workflow-manager —namespace=deis —replicas=1
  ```

    The following should now return the previous cluster count incremented by 1 thanks to our newly reporting `deis-workflow-manager` pod:

    `curl http://wfm-api.${ROUTABLE_IP}.nip.io/v3/clusters/count`
