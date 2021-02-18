package main

/*  Rest-Go-Mail - HTTPS-REST capable (html) e-mail sender agent
    (C) 2021 Péter Deák (hyper80@gmail.com)
    License: GPLv2
*/

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type configItemType struct {
	sval     string
	fval     float64
	bval     bool
	vtype    string
	required bool
	def      string
}

type mailDataType struct {
	from     string
	to       string
	subject  string
	bodyhtml string
}

var config map[string]*configItemType
var knownCertificates map[string]string
var senderChannel chan mailDataType

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05 RestGoMail: ") + string(bytes))
}

func handleFallbackHTTPReq(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's work"))
}

func handleHTTPMailReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.URL.Path != "/sendmail" {
		handleFallbackHTTPReq(w, r)
		return
	}

	var status bool = true
	if body, rdBodyErr := ioutil.ReadAll(r.Body); rdBodyErr != nil {
		log.Printf("Error (%s): %s\n", r.RemoteAddr, rdBodyErr.Error())
	} else {
		status = processRequest(&body, r.RemoteAddr)
	}

	message := "{status: \"Received\"}"
	if status {
		message = "{status: \"Failed\"}"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "utf-8")
	w.Write([]byte(message))
}

func checkClientCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if !getConfigBool("allowOnlyKnownCertificates") {
		return nil
	}

	if len(rawCerts) < 1 {
		return errors.New("client certificate not found")
	}

	for i := 0; i < len(rawCerts); i++ {
		b64cert := base64.StdEncoding.EncodeToString(rawCerts[i])
		if len(b64cert) != 0 {
			for name, certval := range knownCertificates {
				if certval == b64cert {
					log.Println("Accepted client certification named:", name)
					return nil
				}
			}
		}
	}
	return errors.New("not matching client certificate")
}

func initConfig() {
	config = make(map[string]*configItemType)
	config["smtpHost"] = &configItemType{vtype: "string", required: true}
	config["smtpPort"] = &configItemType{vtype: "string", required: true}
	config["smtpAuthRequired"] = &configItemType{vtype: "bool", required: true, def: "false"}
	config["smtpAuthIdentity"] = &configItemType{vtype: "string", required: false, def: ""}
	config["smtpAuthPassword"] = &configItemType{vtype: "string", required: false, def: ""}
	config["smtpAllowedFromAddressOnly"] = &configItemType{vtype: "string", required: false, def: ""}
	config["httpsListenPort"] = &configItemType{vtype: "string", required: true}
	config["tlsKeyFile"] = &configItemType{vtype: "string", required: true}
	config["tlsCertFile"] = &configItemType{vtype: "string", required: true}
	config["allowOnlyKnownCertificates"] = &configItemType{vtype: "bool", required: false, def: "false"}
	config["waitSecondsAfterSmtpReq"] = &configItemType{vtype: "float", required: false, def: "12"}
	config["debugMode"] = &configItemType{vtype: "bool", required: false, def: "false"}
}

func readConfig(jsonfile string) bool {
	confJSONData, confJSONFileErr := ioutil.ReadFile(jsonfile)
	if confJSONFileErr != nil {
		log.Printf("Error, cannot read JSON file: %s\n", confJSONFileErr.Error())
		return true
	}

	var confJSONParsed map[string]interface{}
	confJSONError := json.Unmarshal(confJSONData, &confJSONParsed)
	if confJSONError != nil {
		log.Printf("Error, configuration file has not valid JSON: %s\n", confJSONError.Error())
		return true
	}

	for confName, configItemVal := range config {
		switch configItemVal.vtype {
		case "string":
			sv, svt := getStringByPath(confJSONParsed, "restgomail/"+confName)
			if configItemVal.required && svt == "none" {
				log.Printf("Error, incomplete config JSON, \"%s\" missing", confName)
				return true
			}
			if svt == "none" {
				config[confName].sval = config[confName].def
			} else {
				config[confName].sval = sv
			}
		case "float":
			fv, fvt := getFloat64ByPath(confJSONParsed, "restgomail/"+confName)
			if configItemVal.required && fvt == "none" {
				log.Printf("Error, incomplete config JSON, \"%s\" missing", confName)
				return true
			}
			if fvt == "none" {
				config[confName].fval, _ = strconv.ParseFloat(config[confName].def, 32)
			} else {
				config[confName].fval = fv
			}
		case "bool":
			bv, bvt := getBoolByPath(confJSONParsed, "restgomail/"+confName)
			if configItemVal.required && bvt == "none" {
				log.Printf("Error, incomplete config JSON, \"%s\" missing", confName)
				return true
			}
			if bvt == "none" {
				if config[confName].def == "true" {
					config[confName].bval = true
				} else {
					config[confName].bval = false
				}
			} else {
				config[confName].bval = bv
			}
		}
	}

	loadedCertsCounts := 0
	cm, cmt := getMapByPath(confJSONParsed, "restgomail/knownCertificates")
	if cmt == "map" {
		for name, value := range cm {
			if certstr, isStr := value.(string); isStr {
				certval := certstr
				if strings.HasPrefix(certstr, "@") {
					certbytes, rdCertFilErr := ioutil.ReadFile(certstr[1:])
					if rdCertFilErr != nil {
						log.Printf("Error, cannot read certification file \"%s\"", certstr[1:])
						return true
					}
					certval = string(certbytes)
					reBeg := regexp.MustCompile(`\-+BEGIN CERTIFICATE\-+`)
					certval = reBeg.ReplaceAllString(certval, "")
					reEnd := regexp.MustCompile(`\-+END CERTIFICATE\-+`)
					certval = reEnd.ReplaceAllString(certval, "")
					reWsp := regexp.MustCompile(`\s`)
					certval = reWsp.ReplaceAllString(certval, "")
				}
				knownCertificates[name] = certval
				loadedCertsCounts++
			}
		}
	}

	if getConfigBool("debugMode") {
		log.Printf("%d white list cert loaded.", loadedCertsCounts)
	}
	return false
}

func getConfigString(confname string) string {
	return config[confname].sval
}

func getConfigBool(confname string) bool {
	return config[confname].bval
}

func getConfigFloat64(confname string) float64 {
	return config[confname].fval
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	if len(os.Args) < 2 {
		log.Println("Error: You must pass a configuration JSON file name as parameter.")
		return
	}

	knownCertificates = make(map[string]string)
	initConfig()
	if readConfig(os.Args[1]) {
		return
	}

	senderChannel = make(chan mailDataType, 32)
	go senderAgent(senderChannel)

	log.Printf("Start listening on %s...\n", getConfigString("httpsListenPort"))

	server := &http.Server{
		Addr:    ":" + getConfigString("httpsListenPort"),
		Handler: nil,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:    true,
			ClientAuth:            tls.RequestClientCert,
			VerifyPeerCertificate: checkClientCert,
		},
	}

	http.HandleFunc("/", handleFallbackHTTPReq)
	http.HandleFunc("/sendmail", handleHTTPMailReq)
	if err := server.ListenAndServeTLS(getConfigString("tlsCertFile"),
		getConfigString("tlsKeyFile")); err != nil {
		panic(err)
	}
}

func processRequest(req *[]byte, remote string) bool {
	var parsed map[string]interface{}

	if getConfigBool("debugMode") {
		log.Println("******* BEGIN Request data ************")
		log.Println(string(*req))
		log.Println("******* END Request data **************")
	}
	error := json.Unmarshal(*req, &parsed)
	if error != nil {
		log.Printf("Error (%s) Not valid JSON: %s\n", remote, error.Error())
		return true
	}

	if getConfigBool("debugMode") {
		log.Println("JSON from client: ", remote)
		log.Println("------- BEGIN Parsed JSON -------------")
		log.Println(printParsedJSON(parsed))
		log.Println("------- END Parsed JSON ---------------")
	}

	from, fromType := getStringByPath(parsed, "sendmail/from")
	to, toType := getStringByPath(parsed, "sendmail/to")
	subject, subjectType := getStringByPath(parsed, "sendmail/subject")
	bodyhtml, bodyhtmlType := getStringByPath(parsed, "sendmail/bodyhtml")

	if fromType == "none" || toType == "none" ||
		subjectType == "none" || bodyhtmlType == "none" {
		log.Printf("Error (%s) incomplete request\n", remote)
		if getConfigBool("debugMode") {
			log.Println(" from...", fromType)
			log.Println(" to...", toType)
			log.Println(" subject...", subjectType)
			log.Println(" bodyhtml...", bodyhtmlType)
		}
		return true
	}

	log.Printf("Received send mail request from %s\n", remote)
	allowedSnAddr := getConfigString("smtpAllowedFromAddressOnly")
	if allowedSnAddr != "" && allowedSnAddr != from {
		log.Printf("Error (%s) not allowed sender address\n", remote)
		return true
	}
	senderChannel <- mailDataType{from, to, subject, bodyhtml}
	return false
}

func senderAgent(recvChannel <-chan mailDataType) {
	authRequired := getConfigBool("smtpAuthRequired")
	var mail mailDataType
	for {
		mail = <-recvChannel
		if authRequired {
			sendmailReqAuth(mail)
		} else {
			sendmailNoAuth(mail)
		}
		time.Sleep(time.Second * time.Duration(getConfigFloat64("waitSecondsAfterSmtpReq")))
	}
}

func smtpCreateRawBody(from, to, subject, body string) []byte {
	return []byte("" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"utf-8\";\r\n\r\n" +
		"<html>" +
		"<head><title>" + subject + "</title></head>" +
		"<body>" + body + "</body>" +
		"</html>" +
		"\r\n")
}

func sendmailNoAuth(mail mailDataType) {
	c, err := smtp.Dial(getConfigString("smtpHost") + ":" + getConfigString("smtpPort"))
	if err != nil {
		log.Println("Error(11) :", err.Error())
		return
	}
	defer c.Close()

	c.Mail(mail.from)
	c.Rcpt(mail.to)

	wc, err := c.Data()
	if err != nil {
		log.Println("Error(12) :", err.Error())
		return
	}
	defer wc.Close()

	message := smtpCreateRawBody(mail.from, mail.to, mail.subject, mail.bodyhtml)

	if _, err = wc.Write(message); err != nil {
		log.Println("Error(13) :", err.Error())
		return
	}
	log.Println("Message sent to ", mail.to)
}

func sendmailReqAuth(mail mailDataType) {
	toarray := []string{
		mail.to,
	}

	message := smtpCreateRawBody(mail.from, mail.to, mail.subject, mail.bodyhtml)
	auth := smtp.PlainAuth(getConfigString("smtpAuthIdentity"), mail.from, getConfigString("smtpAuthPassword"), getConfigString("smtpHost"))

	err := smtp.SendMail(getConfigString("smtpHost")+":"+getConfigString("smtpPort"), auth, mail.from, toarray, message)
	if err != nil {
		log.Println("Error(21) :", err.Error())
		return
	}
	log.Println("Message sent to ", mail.to)
}