[![GoDoc](https://godoc.org/github.com/m-lab/bmctool?status.svg)](https://godoc.org/github.com/m-lab/bmctool) [![Build Status](https://travis-ci.org/m-lab/bmctool.svg?branch=master)](https://travis-ci.org/m-lab/bmctool) [![Coverage Status](https://coveralls.io/repos/github/m-lab/bmctool/badge.svg?branch=master)](https://coveralls.io/github/m-lab/bmctool?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/m-lab/bmctool)](https://goreportcard.com/report/github.com/m-lab/bmctool)

# BMCTool
BMCTool is a command line tool to manage Baseband Management Controller credentials on the M-Lab infrastructure.

In particular, it allows to:

* Fetch credentials for the BMC module on a given node from Google Cloud Datastore
* Add a new node
* Update the credentials for an existing node
* Delete credentials for an existing node
* Reboot a node via its BMC, using the credentials in GCD to log in
* Automatically set up SSH forwarding to access the BMC's web interface and virtual console

If an entity is created or updated, output is provided in JSON format.

## Usage

### Common flags

```--project <project_id>```

Overrides the auto-detected project ID (based on the hostname format) with the specified `<project_id>`.

### Fetch credentials

```bmctool get <host>```

Retrieves the credentials for `<host>`.

### Add a new node

```bmctool add <host> <address>```

Set the shell variables `BMCUSER` and `BMCPASS` to the appropriate values before running this command.

Creates the node `<host>` with the provided `BMCUSER`, `BMCPASS` and `<address>`. If the specified node already exists or the required env variables are not set, the command will fail.

### Update an existing node

```bmctool set <host>```

Updates details for node `<host>`. If `<host>` does not exist, the command will fail.

A list of fields that can be updated can be retrieved with `bmctool set --help`.

For example, to set the address to 127.0.0.1

```bmctool set mlab1.lga0t --addr 127.0.0.1```


### Delete an existing node

```bmctool delete <host>```

Deletes the Credentials entity for `<host>`. This works in the same way as Datastore's `Delete` operation, thus if `<host>` does not exist, this command will not fail.

### Reboot a node

```bmctool reboot <host>```

Logs into the BMC for `<host>` and reboots the node. Credentials are fetched from GCD.

### Forward ports via SSH

```bmctool forward <host>```

Runs the system-wide SSH command to connect to a bastion host and forward ports to the BMC. This is normally used to access the BMC's web interface and its Virtual Console, thus the default behavior is to map ports 4443 -> 443 (to not require root privileges) and 5900 -> 5900. For a list of the available options:  `bmctool forward --help`.
