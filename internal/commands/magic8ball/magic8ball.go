package magic8ball

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mozoarella/wombot/internal/types"
)

type Randomizer = types.Randomizer

var (
	responses Randomizer
)

func init() {
	responses.Fill("data/8ballresponses.txt", true)
}

func ShakeTheBall(i *discordgo.InteractionCreate) string {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var shaker string
	var question string
	ballResponse := responses.Random()

	if i.User != nil {
		shaker = i.User.Username
	} else {
		shaker = i.Member.Nick
	}

	if val, ok := optionMap["question"]; ok {
		question = val.StringValue()
	}

	var fullResponse string

	var responseWithoutQ string = `*%v shakes the Magic 8-ball*
	
The Magic 8-ball says: 
**%v**`

	var responseWithQ string = `%v asks the Magic 8-ball: "%v"

*They give the ball a good shake*

The Magic 8-ball says: 
**%v**`

	if len(question) == 0 {
		fullResponse = fmt.Sprintf(responseWithoutQ, shaker, ballResponse)
	} else {
		fullResponse = fmt.Sprintf(responseWithQ, shaker, question, ballResponse)
	}

	return fullResponse
}