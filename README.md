<div align="center">
    <img src="fsguard.svg" alt="FsGuard Wrapper logo" width="200">
    <p>FsGuard is a tool for verifying filesystem integrity at boot time.</p>
</div>

# FsGuard
## Building
Dependencies:
- go

simply run `build.sh` to build the project and append the testing signature.

This will ensure that FsGuard is able to run properly by fetching signatures


## Deploying
### Filelist
FsGuard needs a filelist containg the sha1sum and suid permission of every binary to scan, an example file can be found [here](https://github.com/linux-immutability-tools/FsGuard/blob/main/test_filelist).
A bash oneliner to create an entry for this file could look like this:
```
echo $(sha1sum /path/to/binary | sed 's/  / /g') $(ls -al /path/to/binary | awk 'BEGIN{FS=" "}; {print $1};' | grep s > /dev/null && echo "true" || echo "false")
```

This Filelist can be placed anywhere, as long as FsGuard has access to it when it launches.

### Signing the Filelist
FsGuard expects a minisign signature and filelist to be appended to the binary. An example signature "set" can be found [here](https://github.com/linux-immutability-tools/FsGuard/blob/main/signatures).
A signature set can be generated and added to FsGuard with these commands:
```bash
# Create a new passwordless key pair
minisign -WG
# Signing the filelist
minisign -Sm /path/to/filelist

# Generate the signature set
touch /path/to/signature
echo -n "----begin attach----" >> /path/to/signature
cat /path/to/filelist.minisig >> /path/to/signature
echo -n "----begin second attach----" >> /path/to/signature
tail -n1 ./minisign.pub >> /path/to/signature

# Append the signature set to the FsGuard binary
cat /path/to/signature >> /path/to/FsGuard
```

## Launching FsGuard
### As an init
FsGuard automatically starts the verification if it detects that it is a specific binary. 
This binary name and path can be set with the `InitLocation` property in the [`config/config.go`](https://github.com/linux-immutability-tools/FsGuard/blob/main/config/config.go) file.

Additionally, FsGuard automatically starts a proper init once it completed the verification process. The init it launches can be controlled with the `PostInitExec` property in [`config/config.go`](https://github.com/linux-immutability-tools/FsGuard/blob/main/config/config.go).

### As a pre-init script
FsGuard can also be started in a bash script that gets launched as an init, in this case, FsGuard works like a regular cli application and accepts the filelist location as an argument.
A possible pre-init script could look like this:
```
#!/usr/bin/bash
FsGuard verify /path/to/filelist
exec /path/to/init
```

Make sure to launch the proper init using `exec`, some init systems like systemd will refuse to launch if they are not pid1, `exec` makes sure that the init script "drops" its pid and systemd is able to claim it.

## Reporting issues
When reporting issues you encounter with FsGuard, please make sure to include the config.go file and how FsGuard gets launched.
