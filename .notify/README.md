# plain ci notification api

Plain ci is so plain, it has no built in plugins.
But it provides support for custom actions via its notification api.

Whenever a build reached a noteworthy step, plain-ci will call
external programs or scripts that gives you all freedom you need.

## pre build
**pre build** scripts are called right after the build was enqueuedm but before the
build was started. 

If one of those scripts fails (exit code != 0), the build will fail.

## post build
**post build** scripts are called right after the build has finished.

If one of those script fails (exit code != 0), an error will be logged.