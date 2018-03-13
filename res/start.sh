#!/bin/bash

trap '[ "$?" -eq 0 ] || read -p "Looks like something went wrong in step ´$STEP´... Press any key to continue..."' EXIT

DOCKER_MACHINE=./docker-machine.exe

STEP="Looking for VBoxManage.exe"
if [ ! -z "$VBOX_MSI_INSTALL_PATH" ]; then
  VBOXMANAGE="${VBOX_MSI_INSTALL_PATH}VBoxManage.exe"
else
  VBOXMANAGE="${VBOX_INSTALL_PATH}VBoxManage.exe"
fi

BLUE='\033[1;34m'
GREEN='\033[0;32m'
NC='\033[0m'

#clear all_proxy if not socks address
if  [[ $ALL_PROXY != socks* ]]; then
  unset ALL_PROXY
fi
if  [[ $all_proxy != socks* ]]; then
  unset all_proxy
fi

if [ ! -f "${DOCKER_MACHINE}" ]; then
  echo "Docker Machine is not installed."
  exit 1
fi

if [ ! -f "${VBOXMANAGE}" ]; then
  echo "VirtualBox is not installed."
  exit 1
fi

"${VBOXMANAGE}" list vms | grep \""${MACHINE_NAME}"\" &> /dev/null
VM_EXISTS_CODE=$?

set -e

STEP="Checking if machine $MACHINE_NAME exists"
echo "Machine storage path: ${MACHINE_STORAGE_PATH}"
if [ $VM_EXISTS_CODE -eq 0 -a ! -z ${MACHINE_STORAGE_PATH} -a ! -d "${MACHINE_STORAGE_PATH}/machines/${MACHINE_NAME}" ]; then
  "${DOCKER_MACHINE}" rm -f "${MACHINE_NAME}" &> /dev/null || :
fi
if [ $VM_EXISTS_CODE -eq 1 -a ! -z ${MACHINE_STORAGE_PATH} ]; then
  "${DOCKER_MACHINE}" rm -f "${MACHINE_NAME}" &> /dev/null || :
  rm -rf "${MACHINE_STORAGE_PATH}/machines/${MACHINE_NAME}"
  if [ "${HTTP_PROXY}" ]; then
    PROXY_ENV="$PROXY_ENV --engine-env HTTP_PROXY=$HTTP_PROXY"
  fi
  if [ "${HTTPS_PROXY}" ]; then
    PROXY_ENV="$PROXY_ENV --engine-env HTTPS_PROXY=$HTTPS_PROXY"
  fi
  if [ "${NO_PROXY}" ]; then
    PROXY_ENV="$PROXY_ENV --engine-env NO_PROXY=$NO_PROXY"
  fi
  "${DOCKER_MACHINE}" create -d virtualbox $PROXY_ENV \
    --virtualbox-hostonly-cidr "${MACHINE_HOST_CIDR}" \
    --virtualbox-cpu-count "${MACHINE_CPU}" \
    --virtualbox-memory "${MACHINE_RAM}" \
    --virtualbox-disk-size "${MACHINE_DISK}" \
    --virtualbox-share-folder "\\\?\\${MACHINE_SHARED_PATH}:${MACHINE_SHARED_NAME}" \
    "${MACHINE_NAME}"
fi

STEP="Checking status on $MACHINE_NAME"
VM_STATUS="$( set +e ; ${DOCKER_MACHINE} status ${MACHINE_NAME} )"
if [ "${VM_STATUS}" != "Running" ]; then
  "${DOCKER_MACHINE}" start "${MACHINE_NAME}"
  yes | "${DOCKER_MACHINE}" regenerate-certs "${MACHINE_NAME}"
fi

STEP="Setting env"
eval "$(${DOCKER_MACHINE} env --shell=bash --no-proxy ${MACHINE_NAME})"

STEP="Finalize"
clear
cat << EOF


                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/

EOF
echo -e "${BLUE}docker${NC} is configured to use the ${GREEN}${MACHINE_NAME}${NC} machine with IP ${GREEN}$(${DOCKER_MACHINE} ip ${MACHINE_NAME})${NC}"
echo -e "Shared folder is named ${BLUE}${MACHINE_SHARED_NAME}${NC} and is located in ${GREEN}data/shared${NC}"
echo "For help getting started, check out the docs at https://docs.docker.com"
echo
cd

docker () {
  MSYS_NO_PATHCONV=1 docker.exe "$@"
}
export -f docker

if [ $# -eq 0 ]; then
  echo "Start interactive shell"
  exec "$BASH" --login -i
else
  echo "Start shell with command"
  exec "$BASH" -c "$*"
fi
