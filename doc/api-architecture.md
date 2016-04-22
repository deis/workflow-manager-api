# Workflow Manager API Architecture

This is a living document that reflects architectural intentions. Where implementation is TBD, it will be noted.

## Overview

The Workflow Manager API is intended as a lightweight HTTP interface in front of a persistent back-end data store. While not strictly enforced, we will strive toward a close coupling between API endpoints and common data queries. It will not generally be the responsibility of the API to meaningfully transform, interpret, or obfuscate the underlying data representation.

Moreover, we will assume a more proximate relationship between the API and persistent data, compared to API consumers and the API itself; as such, the API *will*, when appropriate, optimize end-to-end behavior. For example, "Multi Get" API endpoints will favor single network calls from the consumer that represent requests requiring multiple queries to persistent storage.

## Terminology

- A *component* is a single deis component, e.g., `deis-router`
- A *train* is a release cadence type, e.g., "beta" or "stable"
- A *version* is a versioned string attached to a component, e.g., "2.0.0" or "v2-beta"

## API endpoints

- Get a particular release
  - Request type and URL
    - `GET /:apiVersion/versions/:train/:component/:release` *implemented in "versions-train" as of 4/18*
  - 200 Response Body
  ```
  {
    "component": {
      "name": "deis-builder",
      "description": "Deis Builder"
    },
    "version": {
      "train": "stable",
      "version": "2.0.2",
      "released": "2016-03-31T23:54:39Z"
      "data": {
        "description": "release notes here",
        "fixes": "list of bug fixes"
      }
    }
  }
  ```
- Get the set of released component + train + versions
  - Request type and URL
    - `GET /:apiVersion/versions/:train/:component` *implemented in "versions-train" as of 4/18*
  - 200 Response Body
  ```
  {
    "data": [
    {
      "component": {
        "name": "deis-builder",
        "description": "Deis Builder"
      },
      "version": {
        "train": "stable",
        "version": "2.0.2",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-builder",
        "description": "Deis Builder"
      },
      "version": {
        "train": "stable",
        "version": "2.0.1",
        "released": "2016-03-21T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-builder",
        "description": "Deis Builder"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-11T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    }
    ]
  }
  ```
- Get the latest component + train release
  - Request type and URL
    - `GET /:apiVersion/versions/:train/:component/latest`
  - 200 Response Body
  ```
  {
    "component": {
      "name": "deis-builder",
    },
    "version": {
      "train": "stable",
      "version": "2.0.2",
      "released": "2016-03-31T23:54:39Z"
      "data": {
        "description": "release notes here",
        "fixes": "list of bug fixes"
      }
    }
  }
  ```
- Multi Get a collection of latest releases
  - Request type and URL
    - `POST /:apiVersion/versions/latest`
  - Request body
  ```
  {
    "data": [
    {
      "component": {
        "name": "deis-builder"
      },
      "version": {
        "train": "stable"
      }
    },
    {
      "component": {
        "name": "deis-controller"
      },
      "version": {
        "train": "stable"
      }
    },
    {
      "component": {
        "name": "deis-database",
      },
      "version": {
        "train": "beta"
    },
    {
      "component": {
        "name": "deis-minio"
      },
      "version": {
        "train": "canary"
      }
    },
    {
      "component": {
        "name": "deis-registry"
      },
      "version": {
        "train": "stable"
      }
    },
    {
      "component": {
        "name": "deis-router"
      },
      "version": {
        "train": "stable"
      }
    },
    {
      "component": {
        "name": "deis-workflow-manager"
      },
      "version": {
        "train": "stable"
      }
    }
    ]
  }
  ```
  - 200 Response body
  ```
  {
    "data": [
    {
      "component": {
        "name": "deis-builder",
        "description": "Deis Builder"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-controller",
        "description": "Deis Controller"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-database",
        "description": "Deis Database"
      },
      "version": {
        "train": "beta",
        "version": "2.0.0-beta1",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-minio",
        "description": "Deis Minio"
      },
      "version": {
        "train": "canary",
        "version": "2.0.0-canary99",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-registry",
        "description": "Deis Registry"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-router",
        "description": "Deis Router"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    },
    {
      "component": {
        "name": "deis-workflow-manager",
        "description": "Deis Workflow Manager"
      },
      "version": {
        "train": "stable",
        "version": "2.0.0",
        "released": "2016-03-31T23:54:39Z"
        "data": {
          "description": "release notes here",
          "fixes": "list of bug fixes"
        }
      }
    }
    ]
  }
  ```
- Publish a new release
  - Request type and URL
    - `POST /:apiVersion/versions/:train/:component/:release` *implemented in "versions-train" as of 4/18*
  - Request body
  ```
  {
    "component": {
      "name": "deis-builder",
      "description": "Deis Builder"
    },
    "version": {
      "train": "stable",
      "version": "2.0.3",
      "released": "2016-04-11T23:54:39Z"
      "data": {
        "description": "release notes here",
        "fixes": "list of bug fixes"
      }
    }
  }
  ```
  - 200 Response body
  ```
  {
    "component": {
      "name": "deis-builder",
      "description": "Deis Builder"
    },
    "version": {
      "train": "stable",
      "version": "2.0.3",
      "released": "2016-04-11T23:54:39Z"
      "data": {
        "description": "release notes here",
        "fixes": "list of bug fixes"
      }
    }
  }
  ```
- Get a simple "known deis clusters" count
  - Request type and URL
    - `GET /:apiVersion/clusters/count`
  - 200 Response Body
  ```
  10231
  ```
- Filter clusters by age
  - Request type and URL
    - `GET /:apiVersion/clusters/age?checked_in_before=:timestamp&checked_in_after=:timestamp&created_before=:timestamp&created_after=:timestamp`
  - 200 Response Body

  ```
  {
    "data": [
      {
        "id": "8c6da034-c8b1-489a-a55d-a2215d93f934",
        "firstSeen": "2016-03-11T23:54:39Z",
        "lastSeen": "2016-03-31T23:54:39Z",
        "components": [
          {
            "component": {
              "name": "deis-builder",
              "description": "Deis Builder"
            },
            "version": {
              "train": "stable",
              "version": "2.0.0",
              "released": "2016-03-31T23:54:39Z"
            }
          },
          ...
        ]
      }
    ]
  }
  ```
  - Error conditions
	  - `400 Bad Request` if `created_before` > `checked_in_before`
		  - Because you can't have clusters that were checked in before they were created
	  - `400 Bad Request` if `checked_in_after` >= `checked_in_before`
		  - Because you can't have clusters that were checked in before time T-1 and at the same time checked in after time T+1 (or both at time T)
	  - `400 Bad Request` if `created_after` >= `created_before`
		  - Because you can't have clusters that were created after time T+1 and at the same time created before time T-1 (or both at time T)
	  - `400 Bad Request` if `checked_in_before` <= `created_after`
		  - Because you can't have clusters that were checked in before time T-1 and at the same time created after time T+1 (or both at time T)
- Get component metadata for a specific deis cluster
  - Request type and URL
    - `GET /:apiVersion/clusters/:id`
  - 200 Response Body
  ```
  {
    "firstSeen": "2016-03-11T23:54:39Z",
    "lastSeen": "2016-03-31T23:54:39Z",
    "id": "8c6da034-c8b1-489a-a55d-a2215d93f934"
    "components": [
      {
        "component": {
          "name": "deis-builder",
          "description": "Deis Builder"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-31T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-controller",
          "description": "Deis Controller"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-30T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-database",
          "description": "Deis Database"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-29T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-minio",
          "description": "Deis Minio"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-28T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-registry",
          "description": "Deis Registry"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-27T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-router",
          "description": "Deis Router"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-26T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-workflow-manager",
          "description": "Deis Workflow Manager"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-25T23:54:39Z"
        }
      }
    ]
  }
  ```
- Submit deis cluster component metadata
  - Request type and URL
    - `POST /:apiVersion/clusters`
  - Request body
  ```
  {
    "id": "8c6da034-c8b1-489a-a55d-a2215d93f934"
    "components": [
      {
        "component": {
          "name": "deis-builder",
          "description": "Deis Builder"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-31T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-controller",
          "description": "Deis Controller"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-30T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-database",
          "description": "Deis Database"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-29T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-minio",
          "description": "Deis Minio"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-28T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-registry",
          "description": "Deis Registry"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-27T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-router",
          "description": "Deis Router"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-26T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-workflow-manager",
          "description": "Deis Workflow Manager"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-25T23:54:39Z"
        }
      }
    ]
  }
  ```
  - 200 Response Body
  ```
  {
    "firstSeen": "2016-03-11T23:54:39Z",
    "lastSeen": "2016-03-31T23:54:39Z",
    "id": "8c6da034-c8b1-489a-a55d-a2215d93f934"
    "components": [
      {
        "component": {
          "name": "deis-builder",
          "description": "Deis Builder"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-31T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-controller",
          "description": "Deis Controller"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-30T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-database",
          "description": "Deis Database"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-29T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-minio",
          "description": "Deis Minio"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-28T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-registry",
          "description": "Deis Registry"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-27T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-router",
          "description": "Deis Router"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-26T23:54:39Z"
        }
      },
      {
        "component": {
          "name": "deis-workflow-manager",
          "description": "Deis Workflow Manager"
        },
        "version": {
          "train": "stable",
          "version": "2.0.0",
          "released": "2016-03-25T23:54:39Z"
        }
      }
    ]
  }
  ```
