# Local Testing Guide for provider-dynatrace

This guide explains how to run and test `provider-dynatrace` locally securely without needing a full-fledged Kubernetes cluster. We will use a lightweight Kubernetes cluster (`kind`) exclusively for the API server so that we have a place to define Crossplane resources and ProviderConfigs, while the provider itself runs natively on your machine via `make run`.

## Prerequisites

1.  **Docker**: Required to run `kind`.
2.  **Kind**: Kubernetes IN Docker ([installation instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)).
3.  **Go**: Ensure a recent version of Go is installed (>=1.24).
4.  **Dynatrace Environment URL and API Token**: Required to authenticate against your SaaS instance (e.g., `abc12345`).
    - Note that you should have correct permissions assigned to the token so operations can be executed.

## Step 1: Create a Local Control Plane

Create a minimal local cluster using `kind`:

```sh
kind create cluster --name dynatrace-test
```

Install Crossplane into this local cluster:

```sh
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm upgrade --install crossplane crossplane-stable/crossplane \
    --namespace crossplane-system \
    --create-namespace \
    --wait
```

## Step 2: Apply the Custom Resource Definitions (CRDs)

Run the following command in the provider root directory to generate the code and then apply the necessary CRDs to your `kind` cluster:

```sh
make generate
kubectl apply -f package/crds/
```

## Step 3: Configure the ProviderCredentials

Since we are running out-of-cluster, we provide the credentials natively. First, define a `Secret` in the cluster that contains our Dynatrace API token:

1. Create a file `secret.yaml`:
    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: dynatrace-creds
      namespace: crossplane-system
    type: Opaque
    stringData:
      credentials: |
        {
          "dt_env_url": "https://abc12345.live.dynatrace.com",
          "dt_api_token": "YOUR_SUPER_SECRET_TOKEN"
        }
    ```
    *Note: The keys required for authentication may vary depending on the `provider-dynatrace` schema, typically Crossplane uses JSON configuration blobs that match Terraform provider variables. Verify exact JSON keys against the Dynatrace Terraform Provider documentation (https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs).*

2. Apply it:
    ```sh
    kubectl apply -f secret.yaml
    ```

3. Configure the provider to use the secret by creating `providerconfig.yaml`:
    ```yaml
    apiVersion: dynatrace.crossplane.io/v1beta1
    kind: ProviderConfig
    metadata:
      name: default
    spec:
      credentials:
        source: Secret
        secretRef:
          name: dynatrace-creds
          namespace: crossplane-system
          key: credentials
    ```
4. Apply it:
    ```sh
    kubectl apply -f providerconfig.yaml
    ```

## Step 4: Run the Provider Locally

Now, boot up the provider from your host machine pointing to your `kind` cluster's context:

```sh
make run
```
You should see logging indicating the provider has successfully started.

### How references and ordering work

Just like Terraform resolves `dynatrace_iam_group.svc[each.key].id` at apply time, Crossplane resolves resource references automatically via its reconciliation loop:

| Terraform reference | Crossplane equivalent |
|---|---|
| `dynatrace_iam_group.svc[each.key].id` | `groupRef: name: sv-xyz-dev-analyst` |
| `dynatrace_iam_policy_boundary.svc[each.key].id` | `boundariesRefs: [{name: sv-xyz-dev}]` |
| `policy.value` (pre-existing UUID) | `id: "1348b750-..."` (raw UUID — not Crossplane-managed) |

**Ordering is automatic**: when `PolicyBindingsV2` has a `groupRef` or `boundariesRefs`, Crossplane waits until the referenced resource is `READY` before resolving its UUID and proceeding. You do not need to apply resources in a specific sequence.

### Apply everything at once

```sh
kubectl apply -f examples/cluster/iam/iam-v2.yaml
```

Then watch resources converge:

```sh
watch kubectl get managed
```

You should see Groups and Boundaries become `READY` first, then the `PolicyBindingsV2` will resolve their IDs and sync shortly after.


## Step 6: Drift Correction (Reconciliation)

`make run` is the **reconciliation engine**. As long as it is running, it will continuously monitor the state in Dynatrace and compare it against your Custom Resources (CRs).

To test drift correction:
1. Ensure your resources are `READY: True`.
2. Go to the Dynatrace console and manually change a property (e.g., change the `description` of the `sv-xyz-dev-analyst` Group).
3. Wait for the next reconciliation cycle. 
    - **Note**: I have updated the provider to poll every **1 minute** by default for faster local testing.
4. **Trigger immediate reconciliation**: To see the change reverted immediately without waiting a minute, you can:
    - Update a non-functional field in the CR (like adding an annotation).
    - Or simply **restart the `make run` process**.
5. Check the Dynatrace console again; the manual change should be reverted to match the CR.

## Cleanup

```sh
kubectl delete -f examples/cluster/iam/iam-v2.yaml
kind delete cluster --name dynatrace-test
```

## Troubleshooting

If your resources do not transition to `READY: True`, then the Dynatrace API is likely rejecting the provider's request. 

1. **Check the Resource Conditions**
   To see the exact error message returned from the Dynatrace API (e.g., "invalid token" or "insufficient permissions"), describe the object:
   ```sh
   kubectl describe group.iam sv-xyz-dev-analyst
   ```
   Look at the `Events` and `Status.Conditions` sections at the bottom of the output.

2. **Check the Provider Logs**
   If you are running the provider locally via `make run`, simply look at the standard output in that terminal. Every API interaction, drift detection loop, and authentication failure will be logged natively here.
