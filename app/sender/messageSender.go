package sender

import (
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nanu-c/textsecure"
	"github.com/nanu-c/textsecure-qml/app/helpers"
	"github.com/nanu-c/textsecure-qml/app/store"
)

func SendMessageHelper(to, message, file string) (error, *store.Message) {
	var err error
	if file != "" {
		file, err = store.CopyAttachment(file)
		log.Debugln("[axolotl] attachment: " + file)
		if err != nil {
			log.Errorln("Error Attachment:" + err.Error())
			return err, nil
		}
	}
	session := store.SessionsModel.Get(to)
	m := session.Add(message, "", file, "", true, store.ActiveSessionID)
	m.Source = to
	_, savedM := store.SaveMessage(m)
	go SendMessage(session, m)
	return nil, savedM
}
func SendMessage(s *store.Session, m *store.Message) {
	var att io.Reader
	var err error

	if m.Attachment != "" {
		att, err = os.Open(m.Attachment)
		if err != nil {
			return
		} else {
			log.Printf("SendMessage FileOpend")
		}
	}

	ts := SendMessageLoop(s.Tel, m.Message, s.IsGroup, att, m.Flags)
	// log.Debugln(ts, s)
	m.SentAt = ts
	s.Timestamp = m.SentAt
	m.IsSent = true
	//FIXME avoid rerendering the whole qml
	m.HTime = helpers.HumanizeTimestamp(m.SentAt)
	s.When = m.HTime
	store.UpdateMessageSent(m)
	store.UpdateSession(s)
}
func SendMessageLoop(to, message string, group bool, att io.Reader, flags int) uint64 {
	var err error
	var ts uint64
	var count int
	for {
		err = nil
		if flags == helpers.MsgFlagResetSession {
			ts, err = textsecure.EndSession(to, "TERMINATE")
		} else if flags == helpers.MsgFlagGroupLeave {
			err = textsecure.LeaveGroup(to)
		} else if flags == helpers.MsgFlagGroupUpdate {
			_, err = textsecure.UpdateGroup(to, store.Groups[to].Name, strings.Split(store.Groups[to].Members, ","))
		} else if att == nil {
			if group {
				ts, err = textsecure.SendGroupMessage(to, message)
				log.Debugln("[axolotl] send to group ")
			} else {
				ts, err = textsecure.SendMessage(to, message)
			}
		} else {
			if group {
				ts, err = textsecure.SendGroupAttachment(to, message, att)
			} else {
				log.Printf("SendMessageLoop sendAttachment")
				// buf := new(bytes.Buffer)
				// buf.ReadFrom(att)
				// s := buf.String()
				// log.Printf(s)

				ts, err = textsecure.SendAttachment(to, message, att)
			}
		}
		if err == nil {
			break
		}
		log.Println(err)
		//If sending failed, try again after a while
		time.Sleep(3 * time.Second)
		count++
		if count == 2 {
			// return nil, new Error("sending")
			break

		}
	}
	return ts
}
func SendUnsentMessages() {
	// for _, s := range store.SessionsModel.Sess {
	// 	for _, m := range s.Messages {
	// 		if m.Outgoing && !m.IsSent {
	// 			go SendMessage(s, m)
	// 		}
	// 	}
	// }
}