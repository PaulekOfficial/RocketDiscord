package commands

import (
	"RocketDiscord/cache"
	"github.com/andersfylling/disgord"
)

func init() {
	command := newCommand("loop", false, false, onLoopCommand)
	command.setArgumentsRequirements(true, -1, -1)
	command.setHelpMessage("Polecenia pozwalajace na zapętlenie konretnego utworu muzyki na kanale glosowym. Użycie: /loop")
	command.register()
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
