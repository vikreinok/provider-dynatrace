# Regeneration Guide

This guide describes how to update the Terraform Provider version that `provider-dynatrace` depends on, and how to regenerate all corresponding Crossplane models and Custom Resource Definitions.

## 1. Update the Makefile

Open the root `Makefile` and locate the Terraform provider exports:

```makefile
export TERRAFORM_PROVIDER_VERSION ?= 1.91.1
export TERRAFORM_NATIVE_PROVIDER_BINARY ?= terraform-provider-dynatrace_v1.91.1
```

Change the version numbers to the desired newer release (e.g., `1.95.0`). You can find recent releases on the [terraform-provider-dynatrace GitHub repository](https://github.com/dynatrace-oss/terraform-provider-dynatrace/releases).

## 2. Run the Generation Process

From the root of the repository, execute:

```sh
make generate
```

This target handles the following automatically:
1. Downloads the specific terraform release binary (`terraform-provider-dynatrace_v1.x.y`).
2. Generates the provider schema by calling the native Terraform API plugin bindings.
3. Generates Crossplane Go types (APIs) for every resource configured in `config/iam/config.go`.
4. Generates corresponding OpenAPI validation types for the Custom Resource Definitions (CRDs).

## 3. Review the Changes

The `make generate` command frequently changes `apis/` and `package/crds/` directories.
It's recommended to do a standard git review on those changes before committing.

```sh
git status
git diff
```

Pay close attention to changes in new required fields, structurally rearranged inputs, or renaming within the IAM APIs, as they can break existing manifests.
