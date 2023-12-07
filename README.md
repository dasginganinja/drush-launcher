# Drush Launcher - Reviving Drush Support with 💪 and ❤️

## Introduction
🚀 Welcome to `drush-launcher`, a heartfelt tribute and robust successor to the original drush-launcher. This Go-based solution breathes new life into Drush support, especially for the long-awaited Drush 12.

## Why `drush-launcher`?
- 💡 **Drush 12 Compatibility**: Unlike its predecessor, `drush-launcher` embraces Drush 12, ensuring your projects stay up-to-date and efficient.
- 🌍 **System-wide Installation**: Seamlessly integrates into your $PATH, providing universal access to the `drush` executable.
- 🌟 **Community-Driven**: Born from community needs and nurtured by proactive problem-solving.


## Step 1 - Get Binary for drush-launcher
Simple, system-wide setup ensures `drush` is always at your fingertips. Just a few commands, and you're all set!

### Option a: Fetch Built Release
Installation is as simple as running the following command:

`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/dasginganinja/drush-launcher/main/download_and_build.sh)"`

This will download a prebuilt file for your system based on the detected OS and Architecture and will be available at `drush-launcher_*` in the same directory.

### Option b: Build release from source
Run `go build` in the project directory. 

This will create an executable called `drush-launcher`.

> For more details or different OS/platform targets, refer to the [Go build documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies). These can be controlled via the `$GOARCH` and `$GOOS` environment variables when running the build.

## Step 2 - Install into $PATH

### System wide (as the original Drush launcher was)

It is recommended to install this binary to `/usr/local/bin/drush` so this is the preferred fallback `drush` executable in the $PATH.

As an example, one could run this to copy the file.
`sudo cp -n drush-launcher /usr/local/bin/drush`

> If drush launcher already exists at the location it will not be overwritten with the above command because `-n` means to not overwrite the file if it exists.

## Usage

For standard usage, just run this as you normally would drush.

This executable supports one flag for setting the Drupal root via `-r` or `--root`, in case it wasn't detected automatically in the current or parent directories.

All other flags and arguments are passed on, including the aforementioned root options.

## Testing
Execute `go test` to run tests. 

## Contributing
Your feedback and merge requests are highly valued! If you have ideas or improvements that align with the project's purpose, please feel welcome to contribute. Together, we can enhance `drush-launcher` for everyone!

---

Discover more and contribute on our [GitHub repository](https://github.com/dasginganinja/drush-launcher).