# gossh
A simple CLI-based SSH Client written entirely in go.

## Inspiration
I was working on this project with lots of microservices and only had the option to login to remote servers with password. I wanted to bring my servers with me on any machine I was working on. One way was to use paid apps like Termius. Since I'm poor, I Decided to write my own simpler version that DEFINITELY gets the job done :)

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
You can add your remote servers and connect to them without entering your password everytime.
If you provide ```-e``` flag with a passphrase when adding a new server, Your password will be encrypted.
To access a server which it's password is encrypted, again, pass ```-e``` flag and your passphrase after server's name.

List all available servers:
```$ gossh list```

Add a new server:
```$ gossh add -n [server name[no spaces]] -a [server address] -t [Optional][port [default:22]] -u [user] -p [password] -e [Optional][passphrase to encrypt password with]```

Connect to a server:
```$ gossh connect [server name] -e [Optional][passphrase to decrypt password]```

Delete a server (permanently):
```$ gossh delete [server name]```

### Backup your data
Just copy ```servers.ini``` file from ```/usr/local/etc/gossh``` and put it inside the same directory on any machine.
If you've encrypted your passwords (Explained in Usage section), make sure not to forget your passphrase and you will be good to go.
Simple, Right? :)

## Contributing
Pull requests are welcomed. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://github.com/pooladkhay/gossh/blob/main/LICENSE)