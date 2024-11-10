# action-usrmerge

An Immutos action that merges the /usr directory into the root filesystem,
see: [wiki.debian.org/UsrMerge](https://wiki.debian.org/UsrMerge).

We use a go binary here instead of a shell script to avoid chicken and egg
problems that arise when moving shell utilities around.