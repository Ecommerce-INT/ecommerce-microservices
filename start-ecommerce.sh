#!/usr/bin/env bash
# =============================================================================
# start-ecommerce.sh — One-shot local K8s deployment for Linux / macOS / WSL / Git Bash.
#
# Works with Docker or Podman (auto-detected; Docker wins when both are
# available). Force a runtime with:
#   RUNTIME=podman ./start-ecommerce.sh
#
# Delegates to k3d-setup.sh, which installs k3d/kubectl if missing.
# =============================================================================

set -euo pipefail
cd "$(dirname "$0")"

chmod +x k3d-setup.sh
exec ./k3d-setup.sh "$@"
