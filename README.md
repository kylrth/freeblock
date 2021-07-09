# freeblock

Block and unblock websites using the `/etc/hosts` hosts file.

You can download the `freeblock` binary by going to the latest release and downloading the one for your platform and architecture. If you don't see your architecture listed, let me know and I can add it.

On Windows, the default hosts file location is `C:\Windows\System32\drivers\etc\hosts`.

## usage

The `freeblock` binary has three subcommands:

- `block` accepts a list of domains to block, and adds them to the hosts file with a `0.0.0.0` IP address so that they don't resolve. See `freeblock block -h` for more details about how it handles domains already present in the file.
- `unblock` accepts a list of domains to unblock. It does this by commenting out any lines that have that domain set to resolve to `0.0.0.0`. Again, see `freeblock unblock -h` for details.
- `open` accepts a list of domains to temporarily unblock. It does the same thing as `unblock` but then waits until it's killed (with either SIGINT or SIGTERM) to re-block the domains. Currently `open` reblocks the domains by restoring the old version of the file, so any changes made to the hosts file while `open` is running will be lost. If someone complains hard enough I'll change this.
