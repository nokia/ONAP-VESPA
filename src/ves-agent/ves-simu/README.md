# VES Collector simulator
This is a basic implementation of a VES-Collector simulator. It's used to receive events and validate them against the schema.

Currently, only VES-5.3 schema specification is supported.

The simulator alos has a specific control REST API. Commands to be sent to VES-Agent can be set, and the 
simulator stores all the received event in memory so that they can bes retreived using control API.

## Installing
### From binaries (to be completed)

### From sources (to be completed)

## Usage
### Start the simulator
The simulator is started by invoking the **ves-simu** executable binary. It works out of the box with default parameters, but can be customized using flags :
```
./ves-simu -help
Usage of ./ves-simu:
  -cert string
        Path to certificate file if HTTPS is enabled
  -d uint
        Verbosity level. 0, 1 or 2
  -event-size int
        Maximum size in bytes allowed for events (default 2000000)
  -https
        Enable HTTPS instead of plain HTTP
  -jsonlog
        Format log messaegs into JSON
  -key string
        Path to certificate's private key file if HTTPS is enabled
  -logfile string
        Path to a file where to write logs outputs, additionally to standard output
  -max-event int
        Maximum number of events to keep in memory, 0 meaning there's no limit (default 4096)
  -passwd string
        The password for authenticating incoming requests (default "pass")
  -port int
        The port to bind VES simulator to (default 8443)
  -topic string
        An optional topic
  -user string
        The username for authenticating incoming requests (default "user")
```
There's no log file, and everything is printed on standard output

EG:
```bash
./ves-simu -user foo -passwd bar -port 8080
```

> Once started, the simulator is stopped with the **INTERRUPT** signal (eg: CTRL+C)

### Control REST API
You can inject commands to be sent to VES clients with requests like
```http
POST https://localhost:8443/testControl/v5/commandList HTTP/1.1
Content-Type: application/json

{
    "commandList": [
        {
            "command": {
                "commandType": "heartbeatIntervalChange",
                "heartbeatInterval": 15
            }
        },
        {
            "command": {
                "commandType": "measurementIntervalChange",
                "measurementInterval": 30
            }
        }
    ]
}
```