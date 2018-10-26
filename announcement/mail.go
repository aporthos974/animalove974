package announcement

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"net/smtp"
	"reunion/configuration"
)

type MailRequest struct {
	Sender, Message, AnnouncementId string
}

type Mail struct {
	Sender, Recipient, Message string
}

func GetContactMessageHandler(writer http.ResponseWriter, request *http.Request) {
	var mail = Mail{Recipient: configuration.Conf.GetMail().Contact}
	parseBody(request, &mail)
	if mail.Sender == "" || mail.Message == "" {
		log.Panicf("Paramètre manquant, Sender : %s, Message : %s", mail.Sender, mail.Message)
	}

	sendMail(mail, "Bonjour,\n\nVoici un message envoyé de la part de : "+mail.Sender+"\n\n"+mail.Message)
}

func GetMailHandler(writer http.ResponseWriter, request *http.Request) {
	var mailRequest = MailRequest{}
	parseBody(request, &mailRequest)
	if mailRequest.AnnouncementId == "" || mailRequest.Message == "" {
		log.Panicf("Paramètre manquant, announcementId : %s, message : %s", mailRequest.AnnouncementId, mailRequest.Message)
	}
	mail := Mail{Message: mailRequest.Message, Sender: mailRequest.Sender}

	var announcement Announcement
	collection, session := configuration.GetAnnouncementCollection()
	defer session.Close()
	if err := collection.FindId(bson.ObjectIdHex(mailRequest.AnnouncementId)).One(&announcement); err == nil {
		mail.Recipient = announcement.Account.Email
		sendMail(mail, "Bonjour,\n\nVoici un message envoyé de la part de : "+mail.Sender+"\n\n"+mail.Message)
	} else {
		log.Panicf("Erreur lors de l'accès à la bdd : %s", err.Error())
	}
}

func parseBody(request *http.Request, mailRequest interface{}) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Panicf("Erreur dans le contenu de la requête : %s", err.Error())
	}
	err = json.Unmarshal(body, mailRequest)
	if err != nil {
		log.Panicf("Erreur dans le contenu de la requête : %s", err.Error())
	}
}

func sendMail(mailContent Mail, format string) {
	header := make(map[string]string)
	header["From"] = mailContent.Sender
	header["To"] = mailContent.Recipient
	header["Subject"] = "Site des animaux perdus - Message"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	var message string
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(format+"\n\nAnimaLove\n"+configuration.Conf.GetUrl()))

	auth := smtp.PlainAuth("", configuration.Conf.GetMail().Contact, "animalove311212", "smtp.animalove.re")

	err := smtp.SendMail("smtp.animalove.re:587", auth, configuration.Conf.GetMail().Contact, []string{mailContent.Recipient}, []byte(message))
	if err != nil {
		log.Printf("Erreur lors l'envoi de mail : %s", err.Error())
	}

}
