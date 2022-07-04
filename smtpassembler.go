package main

/* SmtpAssembler
   (C) 2021-2022 Péter Deák (hyper80@gmail.com)
   License: GPLv2
*/

import (
	"encoding/base64"
)

type SmtpAssembler struct {
	From           string
	To             string
	Subject        string
	BoundaryString string
	rawdata        string
}

func InitMessage() SmtpAssembler {
	mb := SmtpAssembler{}
	mb.BoundaryString = "======000smtpassemblerboundary000======"
	mb.rawdata = ""
	mb.From = ""
	mb.To = ""
	mb.Subject = "no subject"
	return mb
}

func (mb *SmtpAssembler) AddHtmlFull(html string) {
	mb.rawdata += mb.getBoundary()
	mb.rawdata +=
		"MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"Content-Transfer-Encoding: 8bit\r\n\r\n" + html + "\r\n"
}

func (mb *SmtpAssembler) AddHtmlBody(html string) {
	mb.AddHtmlFull("<html>" +
		"<head><title>" + mb.Subject + "</title></head>" +
		"<body>" + html + "</body>" +
		"</html>")
}

func (mb *SmtpAssembler) AddPlainText(txt string) {
	mb.rawdata += mb.getBoundary()
	mb.rawdata +=
		"MIME-version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"Content-Transfer-Encoding: 7bit\r\n\r\n" + txt + "\r\n"
}

func (mb *SmtpAssembler) AddAttachmentRaw(filename string, content []byte) {
	mb.AddAttachmentBase64(filename, base64.StdEncoding.EncodeToString(content))
}

func (mb *SmtpAssembler) AddAttachmentBase64(filename string, content string) {
	mb.rawdata += mb.getBoundary()
	mb.rawdata +=
		"MIME-Version: 1.0\r\n" +
			"Content-Type: application/octet-stream\r\n" +
			"Content-Transfer-Encoding: base64\r\n" +
			"Content-Disposition: attachment; filename=\"" + filename + "\"\r\n\r\n" +
			content + "\r\n\r\n"
}

func (mb *SmtpAssembler) AddAttachmentsBase64(attachments map[string]string) {
	for filename, contentb64 := range attachments {
		mb.AddAttachmentBase64(filename, contentb64)
	}
}

func (mb *SmtpAssembler) getBoundary() string {
	return "--" + mb.BoundaryString + "\r\n"
}

func (mb SmtpAssembler) GetRawMessage() []byte {
	return []byte("MIME-version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=\"" + mb.BoundaryString + "\"\r\n" +
		"From: " + mb.From + "\r\n" +
		"To: " + mb.To + "\r\n" +
		"Subject: " + mb.Subject + "\r\n\r\n" +
		"This is a message with multiple parts in MIME format.\n" +
		mb.rawdata +
		"--" + mb.BoundaryString + "--\r\n")
}
