# `./frc`
A helpful CLI app for FRC-related tasks.

* Fetch data from The Blue Alliance
* More features coming soon

![Screenshot](screenshot.png)

## Compiling
* Clone the repository and `cd` into it
* Build the executable:

        go build frc.go

* Move it so that it can be executed from anywhere (choose another location in your path if you don't have administrator privileges):

        mv frc /usr/local/bin/frc

## Usage examples
* Get all data on a team:

        frc team -n 254

* Get a specific data point about a team:

        frc team -n 2056 -d country # => Team 2056 is from Canada.

* Get data on an event:

        frc event -k 2013cmp

    (If you omit the year, the current year will be inferred.)

* Get a specific data point about an event:

        frc event -k new -d type # => The Newton Division is a Championship Division (ID 3).

* To get match data, there are two methods; by key

        frc match -k 2017mokc_qm23

    or by identifiers.

        frc match -y 2017 -e mokc -l qm -n 23

## Licensing
This software is available under the terms of the [MIT License](LICENSE).

## Authors
* [Erik Boesen](https://github.com/ErikBoesen)
