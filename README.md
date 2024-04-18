# VPN

This project provides a wrapper around the `openvpn3` cli for ease of use and session management.

### Prerequisites

- Only tested on Linux
- 1Password cli installed and in your path (`op`)
- `openvpn3` client installed

### Usage

For portability this assumes that you store your .ovpn files in a 1Password vault. You specify the prefix via the `op-prefix` option in `~/vpn/.vpn.yaml`, or it will be requested on first run.

```
vpn up -r <region>
```