# go-ass5

## About
go-ass5 is a system that imitates a bidding service. 
go-ass5 is capable of running with three clients and three servers. If there were to be made small changes it could handle more of either one.

## How to run

### Start servers

First the servers has to be started.
First write command:
```bash
go run server/server.go
```
If it is the first server then write `1` after running the above command. If it is the second then write `2` and if it is the third write `3` and hit enter.
Then you have to write at which point you want the bidding to stop. This is expressed as the clock you want it to stop. The format is `<HH MM>` followed by hitting enter.
These steps has to be done for all three servers.

### Start client

Then the clients can be started.
Write the command:
```bash
go run client/client.go
```
Then you have to name the client, followed by enter.
This step has to be done for all three clients. Make sure to use a unique name for each client. For the program to create clean log name the client with four characters.

## Crash server
To crash a server you have to write the command `ctrl + c`. This will crash the server entirely.

## System log

### Server log

The server can output four different log types.
The handshake will look like this:
```bash
<time (yyyy/MM/dd HH:mm:ss)> Handshake <client name>
```
<br/>

The bid will look like this:
```bash
<time (yyyy/MM/dd HH:mm:ss)> Bid <client name> <see bid response below> with <bid amount>
```
If the auction is over it will be followed by:
```bash
but auction over, winner: <winner name> , with <winning amount>
```
<br/>

The Request will look like this:
```bash
<time (yyyy/MM/dd HH:mm:ss)> Request <client name> highest bid is: <highest bid> by: <highest bid name>
```
<br/>

When a client bid and the auction has just finished it will log:
```bash
<time (yyyy/MM/dd HH:mm:ss)> --- Auction is over, <winner name> won with bid <winning amount> ---
```

### Client log

The client can output five different log types.
The bid will look like this:
```bash
---------Bid <bid amount> was <see bid response below>
```
<br/>

The request will look like this:
```bash
---------Current highest bid is <highest current bid>
```
<br/>

If a server crashes it will write:
```bash
Server <server port> is down
```
<br/>

If the client has lost the bid it will write:
```bash
Won the auction with bid <bid amount>
```
And if lost:
```bash
Lost the auction
```

## Bid response
The bid response is what the server can answer to a potential bid.
If the bid amount is larger than the current highest bid it will respond `Success`.
If the bid amount is less than or if the auction is over it will respond `Fail`.
If there happens an exception it will respond `Exception`.
