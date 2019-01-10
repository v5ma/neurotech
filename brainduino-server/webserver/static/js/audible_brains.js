let audioContext = new (window.AudioContext || window.webkitAudioContext);
let masterGainNode = null;
let customWaveform = null;
let sineTerms = null;
let cosineTerms = null;
let volumeControl = document.querySelector("input[name='volume']");
let playing = false;
let playButton = document.querySelector(".play");
playButton.onclick = togglePlay;
function setup() {
	volumeControl.addEventListener("change", changeVolume, false);
	masterGainNode = audioContext.createGain();
	masterGainNode.connect(audioContext.destination);
	masterGainNode.gain.value = volumeControl.value;
	sineTerms = new Float32Array([0, 0, 1, 0, 1]);
	cosineTerms = new Float32Array(sineTerms.length);
	customWaveform = audioContext.createPeriodicWave(cosineTerms, sineTerms);
}
function changeVolume(event) {
	masterGainNode.gain.value = volumeControl.value;
}
function playTone(freq) {
	let osc = audioContext.createOscillator();
	osc.connect(masterGainNode);
	osc.type = 'sine';
	osc.frequency.value = freq;
	osc.start();
		 
	return osc;
}
setup();
var accordium = [];
for (i=0; i<5; i++) {
	accordium[i] = [];
}
function togglePlay() {
	if(playing) {
		for (acc of accordium) {
			acc.stop();
		}
		playing = false;
	} else {
		accordium[0] = playTone(220.0);
		accordium[1] = playTone(330.0);
		accordium[2] = playTone(440.0);
		accordium[3] = playTone(660.0);
		accordium[4] = playTone(880.0);
		playing = true;
	}
}
