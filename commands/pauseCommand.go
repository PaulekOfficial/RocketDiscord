package commands

import (
	"github.com/andersfylling/disgord"
)

func init() {
	command := NewCommand("pause", false, false, onPauseCommand)
	command.SetArgumentsRequirements(true, -1, -1)
	command.SetHelpMessage("Polecenia pozwalajace na zatrzymanie muzyki na kanale glosowym. Użycie: /pause")
	command.Register()
}

func onPauseCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	voiceState := cache.GetVoiceState(guild.ID)
	if voiceState == nil {
		_, err := session.SendMsg(event.Message.ChannelID, ":upside_down: Nie moge zatrzymac muzyki, poniewaz nie jest ona uruchomiona")
		return err
	}

	voiceState.Paused = !voiceState.Paused

	if voiceState.Paused {
		_, err := session.SendMsg(event.Message.ChannelID, ":raccoon: Zatrzymano granie muzyki")
		return err
	} else {
		_, err := session.SendMsg(event.Message.ChannelID, ":notes: Wznowiono granie muzyki")
		return err
	}

	return
}
