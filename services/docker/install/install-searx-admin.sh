#!/bin/sh
set -x
set -e

clear
echo

DIR=$(dirname "$0")
echo "$DIR"
. ${DIR}/internal/common.sh
. ${DIR}/internal/env-aliases.sh
. ${DIR}/internal/vcs-git.sh

# Set temp environment vars
export SEARX_ADMIN_VCS_URI=${SEARX_ADMIN_VCS_URI:-"github.com/kvch/searx-admin.git"}
export SEARX_ADMIN_VCS_BRANCH=${SEARX_ADMIN_VCS_BRANCH:-"master"}
export SEARX_ADMIN_VCS_CLONE_DEPTH=${SEARX_ADMIN_VCS_CLONE_DEPTH:-"1"}
export SEARX_ADMIN_VCS_CLONE_PATH=${SEARX_ADMIN_VCS_CLONE_PATH:-"/app/searx-admin"}

pip install --upgrade pip

if [[ -d ${SEARX_ADMIN_VCS_CLONE_PATH} ]]; then
	rm -fR ${SEARX_ADMIN_VCS_CLONE_PATH}
fi

# Clone, Compile & Install
git clone -b ${SEARX_ADMIN_VCS_BRANCH} --recursive --depth ${SEARX_ADMIN_VCS_CLONE_DEPTH} -- https://${SEARX_ADMIN_VCS_URI} ${SEARX_ADMIN_VCS_CLONE_PATH}
cd ${SEARX_ADMIN_VCS_CLONE_PATH}

pip install --no-cache --no-cache-dir -r ./requirements.txt
