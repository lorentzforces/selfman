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

Documentation needs to be updated quite badly (at this point even for my own sake), especially:
- app flavors and what configuration is valid for each flavor
- how placeholders work and which fields they can be used in
