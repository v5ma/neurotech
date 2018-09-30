Operating instructions
======================

```
cd ~/go/src/github.com/Micah1/neurotech/brainduino-server
go build
./brainduino-server
```

For details see `./brainduino-server --help`

All files in `~/go/src/github.com/Micah1/neurotech/brainduino-server/static` are served upon request by the brainduino-server program. For example, when the browser issues the following request `GET localhost:8080/static/index.html`, brainduino-server will serve the index.html file. `GET localhost:8080` does the same thing. `GET localhost:8080/chartsngraphs` is equivalent to `GET localhost:8080/static/chartsngraphs.html`.

For writting commands to the braindunio try `POST localhost:8080/command/S`. For a full list of brainduino commands see (https://github.com/Micah1/neurotech).
