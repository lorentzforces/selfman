# selfman

For self-managed programs.

## Goal

A reasonable way to manage self-built stuff.

Stuff managed by `selfman` might be a tool that I've written myself, but also might be an open-source project that I have made modifications to.


## An alternative solution for the data

- yaml/toml
- programs/tools/managed-tools
  - name
  - how to download
  - how to update
  - how to build
    - my own make convention
    - is literally just a binary
    - location of output executable
    - an arbitrary script
  - how to make available/install
    - link to artifact
    - run a command

question: what if it's just a binary?
question: how do we configure things so that if the thing follows a "standard" pattern, you can just use that off-the-shelf, and then override selectively things that you want?
  possible answer: Just use pointers on everything. Have a bunch of configurations with names we can pull based on a field set in configuration (if present). If that is, then we can take any further-defined fields and apply those values to the configuration.

## Data model

Base config
- basic things to configure Selfman itself
  - executable location (for sourcing)

App configs
- separate config files in a particular directory or with a naming convention

- [program|tool|managed-tool]
  - name
  - how to download (source)
    - git repo
    - filesystem
    - other things TBD
  - how to update (source)
    - git pull
    - other things...?
  - how to build
    - my own make convention
    - an arbitrary script
    - location of output executable

### Directories

- config (`$XDG_CONFIG_HOME/selfman`)
- application configs (`$XDG_CONFIG_HOME/selfman/apps`)
- application repositories (`$XDG_DATA_HOME/selfman/sources`)
- application binaries or other artifacts (`$XDG_DATA_HOME/selfman/artifacts`)
  - TODO: figure out how we track old versions etc

### Additional TODO

Write a bootstrap bash script that:

- installs go
- downloads the selfman repo into a temp directory
- builds and inits selfman
- OPEN QUESTION: is selfman config portable? what about the DB? do these need to be separate?
  - probably doesn't need to be portable (chosen applications are likely to vary between systems), but exporting an app's configuration probably is useful
- adds ensure-path into selfman config
