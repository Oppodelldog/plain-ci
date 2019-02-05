#!/bin/bash
buildId=${PLAIN_CI_BUILD_ID}
repoURL=${PLAIN_CI_REPO_URL}
commitHash=${PLAIN_CI_COMMIT_HASH}
originator=${PLAIN_CI_ORIGINATOR}
error=${PLAIN_CI_ERROR}

echo "sample post-build notification"
echo "------------------------------------------"
echo Build ID   : ${buildId}
echo Repo URL   : ${repoURL}
echo Commit-Hash: ${commitHash}
echo Originator : ${originator}
echo Error      : ${error}