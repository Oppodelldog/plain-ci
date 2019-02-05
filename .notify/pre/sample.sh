#!/bin/bash

buildId=${PLAIN_CIBUILD_ID}
repoURL=${PLAIN_CIREPO_URL}
commitHash=${PLAIN_CICOMMIT_HASH}
originator=${PLAIN_CIORIGINATOR}
error=${PLAIN_CIERROR}

echo "sample pre-build notification"
echo "------------------------------------------"
echo Build ID   : ${buildId}
echo Repo URL   : ${repoURL}
echo Commit-Hash: ${commitHash}
echo Originator : ${originator}
echo Error      : ${error}