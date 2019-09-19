package commands_test

import (
	"testing"

	"github.com/hkbarton/arango-cli/commands"
)

const rightMark = "\u2713"
const wrongMark = "\u2717"

func TestParse(t *testing.T) {
	cases := []struct {
		input         []string
		expectComamnd commands.Command
	}{
		{
			input:         []string{"list", "db", "-a"},
			expectComamnd: commands.Command{"list", []string{"db"}, map[string]string{"a": ""}},
		},
		{
			input:         []string{"info"},
			expectComamnd: commands.Command{Action: "info"},
		},
		{
			input:         []string{"use", "db"},
			expectComamnd: commands.Command{Action: "use", Args: []string{"db"}},
		},
	}
	t.Log("Given the need to test command parse")
	{
		for _, c := range cases {
			t.Logf("\tWhen pass a slice of string %v", c.input)
			{
				cmd, err := commands.Parse(c.input)
				if err != nil {
					t.Fatal("\t\tShould be able to generate command.", wrongMark, err)
				}
				t.Log("\t\tShould be able to generate command.", rightMark)
				if cmd.String() == c.expectComamnd.String() {
					t.Log("\t\tShould be able to generate command with right value.", rightMark)
				} else {
					t.Errorf("\t\tShould be able to generate command with right value. %s Expected: %s, got %s.",
						wrongMark, c.expectComamnd.String(), cmd.String())
				}
			}
		}
	}
}
