![RestGoMail logo](https://raw.githubusercontent.com/hyper-prog/restgomail/main/images/restgomail.png)

RestGoMail - HTTP-REST Mail gateway in Go
==========================================

RestGoMail is a small daemon/container which able to receive HTML e-mail forward requests 
as HTTP POST in a JSON data, queue the requests, and sends the specified
mails to a SMTP server according to the settings.

Docker images
-------------
Available a docker container with a compiled restgomail daemon:
 Docker hub:

- https://hub.docker.com/r/hyperprog/restgomail

 Downloadable (pullable) image name:
 
    hyperprog/restgomail

Check the docker-compose.yml file under the EXAMPLES directory to learn how to configure it.

Config file sample (JSON) 
-------------------------
You need to give a JSON file as a command line argument to specify
SMTP host and port, the authentications, allowed client certificates and so on...

    {
        "restgomail": {
            "httpsListenPort": "443",
            "smtpHost": "smtp.gmail.com",
            "smtpPort": "587",
            "smtpAuthRequired": true,
            "smtpAuthPassword": "gmailpasswordsample",
            "smtpAllowedFromAddressOnly": "sampleuser@gmail.com",
            "tlsKeyFile": "restgomail.key",
            "tlsCertFile": "restgomail.crt",
            "allowOnlyKnownCertificates": true,
            "knownCertificates": {
                "clientOneContainer": "@one_client.crt",
                "clientTwoContainer": "MIdjr6RfjfuESwekjEDffg..."
            },
            "waitSecondsAfterSmtpReq": 10,
            "debugMode": false
        }
    }

E-mail send request sample (JSON)
---------------------------------
This is a sample JSON passed as HTTP-POST request to the 443 port (configured above)

    {
        "sendmail": {
            "from": "sampleuser@gmail.com",
            "to": "tomyfriend@postmail.com",
            "subject": "This is a test message",
            "bodyhtml": "Probably <i>I can say</i> this thing is <h1>WORK'S FLAWLESS!</h1>!"
        }
    }

