# Spin Version Manager Plugin

This is a plugin that makes it easy to switch between different versions of Spin.

# Installation

The trigger is installed as a Spin plugin. It can be installed from a release or build.

## Install the latest version of the plugin

The latest stable release of the command trigger plugin can be installed like so:

```sh
spin plugins update
spin plugin install verman
```

## Install the canary version of the plugin

The canary release of the command trigger plugin represents the most recent commits on `main` and may not be stable, with some features still in progress.

```sh
spin plugins install --url https://github.com/fermyon/verman-plugin/releases/download/canary/verman.json
```

## Install from a local build

Alternatively, use the `spin pluginify` plugin to install from a fresh build. This will use the pluginify manifest (`spin-pluginify.toml`) to package the plugin and proceed to install it:

```sh
spin plugins install pluginify
cargo build --release
spin pluginify install
```

# Usage

Once the plugin is installed, you can try the below commands:

## Set a different version of Spin

```sh
# Adding the v prefix to the version is optional
spin verman set v2.5.0
```

## List the versions of Spin that are downloaded via the verman plugin

```sh
spin verman ls
```

## Remove version(s) of Spin downloaded via the verman plugin

Remove a single version:

```sh
# Adding the v prefix to the version is optional
spin verman rm v2.5.0
```

Remove all versions:

```sh
spin verman rm all
```

Remove the alternate Spin version, reverting back to the root version of Spin, but preserving all other versions of Spin downloaded locally:

```sh
spin verman rm current
```