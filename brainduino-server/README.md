Operating instructions
======================

Update code
```
cd ~/go/src/github.com/Micah1/neurotech/
git pull origin master
dep ensure
```

Build and run webserver
```
cd ~/go/src/github.com/Micah1/neurotech/brainduino-server/webserver
go build
sudo ./webserver --url 0.0.0.0:80
```

Build and run Brainduino
```
cd ~/go/src/github.com/Micah1/neurotech/brainduino-server/brainduino
go build
sudo ./brainduino --url 0.0.0.0:80 # alternatively sudo ./brainduino --url ml.sensorium.space
```

index.html location
`~/go/src/github.com/Micah1/neurotech/brainduino-server/webserver/static/views/index.html`

To commit changes to github at the end of the night
```
git add . --all
git commit -m "<what's in the commit>"
git push origin master
```

To inspect changes before committing try running `git status` or `git log`


Event Schema
============
```
type Sample struct {
        Name           string
        Channels       []float64
        Timestamp      time.Time
        SequenceNumber uint
}
Example websocket event JSON:
{
  "data": {
    "Name": "sample",
    "Channels": [0.01, 0.32],
    "Timestamp": "2018-10-10T21:05:05.031850444-07:00",
    "SequenceNumber": 1
  }
}

type FFTData struct {
        Name           string
        Channels       [][]float64
        Timestamp      time.Time
        SequenceNumber uint
}
Example websocket event JSON:
{
  "data": {
    "Name": "fft",
    "Channels": [[0.01, 0.32, ..., 0.123], [0.09234, 0.1234123, ..., 0.123543]],
    "Timestamp": "2018-10-10T21:05:05.031850444-07:00",
    "SequenceNumber": 1
  }
}
```


For details see `./brainduino --help` or `./webserver --help`

All files in `~/go/src/github.com/Micah1/neurotech/brainduino-server/webserver/static/views` are served upon request by the brainduino-server program. For example, when the browser issues the following request `GET localhost:8080/static/index.html`, brainduino-server will serve the index.html file. `GET localhost:8080` does the same thing. `GET localhost:8080/chartsngraphs` is equivalent to `GET localhost:8080/static/chartsngraphs.html`.

For writting commands to the braindunio try `POST localhost:8080/command/S`. For a full list of brainduino commands see https://github.com/Micah1/neurotech.
