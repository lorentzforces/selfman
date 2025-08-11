## Big TODOS

Need to update git app handling to rebuild if/when you've gotten fresh commits on a branch. Right now selfman will see that the version name has not changed, and therefore won't rebuild the app since the target artifact already exists. For now this can be clunkily worked around by running the 'remove' operation and then updating the app again, but this is basically the standard usage for a git app and so should be fixed. Probably could be done via a git-aware build operation (new operation type) which is initialized with a commit hash, and then skips the build if the current commit hash matches that hash. Alternatively (and this may be preferable for branch-switching reasons), might consider adding the commit hash to the artifact naming convention so it can be compared later.

Relatedly, the check command needs some work:
- git apps show available branch versions, but not available artifact versions. Both should probably be shown.
- web-fetch apps which are only a single file (like lombok) will always show as an inconsistent state because no files are left in the source directory when the target is moved
- If an app has a lot of versions available, the formatting will probably be crap. This is probably puntable until I have an app which this actually affects, but something like a columnar display (3 columns max or something) may be good.

Since selfman is now basically usable for all my initial use cases, I should add a verbose option to CLI commands and make the default less-verbose.

The remove command should have a `--remove-source` option to also delete the source directory for an app.

There should be a cleanup command which allows the user to check and remove outdated/unused versions.

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
```

> **NOTE:** For the purposes of the source directory, the version label for a git app is always "git"
