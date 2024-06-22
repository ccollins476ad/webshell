# webshell

webshell is a server that listens for http requests from localhost, translates them to commands, and executes the result. By default, it listens for http requests on port 9901.

The intent is to enable browser scripts to execute arbitrary commands. Such scripts are prohibited from executing commands directly, but are allowed to send GET http requests.

## Example

1. Start webshell
```
$ ./webshell
```

2. In another terminal, ask webshell for the current UTC time
```
$ curl localhost:9901/webshell/v1/date/-u
Sun Oct 19 22:49:36 UTC 2025
```

## Request format

webshell accepts http requests of any method. A request encodes the command to execute in the url path, as follows:

```
/webshell/v1/<command>/<arg1>/<arg2>/...
```

The command and arguments are url-escaped. The request body is empty.
