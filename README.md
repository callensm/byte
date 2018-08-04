# Byte

## Installation

You can download the `byte` binary file from the releases page for the current version and place in your `PATH`, or to install from source:

```sh
$ go get -u github.com/callensm/byte
$ cd $GOPATH/src/github.com/callensm/byte && make
```

## Usage
### *send*
```sh
byte send -a [<IP>|<IP>:<PORT>] -s [<PATH_TO_CONTENT>]
```

| Option | Alias | Description                                                                         | Example                                    |
| :----: | :---: | :---------------------------------------------------------------------------------: | :----------------------------------------: |
| --addr | -a    | Set the address of the receiving socket. If no port is specified, default is `4500` | `127.0.0.1`, `127.0.0.1:3000`              |
| --src  | -s    | Specify the relative path to the file or directory of files to be sent              | `~/Desktop/myfiles/`, `~/Desktop/data.txt` |

### *receive*
```sh
byte receive -d [<PATH_TO_DEST_DIR>] -p [PORT?] -a
```

| Option         | Alias | Description                                                                                   | Example                                    |
| :------------: | :---: | :-------------------------------------------------------------------------------------------: | :----------------------------------------: |
| --dir          | -d    | Set the path to the destination directory for files to be written to                          | `~/Desktop/destination/`                   |
| --port         | -p    | Set the port to listen for incoming connections on. If not given, default port is `4500`      | `~/Desktop/myfiles/`, `~/Desktop/data.txt` |
| --auto-approve | -a    | Whether to require approval of file structure being sent. Defaults to `false` if not provided | -                                          |

### *view*
```sh
byte view -r [<PATH_TO_ROOT>]
```

| Option | Alias | Description                                          | Example                  |
| :----: | :---: | :--------------------------------------------------: | :----------------------: |
| --root | -r    | Set the path to the root directory for the file tree | `~/Desktop/destination/` |

## License
[MIT](./LICENSE)
