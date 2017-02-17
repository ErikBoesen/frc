package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"
    "regexp"
    "time"

    "github.com/ErikBoesen/tba"
)

const (
    VERSION = "0.1.0"
)


func main() {
    // Subcommands
    teamCommand := flag.NewFlagSet("team", flag.ExitOnError)
    eventCommand := flag.NewFlagSet("event", flag.ExitOnError)

    // Team subcommand flags.
    teamNumber := teamCommand.Int("n", 0, "Team number. (Required)")
    teamDataPoint := teamCommand.String("d", "", "Data point to display. If unspecified, all team data will be shown.")

    // Event subcommand flags.
    eventKey := eventCommand.String("k", "", "ID of event you want data on. (Required)")
    eventDataPoint := eventCommand.String("d", "", "Data point to display. If unspecified, all event data will be shown.")

    // Verify a subcommand has been provided.
    if len(os.Args) < 2 {
        fmt.Println("Error: subcommand is required.")
        os.Exit(1)
    }

    // Initialize TBA parser
    tba, _ := tba.Init("erikboesen", "tbacli", VERSION)

    // Regex to strip URL junk
    urlRE, _ := regexp.Compile("https?://")

    // Parse flags for appropriate FlagSet
    // FlagSet.Parse() requires a set of arguments to parse as input
    // os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
    switch os.Args[1] {
    case "team":
        teamCommand.Parse(os.Args[2:])
    case "event":
        eventCommand.Parse(os.Args[2:])
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Check which subcommand was Parsed using the FlagSet.Parsed() function; then handle each case accordingly.
    // FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user didn't provide any flags).
    if teamCommand.Parsed() {
        // Team key
        tk := *teamNumber

        // Fetch team data
        team, _ := tba.GetTeam(tk)

        // TODO: Catch nonexistent teams

        // Trim URL fragments from start of event website URL
        team.Website = urlRE.ReplaceAllString(team.Website, "")

        // If user didn't specify a data point, output all of them.
        if *teamDataPoint == "" {
            fmt.Printf("\n    Team %d:\n", tk)
            fmt.Printf("\tNickname: %s\n", team.Nickname)
            //fmt.Printf("\tName/Sponsors: %s\n", team.Name)
            fmt.Printf("\tWebsite: %s\n", team.Website)
            fmt.Printf("\tLocality: %s\n", team.Locality)
            fmt.Printf("\tRookie Year: %d\n", team.RookieYear)
            fmt.Printf("\tRegion: %s\n", team.Region)
            fmt.Printf("\tLocation: %s\n", team.Location)
            fmt.Printf("\tCountry: %s\n", team.CountryName)
            if string(team.Motto[0]) == "\"" {
                fmt.Printf("\tMotto: %s\n\n", team.Motto)
            } else {
                fmt.Printf("\tMotto: \"%s\"\n\n", team.Motto)
            }
        } else {
            switch strings.ToLower(*teamDataPoint) {
            case "name":
                fmt.Printf("\n    Team %d's name: %s\n\n", team.TeamNumber, team.Name)
            case "website":
                fmt.Printf("\n    Team %d's website: %s\n\n", team.TeamNumber, team.Website)
        	case "locality":
                fmt.Printf("\n    Team %d is based in %s.\n\n", team.TeamNumber, team.Locality)
        	case "rookieyear":
                fmt.Printf("\n    Team %d's first year was %d.\n\n", team.TeamNumber, team.RookieYear)
        	case "region":
                fmt.Printf("\n    Team %d plays in the %s region.\n\n", team.TeamNumber, team.Region)
            case "teamnumber":
                fmt.Printf("\n    Team %d's number is %d.\n\n", team.TeamNumber, team.TeamNumber)
        	case "location":
                fmt.Printf("\n    Team %d is based in %s.\n\n", team.TeamNumber, team.Location)
        	case "key":
                fmt.Printf("\n    Team %d's key is %s.\n\n", team.TeamNumber, team.Key)
        	case "country":
                fmt.Printf("\n    Team %d is from %s.\n\n", team.TeamNumber, team.CountryName)
        	case "motto":
                if string(team.Motto[0]) == "\"" {
                    fmt.Printf("\n    Team %d's motto is %s.\n\n", team.TeamNumber, team.Motto)
                } else {
                    fmt.Printf("\n    Team %d's motto is \"%s\".\n\n", team.TeamNumber, team.Motto)
                }
        	case "nickname":
                fmt.Printf("\n    Team %d's nickname is %s.\n\n", team.TeamNumber, team.Nickname)
            default:
                fmt.Printf("\n    Invalid data point \"%s\".\n\n", *teamDataPoint)
            }
        }
    }
    if eventCommand.Parsed() {
        ek := *eventKey
        if _, err := strconv.Atoi(string((*eventKey)[0])); err != nil {
            ek = fmt.Sprintf("%d%s", time.Now().Year(), *eventKey)
        }

        event, _ := tba.GetEvent(ek)

        if event.Key == "" {
            fmt.Printf("\n    Invalid event key \"%s\".\n\n", ek)
            os.Exit(1)
        }

        // Trim URL fragments from start of event website URL
        event.Website = urlRE.ReplaceAllString(event.Website, "")

        event.StartDate = strings.Replace(event.StartDate, "-", "/", -1)
        event.EndDate = strings.Replace(event.EndDate, "-", "/", -1)

        event.VenueAddress = strings.Replace(event.VenueAddress, "\n", ", ", -1)

        if *eventDataPoint == "" {
            fmt.Printf("\n    %d %s (%s):\n\n", event.Year, event.Name, event.Key)
            //fmt.Printf("\tShortname: %s\n", event.ShortName)
            if event.StartDate == event.EndDate {
                fmt.Printf("\tDate: %s\n", event.StartDate)
            } else {
                fmt.Printf("\tDates: %s - %s\n", event.StartDate, event.EndDate)
            }
            if event.Official {
                fmt.Println("\tOfficial: Yes")
            } else {
                fmt.Println("\tOfficial: No")
            }
            fmt.Printf("\tTimezone: %s\n", event.Timezone)
            fmt.Printf("\tWebsite: %s\n", event.Website)
            if event.EventDistrict > 0 {
                fmt.Printf("\tDistrict: %s (ID %d)\n", event.EventDistrictString, event.EventDistrict)
            }
            fmt.Printf("\tLocation: %s\n", event.Location)
            fmt.Printf("\tAddress: %s\n", event.VenueAddress)
            fmt.Printf("\tEvent Type: %s (ID %d)\n\n", event.EventTypeString, event.EventType)
        } else {
            switch strings.ToLower(*eventDataPoint) {
            case "key":
                fmt.Printf("\n    Event %s's key is %s.\n\n", event.Name, event.Key)
            case "name":
                fmt.Printf("\n    Event %s's full name is %s.\n\n", event.Name, event.Name)
            case "shortname":
                fmt.Printf("\n    Event %s's shortname is %s.\n\n", event.Name, event.ShortName)
            case "website":
                fmt.Printf("\n    Event %s's website can be found at %s\n\n", event.Name, event.Website)
            case "date":
            case "dates":
                if event.StartDate == event.EndDate {
                    fmt.Printf("\tThe %s takes place on %s.\n", event.Name, event.StartDate)
                } else {
                    fmt.Printf("\tThe %s takes place from %s to %s.\n", event.Name, event.StartDate, event.EndDate)
                }
            case "startdate":
                fmt.Printf("\n    The %s starts on %s.\n\n", event.Name, event.StartDate)
            case "enddate":
                fmt.Printf("\n    The %s ends on %s.\n\n", event.Name, event.EndDate)
            case "official":
                if event.Official {
                    fmt.Printf("\n    The %s is an official FIRST event.\n\n", event.Name)
                } else {
                    fmt.Printf("\n    The %s is not an official FIRST event.\n\n", event.Name)
                }
        	case "district":
                fmt.Printf("\n    The %s is part of the %s district (ID %d).\n\n", event.Name, event.EventDistrictString, event.EventDistrict)
        	case "location":
                fmt.Printf("\n    The %s is based in %s.\n\n", event.Name, event.Location)
            case "year":
                fmt.Printf("\n    The %s took place in %d.\n\n", event.Name, event.Year)
            case "timezone":
                fmt.Printf("\n    The %s is in the %s timezone.\n\n", event.Name, event.Timezone)
            case "address":
                fmt.Printf("\n    The %s took place at %s.\n\n", event.Name, event.VenueAddress)
            case "type":
                fmt.Printf("\n    The %s is a %s (ID %d).\n\n", event.Name, event.EventTypeString, event.EventType)
            default:
                fmt.Printf("\n    Invalid data point \"%s\".\n\n", *eventDataPoint)
            }
        }
    }
}
