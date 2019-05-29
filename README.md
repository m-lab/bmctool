[![GoDoc](https://godoc.org/github.com/m-lab/bmctool?status.svg)](https://godoc.org/github.com/m-lab/bmctool) [![Build Status](https://travis-ci.org/m-lab/bmctool.svg?branch=master)](https://travis-ci.org/m-lab/bmctool) [![Coverage Status](https://coveralls.io/repos/github/m-lab/bmctool/badge.svg?branch=master)](https://coveralls.io/github/m-lab/bmctool?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/bmctool)](https://goreportcard.com/report/github.com/m-lab/bmctool)

# BMCTool
BMCTool is a command line tool to manage Baseband Management Controller credentials on the M-Lab infrastructure.

In particular, it allows to:

* Fetch credentials for the BMC module on a given node from Google Cloud Datastore (GCD)
* Update the credentials for an existing node
* Add a new node to GCD

Output is provided in JSON format.

## Usage
### Common flags
```--project <project_id>```

Use the specified `<project_id>` to connect to GCD.

### Fetch credentials
```./bmctool get <host>```

Retrieves the credentials for `<host>`.

### Add a new node
(TODO)

Set the shell variables `BMCUSER` and `BMCPASS` to the appropriate values.

```./bmctool add <host> <address>```

Creates the node `<host>` with the provided `BMCUSER`, `BMCPASS` and `<address>`. If the specified node already exists, the command will fail.

### Update an existing node
(TODO)

Set the shell variables `BMCUSER` and `BMCPASS` to the appropriate values.

```./bmctool update <host> <address>```

Update details for node `<host>`. If `<host>` does not exist, the command will fail.
