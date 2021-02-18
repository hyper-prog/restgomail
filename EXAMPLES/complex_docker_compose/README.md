Using of this example
=====================

Run the create_certs.bash to generate certificates for server and client

	./create_certs.bash

Edit the data/server.json according to your settings

Start the container of restgomail

	docker-compose up -d

Now you can send messages to the container with the php sample

	php ./send_sample_request.php

