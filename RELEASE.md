# Release Guide for provider-dynatrace

This guide provides instructions on how to build and release the Dynatrace Crossplane provider as a public Docker image and an OCI-compliant Crossplane package (xpkg).

## Target Hub
- **Docker Hub Organization**: `vikreinok`
- **Image Name**: `provider-dynatrace`
- **XPKG Name**: `provider-dynatrace`

---

## Local Build and Release

Follow these steps to build and push the image from your local development machine.

### 1. Login to Docker Hub
Ensure you are authenticated with your `vikreinok` account:
```sh
docker login
```

### 2. Build the Multi-Arch Docker Image
Upjet uses `docker buildx` under the hood to build for multiple platforms (`linux/amd64` and `linux/arm64`). 
Replace `<version>` with your desired tag (e.g., `v0.1.0`).

```sh
# Build images for all supported platforms
REGISTRY_ORGS=vikreinok VERSION=<version> make build.all
```

### 3. Push the Docker Image
```sh
REGISTRY_ORGS=vikreinok VERSION=<version> make image.push.all
```

### 4. Build and Push the Crossplane Package (XPKG)
The XPKG is the OCI-compliant package that users install via `kubectl crossplane pkg install`.

```sh
# Push the package to Docker Hub
XPKG_REG_ORGS=vikreinok VERSION=<version> make xpkg.push.all
```

---

## GitHub Actions Release Pipeline

To automate the release process, you can use the built-in GitHub Actions workflows.

### 1. Configure Repository Secrets
Go to your GitHub repository settings and add the following secrets:
- `DOCKERHUB_USERNAME`: Your Docker Hub username (`vikreinok`).
- `DOCKERHUB_TOKEN`: A Personal Access Token from Docker Hub.

### 2. Create the Release Workflow
Create or update `.github/workflows/publish-dockerhub.yml` with the following content:

```yaml
name: Publish to Docker Hub

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:
    inputs:
      version:
        description: "Version to publish (e.g. v0.1.0)"
        required: true

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: true

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Publish Artifacts
        run: |
          # Use the tag name or input version
          VERSION=${{ github.event.inputs.version || github.ref_name }}
          REGISTRY_ORGS=vikreinok XPKG_REG_ORGS=vikreinok VERSION=$VERSION make publish
```

### 3. Triggering a Release
- **Automatically**: Push a git tag starting with `v` (e.g., `git tag v0.1.0 && git push origin v0.1.0`).
- **Manually**: Go to the **Actions** tab in GitHub, select the **Publish to Docker Hub** workflow, and click **Run workflow**, providing the version string.

---

## GitHub Container Registry Release (GHCR)

You can also publish the provider to the **GitHub Container Registry (GHCR)**. This is useful for distributing the provider within the GitHub ecosystem.

### 1. Generate a GitHub Personal Access Token (PAT)
To allow the workflow to push to GHCR, you need a PAT with the appropriate permissions:
1. Go to your GitHub **Settings** -> **Developer settings** -> **Personal access tokens** -> **Tokens (classic)**.
2. Click **Generate new token (classic)**.
3. Give it a name (e.g., `ghcr-publish-provider-dynatrace`).
4. Select the following scopes:
    - `write:packages`: To upload the images and XPKGs.
    - `read:packages`: To verify existing packages.
5. Click **Generate token** and copy it immediately.

### 2. Configure GitHub Secrets
Add the token to your GitHub repository secrets:
1. Navigate to your repository on GitHub.
2. Go to **Settings** -> **Secrets and variables** -> **Actions**.
3. Click **New repository secret**.
4. Name: `GHCR_TOKEN`.
5. Value: Paste your generated PAT.

### 3. Trigger the GHCR Publish
The workflow is defined in [`.github/workflows/publish-ghcr.yml`](file:///Users/viktor/jb/provider-dynatrace/.github/workflows/publish-ghcr.yml).

1. Go to the **Actions** tab in your repository.
2. Select the **Publish to GHCR** workflow.
3. Click **Run workflow**.
4. Enter the version (e.g., `v0.1.0`).
5. Click **Run workflow**.

This will build and push both the controller image and the Crossplane package to `ghcr.io/${{ github.actor }}/provider-dynatrace`.

---

## Upbound Marketplace Release (OCI)

To publish your provider to the [Upbound Marketplace](https://marketplace.upbound.io) or any OCI-compliant registry using the official Crossplane workflows, use the **`Publish Provider Package`** workflow.

### 1. Generate Upbound Tokens
To push to the Upbound Marketplace, you need a **Token Pair** from the Upbound Console:
1. Log in to [console.upbound.io](https://console.upbound.io).
2. Navigate to your **Organization** or **Account Settings**.
3. Go to **Settings** -> **Tokens**.
4. Click **Create Token**.
5. Give it a name (e.g., `github-actions-provider-dynatrace`).
6. **Copy the values immediately**:
    - The **Token ID** (this will be your `XPKG_MIRROR_ACCESS_ID`).
    - The **Token Secret** (this will be your `XPKG_MIRROR_TOKEN`).

### 2. Configure GitHub Secrets
Add these tokens to your GitHub repository secrets:
- `XPKG_MIRROR_ACCESS_ID`: Paste the Token ID here.
- `XPKG_MIRROR_TOKEN`: Paste the Token Secret here.

### 3. Trigger the OCI Publish
The workflow is defined in [`.github/workflows/publish-provider-package.yml`](file:///Users/viktor/jb/provider-dynatrace/.github/workflows/publish-provider-package.yml). 

1. Go to the **Actions** tab.
2. Select **Publish Provider Package**.
3. Click **Run workflow**.
4. Enter the version (e.g., `v0.1.0-alpha.1`) and the Go version if different from default.
5. Click **Run workflow**.

This will build the provider, package it as an OCI image (`xpkg`), and push it to the mirror registry configured in the shared workflow (defaults to Upbound Marketplace).

---

## Troubleshooting

### Cross-Platform Build Error (Mac ARM / Apple Silicon)

If you see an error like `No such image: build-xxx/provider-dynatrace-amd64:latest` when running `make build.all` on an ARM Mac, it is because Crossplane's `xpkg build` command requires the target architecture's image to be present in your local Docker daemon. 

On Mac ARM, building and loading an `amd64` image can be flaky or require QEMU emulation.

**Recommendation for Local Development:**
Only build for your host architecture (ARM64) to avoid emulation overhead and errors:

```sh
# Build only for the host platform
make build PLATFORMS=linux_arm64
```

**Recommendation for Releasing:**
Always use the **GitHub Actions pipeline** for multi-arch releases. The CI environment is specifically configured to build and push `amd64` and `arm64` artifacts natively and reliably, ensuring the final public image is correctly balanced for all users.
