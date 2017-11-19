package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/willf/pad"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/frc1418/tbago"
)

const KEY = "EzMD6D489Qttrf80Efz0rF9j3zRVz0pWuE0jfc4RlrUNA1yHDoaow8EN4THKIiJt"

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
		log.Fatal("Subcommand is required.")
	}

	// Initialize TBA parser
	tba, err := tbago.New(KEY)
	if err != nil {
		log.Fatal(err)
	}

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
		tk := *teamNumber

		team, err := tba.Team(tk).Get()
		if err != nil {
			log.Fatal(err)
		}

		DisplayTeam(team)
	} else if eventCommand.Parsed() {
		ek := *eventKey
		if _, err := strconv.Atoi(string((*eventKey)[0])); err != nil {
			ek = fmt.Sprintf("%d%s", time.Now().Year(), *eventKey)
		}

		event, err := tba.Event(ek).Get()
		if err != nil {
			log.Fatal(err)
		}

		event.Address = strings.Replace(event.Address, "\n", ", ", -1)

		DisplayEvent(event)
	} else if matchCommand.Parsed() {
		mk := ""
		if *matchKey != "" {
			mk = *matchKey
		} else {
			if *matchLevel == "qm" {
				mk = fmt.Sprintf("%d%s_%s%d", *matchYear, *matchEvent, *matchLevel, *matchNumber)
			} else {
				mk = fmt.Sprintf("%d%s_%s%dm%d", *matchYear, *matchEvent, *matchLevel, *matchNumber, *matchRound)
			}
		}

		match, err := tba.Match(mk).Get()
		if err != nil {
			log.Fatal(err)
		}

		DisplayMatch(match)
	} else if eventMatchesCommand.Parsed() {
		ek := *eventMatchesKey
		if _, err := strconv.Atoi(string((*eventMatchesKey)[0])); err != nil {
			ek = fmt.Sprintf("%d%s", time.Now().Year(), *eventKey)
		}

		var matches []tbago.Match
		// TODO: Don't discard error
		if *eventMatchesTeam == 0 {
			matches, _ = tba.Event(ek).Matches().Get()
		} else {
			// TODO: This functionality not yet supported by tbago
			//matches, _ = tba.GetTeamEventMatches(*eventMatchesTeam, ek)
		}

		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "No matches found for event '%s'.\n", ek)
			os.Exit(1)
		}

		for _, match := range matches {
			DisplayMatch(match)
		}
	}
}

func ListInfo(header string, titles []string, data []interface{}) {
	fmt.Printf("\n    ")
	c.Printf("%s:\n", header)
	for i := range titles {
		g.Printf("\t%s ", pad.Right(titles[i]+":", 10, " "))
		fmt.Println(data[i])
	}
	fmt.Println()
}

func DisplayTeam(team tbago.Team) {
	header := fmt.Sprintf("Team %d", team.TeamNumber)
	titles := []string{"Nickname", "Website", "Rookie", "Country", "Motto"}
	data := []interface{}{team.Nickname, team.Website, team.RookieYear, team.Country, team.Motto}
	ListInfo(header, titles, data)
}

func DisplayEvent(event tbago.Event) {
	header := fmt.Sprintf("%d %s (%s)", event.Year, event.Name, event.Key)
	titles := []string{"Date", "Timezone", "Website", "Location", "Address", "District", "Type"}
	data := []interface{}{fmt.Sprintf("%s - %s", event.StartDate, event.EndDate), event.Timezone, event.Website, event.LocationName, event.Address, event.District.DisplayName, fmt.Sprintf("%s (ID %d)", event.EventTypeString, event.EventType)}
	ListInfo(header, titles, data)
}

func DisplayMatch(match tbago.Match) {
	header := fmt.Sprintf("%s %s #%d", strings.ToUpper(match.EventKey), strings.ToUpper(match.CompLevel), match.MatchNumber)
	if match.CompLevel == "qm" {
		header += fmt.Sprintf(" (%s)", match.Key)
	} else {
		header += fmt.Sprintf(", Round %d (%s)", match.SetNumber, match.Key)
	}
	titles := []string{"Date/Time", "Alliances"}
	data := []interface{}{time.Unix(int64(match.Time), 0).Format("06/01/02 at 15:01"), ""}
	ListInfo(header, titles, data)

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
			b.Print(" | ")
		}
	}
	b.Printf(" => %d points\n\n", match.Alliances.Blue.Score)
}
