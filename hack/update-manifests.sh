#!/bin/sh

set -e

SRCROOT="$( CDPATH='' cd -- "$(dirname "$0")/.." && pwd -P )"
AUTOGENMSG="# This is an auto-generated file. DO NOT EDIT"

update_image () {
  if [ ! -z "${IMAGE_NAMESPACE}" ]; then
    sed 's| image: \(.*\)/\(argo-rollouts.*\)| image: '"${IMAGE_NAMESPACE}"'/\2|g' "${1}" > "${1}.bak"
    mv "${1}.bak" "${1}"
  fi
}

if [ ! -z "${IMAGE_TAG}" ]; then
  (cd ${SRCROOT}/manifests/base && kustomize edit set image argoproj/argo-rollouts:${IMAGE_TAG})
fi

kustomize version

echo "${AUTOGENMSG}" > "${SRCROOT}/manifests/install.yaml"
KUSTOMIZE_TMP_OVERLAY="${SRCROOT}/manifests/tmp"
mkdir -p "$KUSTOMIZE_TMP_OVERLAY"
(cd "$KUSTOMIZE_TMP_OVERLAY"; kustomize create --resources "../cluster-install" --namespace argo-rollouts)
kustomize build --load_restrictor none "$KUSTOMIZE_TMP_OVERLAY" >> "${SRCROOT}/manifests/install.yaml"
rm -r "$KUSTOMIZE_TMP_OVERLAY"
update_image "${SRCROOT}/manifests/install.yaml"

echo "${AUTOGENMSG}" > "${SRCROOT}/manifests/namespace-install.yaml"
kustomize build --load_restrictor none "${SRCROOT}/manifests/namespace-install" >> "${SRCROOT}/manifests/namespace-install.yaml"
update_image "${SRCROOT}/manifests/namespace-install.yaml"
