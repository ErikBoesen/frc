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

    tbago "github.com/ErikBoesen/tba-go"
)

const VERSION = "0.1.0"

// Used for printing in color
var c = color.New(color.FgCyan, color.Underline)
var b = color.New(color.FgBlue)
var r = color.New(color.FgRed)
var g = color.New(color.FgGreen)

// Nicer names for match levels.
var matchLevelNames = map[string]string {
    "qm": "Qualifier",
    "qf": "Quarterfinal",
    "sf": "Semifinal",
    "f": "Final",
}


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

        // Trim URL fragments from start of event website URL
        team.Website = urlRE.ReplaceAllString(team.Website, "")

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

        // Trim URL fragments from start of event website URL
        event.Website = urlRE.ReplaceAllString(event.Website, "")

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

        var matches []tbago.Match

        if *eventMatchesTeam == 0 {
            matches, _ = tba.GetEventMatches(ek)
        } else {
            matches, _ = tba.GetTeamEventMatches(*eventMatchesTeam, ek)
        }

        if len(matches) == 0 {
            fmt.Fprintf(os.Stderr, "No matches found for event '%s'.\n", ek)
            os.Exit(1)
        }

        for i := 0; i < len(matches); i++ {
            PrintMatch(matches[i])
        }
    }
}

func PrintTeam(team tbago.Team) {
    fmt.Printf("\n    ")
    c.Printf("Team %d:\n", team.TeamNumber)
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
    if len(team.Motto) > 0 {
        g.Print("\tMotto:        ")
        if string(team.Motto[0]) == "\"" {
            fmt.Println(team.Motto)
        } else {
            fmt.Printf("\"%s\"\n", team.Motto)
        }
    }
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
    // Often there will be a time given but no TimeString. This should correct for that.
    rawTime := time.Unix(int64(match.Time), 0)
    matchDate := fmt.Sprintf("%d/%d/%d", rawTime.Day(), rawTime.Month(), rawTime.Year())
    min := ""
    if rawTime.Minute() < 10 {
        min = fmt.Sprintf("0%d", rawTime.Minute())
    } else {
        min = fmt.Sprintf("%d", rawTime.Minute())
    }
    match.TimeString = fmt.Sprintf("%s at %d:%s", matchDate, rawTime.Hour(), min)

    fmt.Printf("\n    ")
    if (match.CompLevel)[len(match.CompLevel)-1] == 'f' {
        c.Printf("%s %s #%d, Round %d (%s):\n", strings.ToUpper(match.EventKey), matchLevelNames[match.CompLevel], match.MatchNumber, match.SetNumber, match.Key)
    } else {
        c.Printf("%s %s #%d (%s):\n", strings.ToUpper(match.EventKey), matchLevelNames[match.CompLevel], match.MatchNumber, match.Key)
    }
    // Once in a while a
    if match.Time > 0 {
        g.Printf("\tDate/Time: ")
        fmt.Println(match.TimeString)
    }
    g.Println("\tAlliances:\n")
    if match.Alliances.Red.Score > match.Alliances.Blue.Score {
        r.Printf("\t 🏆  ")
    } else {
        r.Printf("\t    ")
    }
    for i := 0; i < 3; i++ {
        rTeam := match.Alliances.Red.Teams[i]
        r.Printf("%s", rTeam[3:len(rTeam)])
        if i < 2 {
            r.Print(" | ")
        }
    }
    r.Printf(" => %d points\n", match.Alliances.Red.Score)
    if match.Alliances.Red.Score < match.Alliances.Blue.Score {
        b.Printf("\t 🏆  ")
    } else {
        b.Printf("\t    ")
    }
    for i := 0; i < 3; i++ {
        bTeam := match.Alliances.Blue.Teams[i]
        b.Printf("%s", bTeam[3:len(bTeam)])
        if i < 2 {
            b.Print(" | ")
        }
    }
    b.Printf(" => %d points\n\n", match.Alliances.Blue.Score)
}
