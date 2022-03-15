# gossh
A ~~simple~~ CLI-based SSH Client written entirely in go.

## Inspiration
I was working on this project with lots of microservices and only had the option to login to remote servers with password. I wanted to bring my servers with me on any machine I was working on. One way was to use paid apps like Termius. Since I'm poor, I Decided to write my own ~~simpler~~ version that DEFINITELY gets the job done!

## Installation
1) Clone the repo
2) Build the binary
3) Move the binary to ```/usr/local/bin``` in order to run it from everywhere
4) Finally, Make is executeable
```bash
$ git clone https://github.com/pooladkhay/gossh
$ cd gossh
$ go build .
$ mv gossh /usr/local/bin/gossh
$ chmod 755 /usr/local/bin/gossh
```

## Usage
```
$ gossh [command]

Available Commands:
  add         Adds a new server to the list
  completion  generate the autocompletion script for the specified shell
  connect     Connects to a specific server
  delete      Deletes the specified server from server's list
  help        Help about any command
  list        Lists all available servers

Flags:
  -h, --help      help for gossh
  -v, --version   version for gossh

Use "gossh [command] --help" for more information about a command.
```

### Backup your data
Just copy ```servers.ini``` file from ```~/.gossh``` and put it inside the same directory on any machine.
If you've encrypted your passwords (Explained in Usage section), make sure not to forget your passphrase and you will be good to go.
Simple, Right? :)

## Contributing
Pull requests are welcomed. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://github.com/pooladkhay/gossh/blob/main/LICENSE)
