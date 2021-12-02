package commands

import (
	"github.com/andersfylling/disgord"
	"strings"
)

func init() {
	command := NewCommand("rule34", false, true, onRule34Request)
	command.SetArgumentsRequirements(true, 1, -1)
	command.SetHelpMessage("Polecenie pozwalające pokazać obrazek ze strony rule34. Ostrzeżenie treść jest przeznaczona dla osób pełnoletnich!  Użycie: !rule34 <tags...>")
	command.Register()
}

var censored []string = []string{"loli", "shota", "child", "young"}

func onRule34Request(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	for _, arg := range args {
		for _, censoredWord := range censored {
			if strings.Contains(strings.ToLower(arg), censoredWord) {
				_, err := session.SendMsg(event.Message.ChannelID, ":man_police_officer: <!> Podane tagi rule34 są zakazane <!>")
				return err
			}
		}
	}

	urls, err := utils.GetRule34UrlImageFromTags(args, 4)
	if err != nil {
		return err
	}

	if urls == nil {
		_, err = session.SendMsg(event.Message.ChannelID, ":man_facepalming: Nie znalazlem żadnego dopasowania dla tych tagów.")
		return err
	}

	for _, url := range urls {
		_, err = session.SendMsg(event.Message.ChannelID, url)
		if err != nil {
			return err
		}
	}
	return nil
}
