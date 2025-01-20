# selfman

*Wrangling self-managed tools and applications*

Dealing with self-built or self-managed tools is a bit painful. This project aims to create a straightforward way of managing tools or applications you want to build yourself from source. I started this project after writing a number of utilities for myself, and finding the process of getting them added and updated across my various machines (running various flavors of Linux and MacOS) to be a pain.

I've also run into some friction while fixing bugs on open-source tools I use; when a bug is fixed on the project's main branch, I want to compile the project with all the latest changes. However, once the project has all the fixes and features I want, I want to easily revert back to using standard distribution packages.

This project **does not** aim to be a full-featured package manager. In particular, it does not and will not track dependencies, and many features which would be considered mandatory in a package manager may never be included. It is also not really intended to scale beyond a handful of applications -

## Building the project

### Build requirements:

- a Golang installation (built & tested on go v1.23)
- an internet connection to download dependencies (only necessary if dependencies have changed or this is the first build)
- a `make` installation. This project is built with GNU make v4 or higher; full compatibility with other versions of make (such as that shipped by Apple) is not guaranteed, but it _should_ be broadly compatible.

To build the project, simply run `make` in the project's root directory to build the output executable.

> _Note: running with `make` is not strictly necessary. Reference the provided `Makefile` for typical development commands._
