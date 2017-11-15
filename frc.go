package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"
    "time"
    "github.com/fatih/color"

    tbago "github.com/ErikBoesen/tba-go"
)

const VERSION = "0.1.0"

// Used for printing in color
var (
    c = color.New(color.FgCyan, color.Underline)
    b = color.New(color.FgBlue)
    r = color.New(color.FgRed)
    g = color.New(color.FgGreen)
)


func main() {
    // First-level subcommands.
    teamCommand := flag.NewFlagSet("team", flag.ExitOnError)
    eventCommand := flag.NewFlagSet("event", flag.ExitOnError)
    matchCommand := flag.NewFlagSet("match", flag.ExitOnError)
    eventMatchesCommand := flag.NewFlagSet("eventmatches", flag.ExitOnError)

    // Team subcommand flags.
    teamNumber := teamCommand.Int("n", 0, "Team number. (Required)")

    // Event subcommand flags.
    eventKey := eventCommand.String("k", "", "ID of event you want data on. (Required)")

    // Match subcommand flags.
    matchKey := matchCommand.String("k", "", "Match key.")
    matchYear := matchCommand.Int("y", time.Now().Year(), "Year in which match took place.")
    matchEvent := matchCommand.String("e", "", "Event at which match occurred.")
    matchLevel := matchCommand.String("l", "", "Event level (qm, qf, sf, or f).")
    matchNumber := matchCommand.Int("n", 0, "Match number.")
    matchRound := matchCommand.Int("r", 0, "Match round (only in playoffs).")

    // Event matches subcommand flags.
    eventMatchesKey := eventMatchesCommand.String("k", "", "Key of event whose matches you desire. (Required)")
    eventMatchesTeam := eventMatchesCommand.Int("t", 0, "Number of team whose matches you want to show.")

    // Verify a subcommand has been provided.
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "Error: subcommand is required.")
        os.Exit(1)
    }

    // Initialize TBA parser
    tba, _ := tbago.Init("erikboesen", "frcli", VERSION)

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
    case "eventmatches":
        eventMatchesCommand.Parse(os.Args[2:])
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }

    if teamCommand.Parsed() {
        // Team key
        tk := *teamNumber

        // Fetch team data
        team, err := tba.GetTeam(tk)

        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid team key '%d'.\n", tk)
            os.Exit(1)
        }

        PrintTeam(team)
    } else if eventCommand.Parsed() {
        ek := *eventKey
        if _, err := strconv.Atoi(string((*eventKey)[0])); err != nil {
            ek = fmt.Sprintf("%d%s", time.Now().Year(), *eventKey)
        }

        event, err := tba.GetEvent(ek)

        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid event key '%s'.\n", ek)
            os.Exit(1)
        }

        event.StartDate = strings.Replace(event.StartDate, "-", "/", -1)
        event.EndDate = strings.Replace(event.EndDate, "-", "/", -1)

        event.VenueAddress = strings.Replace(event.VenueAddress, "\n", ", ", -1)

        PrintEvent(event)
    } else if matchCommand.Parsed() {
        mk := ""
        if *matchKey != "" {
            mk = *matchKey
        } else {
            if (*matchLevel)[len(*matchLevel)-1] == 'f' {
                mk = fmt.Sprintf("%d%s_%s%dm%d", *matchYear, *matchEvent, *matchLevel, *matchNumber, *matchRound)
            } else {
                mk = fmt.Sprintf("%d%s_%s%d", *matchYear, *matchEvent, *matchLevel, *matchNumber)
            }
        }

        match, err := tba.GetMatch(mk)

        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid event key '%s'.\n", mk)
            os.Exit(1)
        }

        PrintMatch(match)
    } else if eventMatchesCommand.Parsed() {
        ek := *eventMatchesKey
        if _, err := strconv.Atoi(string((*eventMatchesKey)[0])); err != nil {
            ek = fmt.Sprintf("%d%s", time.Now().Year(), *eventKey)
        }

        if *eventMatchesTeam == 0 {
            matches, err := tba.GetEventMatches(ek)
        } else {
            matches, err := tba.GetTeamEventMatches(*eventMatchesTeam, ek)
        }

        if len(matches) == 0 {
            fmt.Fprintf(os.Stderr, "No matches found for event '%s'.\n", ek)
            os.Exit(1)
        }

        for _, match := range matches {
            PrintMatch(match)
        }
    }
}

func PrintTeam(team tbago.Team) {
    fmt.Printf("\n    ")
    c.Printf("Team %d:\n", team.TeamNumber)
    g.Print("\tNickname:     ")
    fmt.Println(team.Nickname)
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
    fmt.Println(team.Motto)
    fmt.Println()
}

func PrintEvent(event tbago.Event) {
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
}

func PrintMatch(match tbago.Match) {
    levels := map[string]string {
        "qm": "Qualifier",
        "qf": "Quarterfinal",
        "sf": "Semifinal",
        "f": "Final",
    }
    fmt.Printf("\n    ")
    if match.CompLevel == "qm" {
        c.Printf("%s %s #%d (%s):\n", strings.ToUpper(match.EventKey), levels[match.CompLevel], match.MatchNumber, match.Key)
    } else {
        c.Printf("%s %s #%d, Round %d (%s):\n", strings.ToUpper(match.EventKey), levels[match.CompLevel], match.MatchNumber, match.SetNumber, match.Key)
    }
    if match.Time > 0 {
        g.Printf("\tDate/Time: ")
        fmt.Println(time.Unix(match.Time, 0).Format("06/01/02 at 15:01"))
    }
    g.Println("\tAlliances:\n")
    if match.Alliances.Red.Score > match.Alliances.Blue.Score {
        r.Printf("\t 🏆  ")
    } else {
        r.Printf("\t    ")
    }
    for index, team := range match.Alliances.Red.Teams {
        r.Printf("%s", team[3:len(team)])
        if index < 2 {
            r.Print(" | ")
        }
    }
    r.Printf(" => %d points\n", match.Alliances.Red.Score)
    if match.Alliances.Red.Score < match.Alliances.Blue.Score {
        b.Printf("\t 🏆  ")
    } else {
        b.Printf("\t    ")
    }
    for index, team := range match.Alliances.Blue.Teams {
        b.Printf("%s", team[3:len(team)])
        if index < 2 {
            r.Print(" | ")
        }
    }
    b.Printf(" => %d points\n\n", match.Alliances.Blue.Score)
}
