Filesystem layout is as follows:

```
- xdg-data/
  + selfman/
    + artifacts/
    | + [app-name]---[version-label] (binary)
    | + ...
    + sources/
      + [app-name]/
      | + [version-label]/
      |   + (source files)
      + ...
- xdg-config/
  + selfman/
    + config.yaml (optional)
    + apps/
      + [app-name].config.yaml
```

> **NOTE:** For the purposes of the source directory, the version label for a git app is always "git"

Major changes from earlier ideas/current state:
- operations which get the app status will now specify a version label and the status only considers that version
- git apps check for a rev which matches the version label (like a branch or tag)
- we mostly don't concern ourselves with which version is curently in use - we re-create the symlink to the binary no matter what, and all other operations are only concerned with state on the filesystem

**TODO:** Change the filesystem abstraction to handle dealing with filepaths so testing doesn't need to give a shit about path strings (possibly construct a filesystem abstraction object with a reference to an app config, so it can ask the app config for file paths if necessary, and also determine what flavor of app it is, etc)
