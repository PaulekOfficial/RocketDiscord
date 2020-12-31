package commands

import (
	"RocketDiscord/cache"
	"github.com/andersfylling/disgord"
)

func init() {
	command := newCommand("stop", false, false, onStopCommand)
	command.setArgumentsRequirements(true, -1, -1)
	command.setHelpMessage("Polecenia pozwalajace na wylaczenie muzyki na kanale glosowym. Użycie: /stop")
	command.register()
}

func onStopCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) (err error) {
	voiceState := cache.GetVoiceState(guild.ID)
	if voiceState == nil {
		_, err := session.SendMsg(event.Message.ChannelID, ":upside_down: Nie moge wylaczyc muzyki, poniewaz nie jest ona uruchomiona")
		return err
	}
	voiceState.Running = false
	voiceState.Tracks = make(map[int]*cache.MusicBotTrack)
	_, err = session.SendMsg(event.Message.ChannelID, ":raccoon: Zakończono odtwarzanie")
	return
}
