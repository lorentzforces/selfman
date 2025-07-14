Filesystem layout is as follows:

"version-label" escapes path separators by substituting the string "#SLASH#"

e.g. `origin/main` -> `origin#SLASH#main`

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

**TODO:** potentially merge "install" and "update" operations to just be "make it so" for a given app

CURRENT OBSTACLE: jdtls's executable (a python script) expects that it is living in the file structure you download, so moving the executable out doesn't work. Some solutions might include:
- linking directly into the sources dir
  - what does the config look like for this? do we do it just based on app type? (probably not?)
- copying the source dir as an artifact dir, and then linking into that
  - duplication of the whole source dir seems pointless
