package commands

import (
	"github.com/andersfylling/disgord"
	"math/rand"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"

var lastKey string

func init() {
	command := NewCommand("shutdown", false, false, shutdownCommand)
	command.SetArgumentsRequirements(true, 0, 1)
	command.Register()
}

func shutdownCommand(session disgord.Session, event *disgord.MessageCreate, guild *disgord.Guild, args []string) error {
	//if event.Author.ID != "419203785190014987" {
	//	_, err := session.ChannelMessageSend(event.ChannelID, ":rotating_light: Brak wystarczających uprawnień do użycia tego polecenia! Wymagany poziom dostępu :infinity:")
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//}
	//
	//if len(args) == 0 {
	//	_, err := session.ChannelMessageSend(event.ChannelID, ":sob: Prosze nie zabijaj mnie!")
	//	if err != nil {
	//		return err
	//	}
	//
	//	lastKey = randomString(rand.Intn(50))
	//	channel, err := session.UserChannelCreate("419203785190014987")
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = session.ChannelMessageSend(channel.ID, fmt.Sprintf(":comet: Hic tu qua concessum sit, ut disable vestri codice mihi: %s", lastKey))
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if len(args) == 1 && lastKey != "" {
	//	if args[0] != lastKey {
	//		_, err := session.ChannelMessageSend(event.ChannelID, ":rotating_light: In codice provisum non valet, valet in codice intra placet!")
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	_, err := session.ChannelMessageSend(event.ChannelID, ":clock1: Confirmavit codice, ut shutdown satus in ordine ...")
	//	if err != nil {
	//		return err
	//	}
	//	err = session.UpdateStatus(0, "Shutting Down...")
	//	if err != nil {
	//		return err
	//	}
	//	time.Sleep(time.Minute)
	//	err = session.Close()
	//	if err != nil {
	//		return err
	//	}
	//	os.Exit(0)
	//}
	return nil
}

func randomString(length int) string {
	randomCharacters := make([]byte, length)

	for i := range randomCharacters {
		randomCharacters[i] = CHARSET[rand.Intn(length)]
	}

	return string(randomCharacters)
}
