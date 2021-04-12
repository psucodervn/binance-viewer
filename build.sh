#!/usr/bin/env bash

function build_app {
  BINARY=copytrader
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o $BINARY main.go
}

function build_docker {
  ENV=$1
  TAG=$2

  case ${ENV} in
    dev | dev-aws)
      REPOSITORY=psucoder
      TAG="dev-${TAG}"
      ;;
    master | alpha)
      REPOSITORY=psucoder
      TAG="prod-${TAG}"
      ;;
    *)
      echo "Invalid env ${ENV}"
      exit 1
  esac

  IMAGE=${REPOSITORY}/copytrader:${TAG}
  echo "Build ${IMAGE}"

  docker build -t "${IMAGE}" -f ./Dockerfile .
  docker push "${IMAGE}"
}

TAG="v$(date -u +"%Y%m%d")-$(git rev-parse --short HEAD)"
BRANCH=${1:-$(git rev-parse --abbrev-ref HEAD)}
COMMIT_MESSAGE=$(git log -1 --pretty=format:%B)

# if the commit message contains pattern [no-build]
if [[ ${COMMIT_MESSAGE} == *"[no-build]"* || ${COMMIT_MESSAGE} == *"[skip-build]"* ]]; then
  echo "Skip build base on git commit message"
  exit 0
fi

build_app

build_docker "${BRANCH}" "${TAG}"

exit 0
