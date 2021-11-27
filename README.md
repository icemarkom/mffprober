# Minimalist Prober for Modern Forms Fans

This came from discussion on
[Home Assistant forum](https://community.home-assistant.io/t/modern-forms-smart-fans-integration/109318/36),
where a few users of Modern Forms fans were finding their fans to stop
responding over time. I wanted to build a simple prober service outside the HA
codebase to see if excessive probing was causing the problem, or perhaps
something else was the cause.

Note: All API commands are from the code-base of the
[Modern Forms integration](https://github.com/jimpastos/ha-modernforms), since
I was unable to find any documented API elsewhere.

## Usage

NOTE: This is assuming you are using the binary release from the
[Releases](https://github.com/icemarkom/mffprober/releases) page.

Basic usage:

```shell
mffprober hostname.or.ip.address
```

There are a few other flags of interest:

- `count`: Maximum number of probes (0 means unlimited).
- ~~`exit-on-error`: DEPRECATED: Use maxfail=0 to disable, or a positive value to control (default true).~~
- ~~`host string`: DEPRECATED: Specify host name/address after the flags.~~
- `interval`: Polling interval in seconds (default 10 seconds).
- `maxfail`: Maximum number of failed probes (0 means unlimited) (default 1).
- `quiet`: Log only polling errors.
- `reboot`: Reboot fan. Ignores most flags.
- `timeout`: Polling probe timeout in seconds (default 1 second).
- `version`: Show version.

## Binary Releases
Binary releases can be found on the
[Releases](https://github.com/icemarkom/mffprober/releases) page.

Currently, the binaries are available for:

- Linux (tgz, apk, rpm, deb)
  - amd64
  - armhf
  - arm64
- Windows
  - 32bit
  - 64bit
- macOS
  - Universal

### Installing Linux binaries using apt

Linux users who use Debian-based systems (Debian, Ubuntu, etc.) can also use
apt to install, and keep up to date with the releases, as with any other package.

Releases are signed with my personal
[GPG key](https://github.com/icemarkom/gpg-key/releases/latest/download/markom@gmail.com.asc).
In order to install packages, you need to add this key for use with apt.

**WARNING:** Most of the instructions on how to do that on the Internet are inherently unsafe,
alas there is no commonly agreed-upon alternative
([1](https://askubuntu.com/a/1307181/1546941),
[2](https://unix.stackexchange.com/questions/332672/how-to-add-a-third-party-repo-and-key-in-debian/582853#582853)).
What follows is my _suggested_ approach.

- Download key to `/usr/local/share/keyrings`
  ```shell
  sudo mkdir -p /usr/local/share/keyrings/
  curl -Ls https://github.com/icemarkom/gpg-key/releases/latest/download/markom@gmail.com.asc \
  | sudo tee /usr/local/share/keyrings/markom@gmail.com.asc
  ```
- Add repository to `/etc/apt/sources.list.d`
  ```shell
  echo 'deb [signed-by=/usr/local/share/keyrings/markom@gmail.com.asc] https://github.com/icemarkom/mffprober/releases/latest/download/ /' \
  | sudo tee /etc/apt/sources.list.d/mffprober.list
  ```
- Update apt
  ```shell
  sudo apt update
  ```
- Install `mffprober`
  ```shell
  sudo apt install mffprober
  ```
