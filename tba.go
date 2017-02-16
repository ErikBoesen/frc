package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "regexp"

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

        // Fetch team data
        team, _ := tba.GetTeam(*teamNumber)

        // Trim URL fragments from start of event website URL
        team.Website = urlRE.ReplaceAllString(team.Website, "")

        // If user didn't specify a data point, output all of them.
        if *teamDataPoint == "" {
            fmt.Printf("\n    Team %d:\n", *teamNumber)
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
                fmt.Printf("Team %d's name: %s\n", *teamNumber, team.Name)
            case "website":
                fmt.Printf("Team %d's website: %s\n", *teamNumber, team.Website)
        	case "locality":
                fmt.Printf("Team %d's locality: %s\n", *teamNumber, team.Locality)
        	case "rookieyear":
                fmt.Printf("Team %d's first year was: %d\n", *teamNumber, team.RookieYear)
        	case "region":
                fmt.Printf("Team %d plays in the %s region.\n", *teamNumber, team.Region)
            case "teamnumber":
                fmt.Printf("Team %d's number is %d.\n", *teamNumber, team.TeamNumber)
        	case "location":
                fmt.Printf("Team %d comes from %s.\n", *teamNumber, team.Location)
        	case "key":
                fmt.Printf("Team %d's key is %s.\n", *teamNumber, team.Key)
        	case "countryname":
                fmt.Printf("Team %d is from the following country: %s\n", *teamNumber, team.CountryName)
        	case "motto":
                fmt.Printf("Team %d's motto is \"%s\".\n", *teamNumber, team.Motto)
        	case "nickname":
                fmt.Printf("Team %d's nickname is %s.\n", *teamNumber, team.Nickname)
            default:
                fmt.Printf("Unrecognized data point \"%s\".\n", *teamDataPoint)
            }
        }
    }
    if eventCommand.Parsed() {

        event, _ := tba.GetEvent(*eventKey)

        if event.Key == "" {
            fmt.Printf("\n    Invalid event key \"%s\".\n\n", *eventKey)
            os.Exit(1)
        }

        // Trim URL fragments from start of event website URL
        event.Website = urlRE.ReplaceAllString(event.Website, "")

        event.StartDate = strings.Replace(event.StartDate, "-", "/", -1)
        event.EndDate = strings.Replace(event.EndDate, "-", "/", -1)

        if *eventDataPoint == "" {
            fmt.Printf("\n    %d %s (%s):\n", event.Year, event.Name, event.Key)
            //fmt.Printf("\tShortname: %s\n", event.ShortName)
            if event.StartDate == event.EndDate {
                fmt.Printf("\tDate: %s\n", event.StartDate)
            } else {
                fmt.Printf("\tDates: %s - %s\n", event.StartDate, event.EndDate)
            }
            fmt.Printf("\tTimezone: %s\n", event.Timezone)
            fmt.Printf("\tWebsite: %s\n", event.Website)
            if event.EventDistrict > 0 {
                fmt.Printf("\tDistrict: %s (ID %d)\n", event.EventDistrictString, event.EventDistrict)
            }
            fmt.Printf("\tLocation: %s\n", event.Location)
            fmt.Printf("\tAddress: %s\n", strings.Replace(event.VenueAddress, "\n", ", ", -1))
            fmt.Printf("\tEvent Type: %s (ID %d)\n\n", event.EventTypeString, event.EventType)

            /*Official            bool          `json:"official"`
            ShortName           string        `json:"short_name"`
            FacebookEid         interface{}   `json:"facebook_eid"`
            Webcast             []interface{} `json:"webcast"`
            Alliances           []struct {
                Declines []interface{} `json:"declines"`
                Picks    []string      `json:"picks"`
            } `json:"alliances"`*/
        } else {
            fmt.Println("Event data point fetching coming soon, for now just get the whole event.")
        }
    }
}
