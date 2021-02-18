<?php

$json =
    '{
      "sendmail": {
        "from": "sampleuser@gmail.com",
        "to": "target.address@mailbox.com",
        "subject": "Test from PHP",
        "bodyhtml": "This is received through <h1>RestGoMail</h1>"
      }
    }';

$r = SendHttpPostWithCert("https://127.0.0.1:44325/sendmail",$json,"rgmclient.crt","rgmclient.key");
print("Received:".$r);

function SendHttpPostWithCert($url, $postdata, $certfile, $keyfile)
{
    $curl = curl_init();

    curl_setopt($curl, CURLOPT_POST, 1);
    curl_setopt($curl, CURLOPT_POSTFIELDS, $postdata);

    curl_setopt($curl, CURLOPT_SSLKEY, $keyfile);
    curl_setopt($curl, CURLOPT_SSLCERT, $certfile);

    curl_setopt($curl, CURLOPT_CERTINFO, true);

    curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, FALSE);
    curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, FALSE);

    curl_setopt($curl, CURLOPT_CONNECTTIMEOUT, 15);
    curl_setopt($curl, CURLOPT_TIMEOUT, 60);

    curl_setopt($curl, CURLOPT_URL, $url);
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);

    $result = curl_exec($curl);
    $info = curl_getinfo($curl,CURLINFO_CERTINFO);

    print "You can check this certificate:"
    print "\n=============================\n";
    $c = $info[0]['Cert'];
    $c = preg_replace('/\-+BEGIN CERTIFICATE\-+/','',$c);
    $c = preg_replace('/\-+END CERTIFICATE\-+/','',$c);
    $c = trim($c);
    print $c;
    print "\n=============================\n";

    curl_close($curl);

    return $result;
}

