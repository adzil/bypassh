# BypaSSH: Remote SSH using WSL from Visual Studio Code

BypaSSH forwards any SSH command from Windows to WSL. It is intended to be used by the Visual Studio Code's Remote-SSH extension to remotely connect from WSL's OpenSSH which is currently not supported natively.

BypaSSH works by translating Windows-style path in SSH arguments into WSL-style path. It also forwards standard I/O pipe and interrupt signal between BypaSSH parent process and WSL's OpenSSH.

Using WSL's OpenSSH enable us to seamlessly share SSH configuration (and private key) between the Visual Studio Code which runs in Windows and other CLI applications that utilizes SSH and runs in WSL.

## Quick Start Guide

Install the latest BypaSSH binary using `go install`:

```
go install github.com/adzil/bypassh@latest
```

Then, add the following line to your VSCode configuration JSON file:

```js
{
    // ... Other config goes here ...
    
    "remote.SSH.path": "C:\\Path\\to\\Your\\Go\\bin\\bypassh.exe",
    
    // Put this if your SSH config is located in WSL.
    "remote.SSH.configFile": "\\\\wsl$\\Ubuntu\\home\\me\\.ssh\\config",
    
    // Or put this instead if your SSH config is located in Windows.
    "remote.SSH.configFile": "C:\\Users\\Me\\.ssh\\config"


    // ... Other config also goes here...
}
```

Now you can proceed to connect to your remote server through the Remote-SSH extension.

## Configuration

If you have different distro name than `Ubuntu`, you need to make `bypassh.json` file in the same directory as the `bypassh.exe` binary and set the `distro` field:

```json
{
    "distro": "your-distro-name"
}
```

There's also configurable path for WSL and SSH binary (see `bypassh.exe -h` for more info) but for most cases you probably wouldn't have to touch it.

## Known Bugs

The Remote-SSH extension is stuck on establishing a connection when the remote host is not yet known even though the user already presses the "yes" option. This issue can be fixed by manual SSH through CLI, permanently add it to the known host list, and then try to reconnect again.

## WSL1 Support

BypaSSH **will not work** with WSL1 because it requires WSL2 network mount from Windows (e.g. `\\wsl$\Ubuntu\`) that did not available on earlier version of WSL.
