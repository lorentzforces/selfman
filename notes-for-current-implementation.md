Folders are as follows:

- xdg-data
  - selfman
    - apps-tracking.yaml
    - artifacts
      - [appname] (binary)
      - ...
    - sources
      - [appname]
        - (source files)
      - ...
- xdg-config
  - selfman
    - config.yaml (optional)
    - apps
      - [appname].config.yaml

Apps-tracking format:

```yaml
apps:
- selfman:
  current-version: main
  versions:
  - main
- ensure-path:
  current-version: main
  versions:
  - main
```
