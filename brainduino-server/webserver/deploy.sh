go build
rm /srv/www/brainduino/bin/webserver
cp webserver /srv/www/brainduino/bin/
rm -rf /srv/www/brainduino/static
cp -r static /srv/www/brainduino/
