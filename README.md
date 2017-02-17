# ./tba
CLI app for getting data from The Blue Alliance.

## Compiling
* Clone the repository and `cd` into it
* Build the executable:

        go build tba.go

* Move it so that it can be executed from anywhere (choose another location in your path if you don't have administrator privileges):

        mv tba /usr/local/bin/tba

## Usage examples
* Get all data on a team:

        tba team -n 254

* Get a specific data point about a team:

        tba team -n 2056 -d country # => Team 2056 is from Canada.

* Get data on an event:

        tba event -k 2013cmp

    (If you omit the year, the current year will be inferred.)

* Get a specific data point about an event:

        tba event -k new -d type # => The Newton Division is a Championship Division (ID 3).

## Licensing
This software is available under the terms of the [MIT License](LICENSE).

## Authors
* [Erik Boesen](https://github.com/ErikBoesen)
