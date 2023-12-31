package textmanipulation

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mozoarella/tibby/internal/types"
	"github.com/mozoarella/tibby/internal/utils"
)

type Randomizer = types.Randomizer

var (
	furryFlavourText Randomizer
)

func Uwuify(i *discordgo.InteractionCreate) string {

	msg := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	input := msg.Content
	msgUrl := utils.ReturnInteractionMessageUrl(i)

	if len(msg.Embeds) > 0 {
		return fmt.Sprintf("I can't mock embeds, sorry.\n\n[Go to the original message](%v)", msgUrl)
	}

	if len(msg.Content) < 1 {
		return fmt.Sprintf("I can't mock messages without text, sorry.\n\n[Go to the original message](%v)", msgUrl)
	}

	/*wordReplacer := strings.NewReplacer(
	"friend", "fwiendo",
	"ove", "uv",
	"hugs", "huggies",
	"hug", "huggy",
	"kisses", "smoochies",
	"bird", "birb",
	"chicken", "chinkem")*/
	letterReplacer := strings.NewReplacer(
		"r", "w",
		"R", "W",
		"l", "w",
		"L", "W")

	replaced := letterReplacer.Replace(input)
	flavoured := addFlavour(replaced)

	var output_format string = `%v
	
	[Go to the original message](%v)
	`

	built_message := fmt.Sprintf(output_format, flavoured, msgUrl)
	return built_message
}

func addFlavour(input string) string {
	furryFlavourText.Append(" UwU ", " OwO ", " ^_^ ", " :3 ", ` \*nuzzles u\* `, "~ ", " nya~ ", " (・`ω´・) ")

	textSlice := strings.Split(input, "")

	for k, v := range textSlice {
		if v == "!" {
			textSlice[k] = furryFlavourText.Random()
		}
	}
	return strings.Join(textSlice, "")
}
