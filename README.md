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

Basic usage: 

```shell
cd main
go run main.go --host={{ip_address}}
```

There are a few other flags of interest:

- `interval`: How often should the prober attempt to connect to the fan. Value
  is in seconds, and the default is 10 seconds.
- `exit-on-error`: Leave the polling loop on first error (of any kind). If not
  set, prober will continue attempting to talk to the fan, even if it receives
  errors. Defaults to `true`.
- `timeout`: HTTP client timeout in seconds. Defaults to 1.
- `quiet`: Log only probing and other errors. No other output is produced.
  Defaults to `false`.
