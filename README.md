# chatProject
little project for GDG Marbella Go course.

## Dependencies

`go get github.com/marcusolsson/tui-go`

## Usage

### Server

Go to  server folder and execute 
```bash
$ go run server.go
```

After that you will see in the cosole something like this, if everything works fine ;-) 

```bash
2020/07/08 15:33:01 Chat server started and listening on :8887
```

### Client

Go to `client` folder, and use the below commmand to launch clients:

```bash
$ go run client.go <username>
``` 
It is important to use a username, if not the program will fail. Error handling is very very poor for now.

Of course you can launch severals clients.

## Features

* TCP chat client using a GUI console module (tui-go).
* Chat with more than one user.
* Notify new users.
* Send private message. To do this, use the below sintax: `@<username> message`. This message only will appear into <username> chat client.

