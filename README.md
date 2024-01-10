# Cloudentity Configuration As Code (Early Access)

This repository contains a CLI tool for managing Cloudentity configuration.

## Installation

### As a go package

```bash
go install github.com/cloudentity/cac
```

## Commands

### Help

Prints help message with available commands and their parameters.

```bash
cac --help
```

### Pull

Pull configuration from Cloudentity and save it to a directory structure.

```bash
cac --config examples/e2e/config.yaml pull --workspace cdr_australia-demo-c67evw7mj4
```

#### Sample output

```
/tmp/e2e-data
└── /tmp/e2e-data/workspaces
    └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/bank2.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/bank.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/Consent_Page_Bank_Client.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/Data_Holder.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/Financroo.yaml
        │   └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/clients/xxx.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/idps
        │   └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/idps/test.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Block_API.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Block_DCR.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Block_Developer.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Block_Machine.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Block_User.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/MFA_User.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/mtls2.rego
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/mtls2.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/mtls.rego
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/mtls.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-1_API.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-1_User.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-2_API.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-2_User.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-3_API.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/NIST-AAL-3_User.yaml
        │   └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policies/Unlock_DCR.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/scripts
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/scripts/debug.js
        │   └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/scripts/debug.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services/CDR_Australia.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services/OAuth2.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services/Profile.yaml
        │   ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services/Transient_One-Time_Passwords.yaml
        │   └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/services/User_Privacy_&_Consent.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/claims.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/consent.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/policy_execution_points.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/scopes.yaml
        ├── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/script_execution_points.yaml
        └── /tmp/e2e-data/workspaces/cdr_australia-demo-c67evw7mj4/server.yaml
```

### Push

Merge configuration from a directory structure and push it into Cloudentity.

#### Push configuration from multiple directories

To push configration from multiple directories, either pass an array to the `storage.dir_path` or use `STORAGE_DIR_PATH` with multiple paths split by a comma.

Configurations are merged in the reverse order, so the first path has the highest priority, and will override everything else.

```bash
cac --config examples/e2e/config.yaml push --workspace cdr_australia-demo-c67evw7mj4
```

## Templates

Templates are used to generate configuration files. They are using [Go template language](https://golang.org/pkg/text/template/).

### Functions

We use [Sprig](http://masterminds.github.io/sprig/) library to extend Go template language with additional functions and also provide several custom functions.

|  Function | Description                                               |
|----------:|:----------------------------------------------------------|
|   include | Includes a template file                                  |
|       env | Reads an environment variable, and fails if it is not set |