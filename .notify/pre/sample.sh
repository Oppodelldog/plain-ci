#!/bin/bash

buildId=${SIMPLE_CI_BUILD_ID}
repoURL=${SIMPLE_CI_REPO_URL}
commitHash=${SIMPLE_CI_COMMIT_HASH}
originator=${SIMPLE_CI_ORIGINATOR}
error=${SIMPLE_CI_ERROR}

echo "sample pre-build notification"
echo "------------------------------------------"
echo Build ID   : ${buildId}
echo Repo URL   : ${repoURL}
echo Commit-Hash: ${commitHash}
echo Originator : ${originator}
echo Error      : ${error}