package commands

import (
	"github.com/PaulekOfficial/RocketDiscord/cache"
	"github.com/andersfylling/disgord"
)

func init() {
	command := NewCommand("loop", false, false, onLoopCommand)
	command.SetArgumentsRequirements(true, -1, -1)
	command.SetHelpMessage("Polecenia pozwalajace na zapętlenie konretnego utworu muzyki na kanale glosowym. Użycie: /loop")
	command.Register()
}

func onLoopCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	voiceState := cache.GetVoiceState(guild.ID)
	if voiceState == nil {
		_, err := session.SendMsg(event.Message.ChannelID, ":upside_down: Nie moge zapętlić muzyki, poniewaz nie jest ona uruchomiona")
		return err
	}

	voiceState.LoopTrack = !voiceState.LoopTrack

	if voiceState.LoopTrack {
		_, err := session.SendMsg(event.Message.ChannelID, ":raccoon: Zapętlono utwór")
		return err
	} else {
		_, err := session.SendMsg(event.Message.ChannelID, ":notes: Wyłączono zapętlenie utworu")
		return err
	}

	return
}
