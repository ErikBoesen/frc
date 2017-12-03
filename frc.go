package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/fatih/color"
	"github.com/willf/pad"
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
	// Initialize TBA parser
	tba, err := tbago.New(KEY)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "frc"
	app.Usage = "handle FRC-related tasks in the command line."
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name: "team",
			Aliases: []string{"t"},
			Usage: "Get data on a team",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					log.Fatal("Usage: " + c.Command.Usage)
				}
				num, err := strconv.Atoi(c.Args()[0])
				if err != nil {
					log.Fatal("Invalid team number.\n")
					os.Exit(1)
				}
				team, err := tba.Team(num).Get()
				if err != nil {
					log.Fatal("TBA request unsuccessful.\n")
					os.Exit(1)
				}
				DisplayTeam(team)
				return nil
			},
		},
		{
			Name: "event",
			Aliases: []string{"e"},
			Usage: "Get data on an event",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					log.Fatal("Usage: " + c.Command.Usage)
					os.Exit(1)
				}

				event, err := tba.Event(c.Args()[0]).Get()
				if err != nil {
					log.Fatal("Invalid event key or failed TBA request.\n")
					os.Exit(1)
				}

				DisplayEvent(event)
				return nil
			},
		},
		{
			Name: "match",
			Aliases: []string{"m"},
			Usage: "Get data on a match",
			Action: func(c *cli.Context) error {
				match, err := tba.Match(c.Args()[0]).Get()
				if err != nil {
					log.Fatal("Invalid match key or failed TBA request.\n")
					os.Exit(1)
				}

				DisplayMatch(match)
				return nil
			},
		},
		{
			Name: "eventmatches",
			Aliases: []string{"em"},
			Usage: "Get data on the matches at an event",
			Action: func(c *cli.Context) error {
				var matches []tbago.Match
				// TODO: Don't discard error
				if len(c.Args()) > 1 {
					matches, _ = tba.Event(c.Args()[0]).Matches().Get()
				} else {
					team, err := strconv.Atoi(c.Args()[1])
					if err != nil {
						log.Fatal("Invalid team key.")
						os.Exit(1)
					}
					matches, _ = tba.Team(team).Event(c.Args()[0]).Matches().Get()
				}

				if len(matches) == 0 {
					log.Fatal("No matches found.\n")
					os.Exit(1)
				}

				for _, match := range matches {
					DisplayMatch(match)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
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
	data := []interface{}{fmt.Sprintf("%s - %s", event.StartDate, event.EndDate), event.Timezone, event.Website, event.LocationName, strings.Replace(event.Address, "\n", ", ", -1), event.District.DisplayName, fmt.Sprintf("%s (ID %d)", event.EventTypeString, event.EventType)}
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
	data := []interface{}{time.Unix(match.Time, 0).Format("06/01/02 at 15:01"), ""}
	ListInfo(header, titles, data)

	if match.Alliances.Red.Score > match.Alliances.Blue.Score {
		r.Printf("\t 🏆  ")
	} else {
		r.Printf("\t    ")
	}
	for index, team := range match.Alliances.Red.TeamKeys {
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
	for index, team := range match.Alliances.Blue.TeamKeys {
		b.Printf("%s", team[3:len(team)])
		if index < 2 {
			b.Print(" | ")
		}
	}
	b.Printf(" => %d points\n\n", match.Alliances.Blue.Score)
}
