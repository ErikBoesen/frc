package main

import (
    "flag"
    "fmt"
    "os"
    "strings"
    "strconv"

    "github.com/CarlColglazier/tba"
)

const (
    VERSION = "0.1.0"
)

func main() {
    // Subcommands
    teamCommand := flag.NewFlagSet("team", flag.ExitOnError)

    // Count subcommand flag pointers
    // Adding a new choice for --metric of 'substring' and a new --substring flag
    teamDataPointPtr := teamCommand.String("d", "", "Data point to fetch. (Required)")

    // List subcommand flag pointers.
    /*listTextPtr := listCommand.String("text", "", "Text to parse. (Required)")
    listMetricPtr := listCommand.String("metric", "chars", "Metric <chars|words|lines>. (Required)")
    listUniquePtr := listCommand.Bool("unique", false, "Measure unique values of a metric.")*/

    // Verify a subcommand has been provided.
    // os.Arg[0] = main command
    // os.Arg[1] = subcommand
    if len(os.Args) < 2 {
        fmt.Println("Error: subcommand is required.")
        os.Exit(1)
    }

    // _ recieves any error. We should probably be outputting errors, but we'll keep it simple for now
    tba, _ := tba.Init("erikboesen", "tbacli", VERSION)

    // Parse flags for appropriate FlagSet
    // FlagSet.Parse() requires a set of arguments to parse as input
    // os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
    switch os.Args[1] {
    case "team":
        teamCommand.Parse(os.Args[2:])
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Check which subcommand was Parsed using the FlagSet.Parsed() function; then handle each case accordingly.
    // FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user didn't provide any flags).
    if teamCommand.Parsed() {

        // TODO: Handle if the user forgets the team parameter.

        teamNumber, _ := strconv.Atoi(os.Args[2])

        team, _ := tba.GetTeam(teamNumber)

        fmt.Println(strings.ToLower(*teamDataPointPtr))

        // TODO: Find an alternative to this method.
        switch strings.ToLower(*teamDataPointPtr) {
        case "website":
            fmt.Println(team.Website)
        case "name":
            fmt.Println(team.Name)
    	case "locality":
            fmt.Println(team.Locality)
    	case "rookieyear":
            fmt.Println(team.RookieYear)
    	case "region":
            fmt.Println(team.Region)
        case "teamnumber":
            fmt.Println(team.TeamNumber)
    	case "location":
            fmt.Println(team.Location)
    	case "key":
            fmt.Println(team.Key)
    	case "countryname":
            fmt.Println(team.CountryName)
    	case "motto":
            fmt.Println(team.Motto)
    	case "nickname":
            fmt.Println(team.Nickname)
        default:
            fmt.Println("Test")
        }


        // Print
        fmt.Printf("textPtr: %s, metricPtr: %s\n", os.Args[2], *teamDataPointPtr)
    }
}
