# Spin Version Manager Plugin

This is a plugin that makes it easy to switch between different versions of Spin.

# Installation

## Install the latest version of the plugin

The latest stable release of the verman plugin can be installed like so:

```sh
spin plugins update
spin plugin install verman
```

## Install the canary version of the plugin

The canary release of the verman plugin represents the most recent commits on `main` and may not be stable, with some features still in progress.

```sh
spin plugins install --url https://github.com/fermyon/verman-plugin/releases/download/canary/verman.json
```

## Install from a local build

Alternatively, use the `spin pluginify` plugin to install from a fresh build. This will use the pluginify manifest (`spin-pluginify.toml`) to package the plugin and proceed to install it:

```sh
spin plugins install pluginify
go build -o verman main.go
spin pluginify --install
```

# Usage

Once the plugin is installed, you'll need to prepend the verman `current_version` directory to your path, allowing the verman plugin to temporarily override your current version of Spin:

```sh
# To permanently prepend the "current_version" directory to your path, add this command to your .zshrc/.bashrc file.
export PATH="$HOME/.spin_verman/versions/current_version:$PATH"
```

Once the path is prepended, you can try the below commands:

## List available versions of Spin

```sh
# list all available versions of spin
spin verman list-remote
```

## Download a specific version of Spin

Specify the desired version:

```sh
# You can download multiple versions at once
spin verman get canary v2.5.0 2.7.0
```

Get the latest stable version:

```sh
spin verman get latest
```

## Create an alias for a local build of Spin

Specify the alias name and path to Spin binary:

```sh
spin verman alias myalias /path/to/spin
```

## Set a different version of Spin

Set a specific version:

```sh
# Adding the v prefix to the version is optional
spin verman set v2.5.0
```

Set the latest version:

```sh
spin verman set latest
```

Set to an alias for a local build:

```sh
spin verman set myalias
```

## Using `.spin-version` to Download and Set the desired Spin version

You can specify the desired version of Spin in a `.spin-version` file. The `verman` plugin is able to download and set the current version from a `.spin-version` file:

```bash
# Specify desired version in a .spin-version file
echo "2.7.0" > .spin-version

# Download the spin version specified in .spin-version
spin verman get

# Set the version of spin specified in .spin-version
spin verman set
```

*Note: Arguments are provided to either `spin verman get` or `spin verman set` have higher priority compared to `.spin-version`.*

## Update a version of Spin in the `~/.spin_verman` directory

```sh
# "canary" is currently the only subcommand supported by "update"
spin verman update canary
```

## List the versions of Spin that are downloaded via the verman plugin

```sh
spin verman list
```

## Remove a version of Spin downloaded via the verman plugin

Remove a single version:

```sh
# Adding the v prefix to the version is optional
spin verman remove v2.5.0
```

Remove an alias for a local build:

```sh
spin verman rm myalias
```

Remove all versions:

```sh
spin verman remove all
```

Remove the alternate Spin version, reverting back to the root version of Spin, but preserving all other versions of Spin downloaded locally:

```sh
spin verman remove current
```