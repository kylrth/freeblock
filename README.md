# freeblock

Block and unblock websites using the `/etc/hosts` hosts file (`C:\Windows\System32\drivers\etc\hosts` on Windows).

Freeblock blocks websites by adding entries with a `0.0.0.0` IP address so that they don't resolve.

## installation

Go to [Releases](https://github.com/kylrth/freeblock/releases) and download the latest release for your platform and architecture to somewhere on your PATH, e.g.:

```sh
wget 'https://github.com/kylrth/freeblock/releases/download/v1.1.0/freeblock-linux-amd64' -O - | sudo tee /usr/bin/freeblock > /dev/null
```

## usage

The `freeblock` binary has three subcommands:

- `block` accepts a list of domains to block. See `freeblock block -h` for more details about how it handles domains already present in the file.
- `unblock` accepts a list of domains to unblock. It does this by commenting out any lines that have that domain set to resolve to `0.0.0.0`. Again, see `freeblock unblock -h` for details.
- `open` accepts a list of domains to temporarily unblock. It does the same thing as `unblock` but then waits until it's killed (with either SIGINT or SIGTERM) to re-block the domains. Currently, `open` re-blocks the domains by restoring the old version of the file, so any changes made to the hosts file while `open` is running will be lost.

### time ranges

If you add a comment to a line in `/etc/hosts` like this:

```hosts
0.0.0.0 www.reddit.com  #freeblock:09-17
```

freeblock will not unblock Reddit between 9am and 5pm. Only hours are supported, not minutes.
