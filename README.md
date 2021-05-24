# Share-Price-Collector

## 1. Summary

A command line program to retrieve share prices from the lse.co.uk website and record them either in a postgres database or to a file.

## 2. Features 

### 2.1. Implemented

 - Read config file (JSON).
 - Read config stored in Postgres.

### 2.2. Partially Implemented

 - Store results in Postgres (optional).

### 2.3. Planned Implementation

 - Output results to screen (optional).
 - Output results to csv file (optional).
 - Output results to json (optional).
 - Update Postgres DB config with JSON config file (optional).
 - Add/remove share from command line.

## 3. Usage

### 3.1. Syntax

Currently, executing the following command will read the config file and display the results to the screen.

```
go run main.go
```

Although there are flags implemented to read the datatabase, database creation scripts are yet to be created and tested.

#### 3.1.1. Flags


|Flag |Description | Implemented? |
--- | --- | ---
|`-dbname`|Database name.| Yes |
|`-user`|Database username| Yes |
|`-passwd`|Database password| Yes |
|`-port`|Database port (default is 5432).| Yes |
|`-output [screen]`|Output to screen.| No |
|`-output [csv]`|Output to screen.| No |
|`-output [json]`|Output to screen.| No |

### 3.2. The JSON config file

This config file only consists of 5 different keys. The initial `startPoll`,  `endPoll` and `interval` keys are defaults. These will be used if any subsequent corresponding keys associated with each individual share are blank.

The keys are:

 - `code` - The 4 character share code that uniquely defines the share.
 - `description` - The description of the share (can be anything you like).
 - `startPoll` - The time to start the polling (24hr format).
 - `endPoll` - The time to end the polling (24hr format).
 - `interval` - The time between each poll (in minutes).

 The format of the file is shown next

```
{
    "startPoll": "08:00",
    "endPoll": "17:00",
    "interval": 5,

    "shareCodes":[
        {
            "code":"SHRA",
            "description": "Share A",
            "startPoll": "08:00",
            "endPoll": "17:00",
            "interval": 5
        },
        {
            "code":"SHRB",
            "description": "Share B",
            "startPoll": "08:30",
            "endPoll": "17:30",
            "interval": 6
        },
        {
            "code":"SHRC",
            "description": "Share C",
            "startPoll": "09:00",
            "endPoll": "18:00",
            "interval": 7
        },
        {
            "code":"SHRD",
            "description": "Share D",
            "startPoll": "09:00",
            "endPoll": "18:00",
            "interval": 7
        }
    ]
 }
 ```