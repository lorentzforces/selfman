## Big TODOS

List command should tolerate bad configs and display those errors alongside a list of the valid configurations.

Consider: should commands which operate on a specific app ignore any parsing or configuration issues with non-specified apps?

Related to git app woes, the check command needs some work:
- git apps show available branch versions, but not available artifact versions. Both should probably be shown.
- If an app has a lot of versions available, the formatting will probably be crap. This is probably puntable until I have an app which this actually affects, but something like a columnar display (3 columns max or something) may be good.

There should be a cleanup command which allows the user to check and remove outdated/unused versions. (and potentially some way of nuking absolutely everything related to an app)

Documentation needs to be updated quite badly (at this point even for my own sake), especially:
- app flavors and what configuration is valid for each flavor
- how placeholders work and which fields they can be used in

## Implementation notes

Filesystem layout is as follows:

"version-label" escapes path separators by substituting the string "%SLASH%"

e.g. `origin/main` -> `origin%SLASH%main`

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
- lib-dir/ (usually ~/.local/lib)
  + [app-name] (links to source dir)
```

> **NOTE:** For the purposes of the source directory, the version label for a git app is always "git"
