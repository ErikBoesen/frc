package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"
    "regexp"
    "time"
    "github.com/fatih/color"

    "github.com/ErikBoesen/tba-go"
)

const (
    VERSION = "0.1.0"
)


func main() {
    // Subcommands
    teamCommand := flag.NewFlagSet("team", flag.ExitOnError)
    eventCommand := flag.NewFlagSet("event", flag.ExitOnError)
    matchCommand := flag.NewFlagSet("match", flag.ExitOnError)

    // Team subcommand flags
    teamNumber := teamCommand.Int("n", 0, "Team number. (Required)")
    teamDatum := teamCommand.String("d", "", "Data point to display. If unspecified, all team data will be shown.")

    // Event subcommand flags
    eventKey := eventCommand.String("k", "", "ID of event you want data on. (Required)")
    eventDatum := eventCommand.String("d", "", "Data point to display. If unspecified, all event data will be shown.")

    // Match subcommand flags
    matchKey := matchCommand.String("k", "", "Match key.")
    matchYear := matchCommand.Int("y", time.Now().Year(), "Year in which match took place.")
    matchEvent := matchCommand.String("e", "", "Event at which match occurred.")
    matchLevel := matchCommand.String("l", "", "Event level. Valid choices include 'qm', 'qf', 'sf', and 'f'.")
    matchNumber := matchCommand.Int("n", 0, "Match number.")
    matchRound := matchCommand.Int("r", 0, "Match round (only in playoffs).")
    matchDatum := matchCommand.String("d", "", "Specific datum to fetch.")

    // Verify a subcommand has been provided.
    if len(os.Args) < 2 {
        fmt.Println("Error: subcommand is required.")
        os.Exit(1)
    }


    c := color.New(color.FgCyan, color.Underline)
    b := color.New(color.FgBlue) // TODO: Bold for team numbers?
    r := color.New(color.FgRed)
    g := color.New(color.FgGreen)

    // Initialize TBA parser
    tba, _ := tba.Init("erikboesen", "frcli", VERSION)

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
    case "match":
        matchCommand.Parse(os.Args[2:])
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

        // If user didn't specify a datum, output all of them.
        if *teamDatum == "" {
            fmt.Printf("\n    ")
            c.Printf("Team %d:\n", tk)
            g.Print("\tNickname:     ")
            fmt.Println(team.Nickname)
            //g.Printf("\tFull Name:   ")
            //fmt.Println(team.Name)
            g.Print("\tWebsite:      ")
            fmt.Println(team.Website)
            g.Print("\tLocality:     ")
            fmt.Println(team.Locality)
            g.Print("\tRookie Year:  ")
            fmt.Println(team.RookieYear)
            g.Print("\tRegion:       ")
            fmt.Println(team.Region)
            g.Print("\tLocation:     ")
            fmt.Println(team.Location)
            g.Print("\tCountry:      ")
            fmt.Println(team.CountryName)
            g.Print("\tMotto:        ")
            if string(team.Motto[0]) == "\"" {
                fmt.Println(team.Motto)
            } else {
                fmt.Printf("\"%s\"", team.Motto)
            }
            fmt.Println("\n")
        } else {
            switch strings.ToLower(*teamDatum) {
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
                fmt.Printf("\n    Invalid datum \"%s\".\n\n", *teamDatum)
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

        if *eventDatum == "" {
            fmt.Printf("\n    ")
            c.Printf("%d %s (%s):\n", event.Year, event.Name, event.Key)
            if event.StartDate == event.EndDate {
                g.Print("\n\tDate: ")
                fmt.Printf("%s\n", event.StartDate)
            } else {
                g.Print("\tDates: ")
                fmt.Printf("%s - %s\n", event.StartDate, event.EndDate)
            }
            g.Print("\tOfficial: ")
            if event.Official {
                fmt.Println("Yes")
            } else {
                fmt.Println("No")
            }
            g.Printf("\tTimezone: ")
            fmt.Println(event.Timezone)
            g.Printf("\tWebsite: ")
            fmt.Println(event.Website)
            if event.EventDistrict > 0 {
                g.Print("\tDistrict: ")
                fmt.Printf("%s (ID %d)\n", event.EventDistrictString, event.EventDistrict)
            }
            g.Print("\tLocation: ")
            fmt.Println(event.Location)
            g.Print("\tAddress: ")
            fmt.Println(event.VenueAddress)
            g.Print("\tEvent Type: ")
            fmt.Printf("%s (ID %d)\n\n", event.EventTypeString, event.EventType)
        } else {
            switch strings.ToLower(*eventDatum) {
            case "key":
                fmt.Printf("\n    The %s's key is %s.\n\n", event.Name, event.Key)
            case "name":
                fmt.Printf("\n    The event's full name is %s.\n\n", event.Name, event.Name)
            case "shortname":
                fmt.Printf("\n    The %s's shortname is %s.\n\n", event.Name, event.ShortName)
            case "website":
                fmt.Printf("\n    The %s's website can be found at %s\n\n", event.Name, event.Website)
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
                fmt.Printf("\n    Invalid datum \"%s\".\n\n", *eventDatum)
            }
        }
    }
    if matchCommand.Parsed() {
        playoffs := (*matchLevel)[len(*matchLevel)-1] == 'f'
        mk := ""
        if *matchKey != "" {
            mk = *matchKey
        } else {
            if playoffs {
                mk = fmt.Sprintf("%d%s_%s%dm%d", *matchYear, *matchEvent, *matchLevel, *matchNumber, *matchRound)
            } else {
                mk = fmt.Sprintf("%d%s_%s%d", *matchYear, *matchEvent, *matchLevel, *matchNumber)
            }
        }

        match, _ := tba.GetMatch(mk)

        matchLevelNames := map[string]string {
            "qm": "Qualifier",
            "qf": "Quarterfinal",
            "sf": "Semifinal",
            "f": "Final",
        }

        if match.Key == "" {
            fmt.Printf("\n    Invalid event key \"%s\".\n\n", mk)
            os.Exit(1)
        }

        // Often there will be a time given but no TimeString. This should correct for that.
        rawTime := time.Unix(int64(match.Time), 0)
        matchDate := fmt.Sprintf("%d/%d/%d", rawTime.Day(), rawTime.Month(), rawTime.Year())
        min := ""
        if rawTime.Minute() < 10 {
            min = fmt.Sprintf("0%d", rawTime.Minute())
        } else {
            min = fmt.Sprintf("%d", rawTime.Minute())
        }
        match.TimeString = fmt.Sprintf("%d:%s", rawTime.Hour(), min)

        if *matchDatum == "" {
            fmt.Printf("\n    ")
            if playoffs {
                c.Printf("%s %s #%d, Round %d (%s):\n", strings.ToUpper(match.EventKey), matchLevelNames[match.CompLevel], match.MatchNumber, match.SetNumber, match.Key)
            } else {
                c.Printf("%s %s #%d (%s):\n", strings.ToUpper(match.EventKey), matchLevelNames[match.CompLevel], match.MatchNumber, match.Key)
            }
            // TODO: Integrate event name into header
            // TODO: Show score breakdown
            g.Printf("\tDate: ")
            fmt.Println(matchDate)
            g.Printf("\tTime: ")
            fmt.Println(match.TimeString)
            g.Println("\tTeams: ")
            fmt.Println()
            for i := 0; i < 3; i++ {
                rTeam := match.Alliances.Red.Teams[i]
                bTeam := match.Alliances.Blue.Teams[i]
                b.Printf("\t    %s    ", rTeam[3:len(rTeam)])
                r.Printf("\t%s\n", bTeam[3:len(bTeam)])
            }
            fmt.Println()
        } else {
            switch strings.ToLower(*matchDatum) {
            case "key":
                fmt.Printf("\n    The match's key is %s.\n\n", match.Key)
            case "matchnumber":
                fmt.Printf("\n    Match %s's number is %d.\n\n", match.Key, match.MatchNumber)
            case "video":
                if (match.Videos[0].Type == "youtube") {
                    fmt.Printf("\n    A video of match %s can be found at https://youtube.com/watch?v=%s.\n\n", match.Key, match.Videos[0].Key)
                }
            case "time":
                fmt.Printf("\n    Match %s took place at %s on %s (Timestamp %d)\n\n", match.Key, match.TimeString, matchDate, match.Time)
            case "date":
                fmt.Printf("\t    Match %s took place on %s.\n", match.Key, matchDate)
            case "round":
                fmt.Printf("\n    Match %s is set number/round #%d.\n\n", match.Key, match.SetNumber)
        	case "teams":
                fmt.Printf("\n    In match %s, the red alliance was composed of %s, %s, and %s. The blue alliance was composed of %s, %s, and %s.", match.Alliances.Red.Teams[0], match.Alliances.Red.Teams[1], match.Alliances.Red.Teams[2], match.Alliances.Blue.Teams[0], match.Alliances.Blue.Teams[1], match.Alliances.Blue.Teams[2])
            case "event":
                fmt.Printf("\n    Match %s took place at event %s.\n\n", match.Key, match.EventKey)
            default:
                fmt.Printf("\n    Invalid datum \"%s\".\n\n", *matchDatum)
            }
        }
    }
}
