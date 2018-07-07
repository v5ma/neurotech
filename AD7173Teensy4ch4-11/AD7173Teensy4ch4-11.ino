// Teensy 24bit 2ch (4ch) brain-duino with filter frequency control LED
// In Arduino IDE select Tools > Board > Teensy 3.1/3.2

#include <SPI.h>
#include <AD7173.h>
#include <Arduino.h>
#include <Audio.h>
#include <Wire.h>

#define HWSERIAL Serial1
/**** df 27 june 2018 - restored serial output (data)
* #define HWSERIAL Serial
****/

AudioSynthWaveformSine sine1;
AudioOutputAnalog dac1;
AudioConnection patchCord2 (sine1, dac1);

long freq = 12;			// 18000 is not good shape, 1000 is nice, use for test impedance 12Hz
unsigned long baud = 19200;
unsigned long baud1 = 230400;	// 460800


const int ledPin = 9;		// 13 // use for clock
const byte chipselect = 10;

// Create an IntervalTimer object 
IntervalTimer myTimer;
IntervalTimer myTimer1;

const int analogInPin = 4;	// Analog input pin A0 -->> A4
const int analogInPin2 = 5;	// Analog input pin A1 -->> A5
const int analogInPin3 = 6;	// Analog input pin A2 -->> A6
const int analogInPin4 = 7;	// Analog input pin A3 -->> A7

int sensorValue = 0;		// value read from the pot
int outputValue = 0;		// value output to the PWM (analog out)
int sensorValue2 = 0;		// value read from the pot
int outputValue2 = 0;		// value output to the PWM (analog out)
int sensorValue3 = 0;		// value read from the pot
int outputValue3 = 0;		// value output to the PWM (analog out)
int sensorValue4 = 0;		// value read from the pot
int outputValue4 = 0;		// value output to the PWM (analog out)
int sensorValue5 = 0;		// value read from the pot
int outputValue5 = 0;		// value output to the PWM (analog out)
int sensorValue6 = 0;		// value read from the pot
int outputValue6 = 0;		// value output to the PWM (analog out)

String channel1 = "abc";
String channel2 = "abc";
String channel3 = "abc";
String channel4 = "abc";
String channel5 = "abc";
String channel6 = "abc";
String zeroChar = "0";

byte buffer[20];
String stringOne;

int ledON;
int ledON2;

const int ledPin2 = 7;		// Teensy // 6 : Arduino // use for LED SW : filter mode, ON is use filter ( 30Hz, 40Hz, 50Hz , 100, 150Hz 200Hz )
const int ledPin4 = 4;		// Teensy . // 2 : Arduino
const int ledPin5 = 5;		// Teensy . // 3 : Arduino
const int ledPin6 = 6;		// Teensy . // 4 : Arduino
const int ledPin7 = 8;		// Teensy . // 5 : Arduino

String txtMsg = "";
boolean stringComplete = false;	// whether the string is complete
String inputString = "";

byte valueR[1];

int InputX = 0;
float sineP = 0.1;
long data;
int data1, data2, data3, data4;
int countD;
int dataX;
byte value0[3];
byte value1[3];
byte value2[3];
byte value3[3];
byte value4[3];
long loopcount;
long MaxCount;

// 2 = 2 channel system, 4 = 4 channel system
int Nchannel = 2;

// final DC offset Calibration, 1 = 0.5 µV

// NeurofoxKeid
// long CH0offset = -3254; // -4254
// long CH1offset = 2937; // 3837

// Libertas
long CH0offset = 3500; // 4,794
long CH1offset = -6000; // -7742

// NeurofoxArcturus
// long CH0offset = -10000;	// -13,895
// long CH1offset = 11000;		// 14527

// NeurofoxCursa
// long CH0offset = 4700; // 6,016
// long CH1offset = 6000; // 8560

long CH2offset = 0;
long CH3offset = 0;
long CH4offset = 0;

// Function definitions here:
void setupInput (int);
void setupAD7173 ();
void setupMyTimer (int);
void setupMyTimer1 (int);
void setupInputX (int);
void HEXprint (byte);
void HEXprint24bit (byte, byte, byte, byte, byte, byte);
void SimpleHEXprint24bit (byte, byte, byte, byte, byte, byte);
void HEXprint24bit4ch (byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte);
void HEXprint24bit5ch (byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte, byte);
void sendData (int);
void getData ();
void callbackX ();
void callback ();
void setupSampleF (int);
void impedanceCheck (int);

void setup ()
{

  pinMode (ledPin, OUTPUT);
  pinMode (ledPin2, OUTPUT);

  pinMode (ledPin4, OUTPUT);
  pinMode (ledPin5, OUTPUT);
  pinMode (ledPin6, OUTPUT);
  pinMode (ledPin7, OUTPUT);


  // this is for Teensy 3.2 : 72 MHz

  myTimer.begin (callback, 25);	// blinkLED to run 39965.625 Hz
  myTimer.priority (0);

  Serial.begin (baud);		// USB, communication to PC or Mac
  HWSERIAL.begin (baud1);	// communication to hardware serial

  
  // 29 July 2015
  //     analogReference(EXTERNAL);

  pinMode (chipselect, OUTPUT);
  digitalWrite (chipselect, LOW);

  delay (200);

  // setup for direct mode
  setupAD7173 ();


  /* wait for ADC */
  delay (200);

  // set as filter mode
  setupInput (0);		// 1:direct mode, not use filter, 0 filter


  digitalWrite (ledPin2, LOW);	// direct mode, not use filter

  // set OFF
  impedanceCheck (4);

}

void setupInput (int maxF)
{

  InputX = maxF;

  // 0 is filter, 1 is direct
  if (Nchannel == 2) {
    if (InputX == 1) {
      AD7173.set_channel_config (CH0, true, SETUP2, AIN8, AIN9);
      AD7173.set_channel_config (CH1, true, SETUP3, AIN10, AIN11);

      AD7173.set_setup_config (SETUP2, BIPOLAR);
      AD7173.set_setup_config (SETUP3, BIPOLAR);

      AD7173.set_offset_config (OFFSET2, CH2offset);
      AD7173.set_offset_config (OFFSET3, CH3offset);
    }
  else if (InputX == 0) {
    AD7173.set_channel_config (CH0, true, SETUP0, AIN4, AIN5);
    AD7173.set_channel_config (CH1, true, SETUP1, AIN6, AIN7);

    AD7173.set_setup_config (SETUP0, BIPOLAR);
    AD7173.set_setup_config (SETUP1, BIPOLAR);

    AD7173.set_offset_config (OFFSET0, CH0offset);
    AD7173.set_offset_config (OFFSET1, CH1offset);
    }
  }
  countD = 0;
}

void setupAD7173 ()
{
  /* initiate ADC, return true if device ID is valid */
  AD7173.init ();
  AD7173.sync ();
  /* reset ADC registers to the default state */
  AD7173.reset ();
  AD7173.sync ();

  /* check if the ID register of the ADC is valid */
  if (AD7173.is_valid_id ())
    Serial.println ("AD7173 ID is valid");
  else
    Serial.println ("AD7173 ID is invalid");

  /* set ADC input channel configuration */
  /* enable channel 0 and channel 1 and connect each to 2 analog inputs for bipolar input */
  /* CH0 - CH15 */
  /* true/false to enable/disable channel */
  /* SETUP0 - SETUP7 */
  /* AIN0 - AIN16 */


// setup for 2channle 24bit
//
  if (Nchannel == 2) {
    if (InputX == 1) {
      AD7173.set_channel_config (CH0, true, SETUP2, AIN8, AIN9);
      AD7173.set_channel_config (CH1, true, SETUP3, AIN10, AIN11);

      AD7173.set_setup_config (SETUP2, BIPOLAR);
      AD7173.set_setup_config (SETUP3, BIPOLAR);

      AD7173.set_offset_config (OFFSET2, CH2offset);
      AD7173.set_offset_config (OFFSET3, CH3offset);
    } else if (InputX == 0) {
	  AD7173.set_channel_config (CH0, true, SETUP0, AIN4, AIN5);
	  AD7173.set_channel_config (CH1, true, SETUP1, AIN6, AIN7);

	  AD7173.set_setup_config (SETUP0, BIPOLAR);
	  AD7173.set_setup_config (SETUP1, BIPOLAR);

	  AD7173.set_offset_config (OFFSET0, CH0offset);
	  AD7173.set_offset_config (OFFSET1, CH1offset);
    }
  }

// setup for 4channle 24bit
  else if (Nchannel == 4) {
    AD7173.set_channel_config (CH0, true, SETUP0, AIN4, AIN5);
    AD7173.set_channel_config (CH1, true, SETUP1, AIN6, AIN7);
    AD7173.set_channel_config (CH2, true, SETUP2, AIN8, AIN9);
    AD7173.set_channel_config (CH3, true, SETUP3, AIN10, AIN11);

    AD7173.set_setup_config (SETUP0, BIPOLAR);
    AD7173.set_setup_config (SETUP1, BIPOLAR);
    AD7173.set_setup_config (SETUP2, BIPOLAR);
    AD7173.set_setup_config (SETUP3, BIPOLAR);

    AD7173.set_offset_config (OFFSET0, CH0offset);
    AD7173.set_offset_config (OFFSET1, CH1offset);
    AD7173.set_offset_config (OFFSET2, CH2offset);
    AD7173.set_offset_config (OFFSET3, CH3offset);
  } else if (Nchannel == 5) {
    AD7173.set_channel_config (CH0, true, SETUP0, AIN4, AIN5);
    AD7173.set_channel_config (CH1, true, SETUP1, AIN6, AIN7);
    AD7173.set_channel_config (CH2, true, SETUP2, AIN8, AIN9);
    AD7173.set_channel_config (CH3, true, SETUP3, AIN10, AIN11);
    AD7173.set_channel_config (CH4, true, SETUP4, AIN12, AIN13);

    AD7173.set_setup_config (SETUP0, BIPOLAR);
    AD7173.set_setup_config (SETUP1, BIPOLAR);
    AD7173.set_setup_config (SETUP2, BIPOLAR);
    AD7173.set_setup_config (SETUP3, BIPOLAR);
    AD7173.set_setup_config (SETUP4, BIPOLAR);

    AD7173.set_offset_config (OFFSET0, CH0offset);
    AD7173.set_offset_config (OFFSET1, CH1offset);
    AD7173.set_offset_config (OFFSET2, CH2offset);
    AD7173.set_offset_config (OFFSET3, CH3offset);
    AD7173.set_offset_config (OFFSET4, CH4offset);
  }
  /*
     set the ADC SETUP0 coding mode to BIPLOAR output
     SETUP0 - SETUP7
     BIPOLAR, UNIPOLAR

     set ADC OFFSET0 offset value
     OFFSET0 - OFFSET7
    
     this function put gain setup with offset setup
     this Gain0 setup is for only one AD7173 : somethings Gain is strange.
     ->> AD7173.set_gain_config(GAIN0, 32768);
     byte valueX[3] = {0x00, 0x80, 0x00}; -->> 32768
     {0x00, 0x50, 0x00}; -->> 20480
     byte valueX[3] = {0x50, 0x00, 0x00}; -->> 5242880 --=> reset
     byte valueX[3] = {0x55, 0x55, 0x55}; -->> 5592405 --=> nominal
     {0x10, 0x00, 0x00}; -->> 1048576
     {0x01, 0x00, 0x00}; -->> 65536
     use this need modify 3 files : AD7173.h AD7173.cpp keywords.txt
    


     SPS_1007 2ch makes 500 Hz sampling par channel
     SPS_1007 4ch makes 250 Hz sampling par channel
     SPS_2597 need to use external Bluetooth deveice
     SPS_2597 2ch makes 1300 Hz sampling par channel
     SPS_2597 4ch makes 650 Hz sampling par channel

     set the ADC FILTER0 ac_rejection to false and samplingrate to 1007 Hz
     FILTER0 - FILTER7
     SPS_1, SPS_2, SPS_5, SPS_10, SPS_16, SPS_20, SPS_49, SPS_59, SPS_100, SPS_200
     SPS_381, SPS_503, SPS_1007, SPS_2597, SPS_5208, SPS_10417, SPS_15625, SPS_31250
  */

  if ((Nchannel == 4) || (Nchannel == 5)) {
      AD7173.set_filter_config (FILTER0, SPS_1007);
      AD7173.set_filter_config (FILTER1, SPS_1007);
      AD7173.set_filter_config (FILTER2, SPS_1007);
      AD7173.set_filter_config (FILTER3, SPS_1007);
      AD7173.set_filter_config (FILTER4, SPS_1007);
  } else {
      AD7173.set_filter_config (FILTER0, SPS_503);
      AD7173.set_filter_config (FILTER1, SPS_503);
      AD7173.set_filter_config (FILTER2, SPS_503);
      AD7173.set_filter_config (FILTER3, SPS_503);
      AD7173.set_filter_config (FILTER4, SPS_503);
  }

  /* 
     Set the ADC data and clock mode
     CONTINUOUS_CONVERSION_MODE, SINGLE_CONVERSION_MODE
     in SINGLE_CONVERSION_MODE after all setup channels are sampled the ADC goes into STANDBY_MODE
     to exit STANDBY_MODE use this same function to go into CONTINUOUS or SINGLE_CONVERSION_MODE
     INTERNAL_CLOCK, INTERNAL_CLOCK_OUTPUT, EXTERNAL_CLOCK_INPUT, EXTERNAL_CRYSTAL
     ex clock : 10min9sec, int clock : 10min7sec when set 500 Hz 10min
    
  */
  // AD7173.set_adc_mode_config(CONTINUOUS_CONVERSION_MODE, INTERNAL_CLOCK);
  AD7173.set_adc_mode_config (CONTINUOUS_CONVERSION_MODE, EXTERNAL_CRYSTAL);


  // enable or disable CONTINUOUS_READ_MODE, to exit use AD7173.reset(); */
  // AD7173.reset(); return all registers to default state, so everything has to be setup again */
  AD7173.set_interface_mode_config (false);

}


/*
Set up high-cut frequency
*/
void setupMyTimer (int speedA)
{
  myTimer.end ();
  myTimer.begin (callback, speedA);
  myTimer.priority (0);
}

/*
Set up LED blink speed -- high-cut frequency indicator
*/
void setupMyTimer1 (int speedA)
{
  ledON2 = 0;
  myTimer1.end ();
  myTimer1.begin (callbackX, speedA);
}

/*
Sets high-cut frequency + LED blink frequency based on cases below
*/
void setupInputX (int maxF)
{
  switch (maxF) {
    case 0:
      setupInput (0);
      setupMyTimer (160);	//  run 6270 -> 3134.8 Hz -->> 32Hz
      setupMyTimer1 (450000);
      break;
    case 1:
      setupInput (0);
      setupMyTimer (25);	//  run 39965.625 Hz -> 20000 -->> 200Hz
      setupMyTimer1 (16000);
      break;
    case 2:
      setupInput (0);
      setupMyTimer (33);	//  run 30000Hz -> 15000 -->> 150Hz
      setupMyTimer1 (33000);
      break;
    case 3:
      setupInput (0);
      setupMyTimer (50);	//   run 20000 -> 10000 -->> 100Hz
      setupMyTimer1 (50000);
      break;
    case 4:
      setupInput (0);
      setupMyTimer (100);	//   run 10000 -> 5000 -->> 50Hz
      setupMyTimer1 (100000);
      break;
    case 5:
      setupInput (0);
      setupMyTimer (128);	//   blinkLED  to run 8000 -> 4000 -->> 40Hz
      setupMyTimer1 (200000);
      break;
    case 6:
      setupInput (1);		// direct mode
      setupMyTimer (25);	//  run 39965.625 Hz -> 20000 -->> 200Hz
      myTimer1.end ();
      ledON2 = 2;
      digitalWrite (ledPin2, LOW);
      break;
    default:
      ledON2 = 2;
  }
}


/*
Ensure 2 ASCII characters are always printed
inval <= 00001000 print leading 0 for most significant four bits
*/
void HEXprint (byte inVal)
{
  if (inVal >= 16) {
    HWSERIAL.print (inVal, HEX);
  } else {
    HWSERIAL.print (0);
    HWSERIAL.print (inVal, HEX);
  }
}

/*
Protocol explained:
- Each channel has 24bits of information
- A hex value is 4bits
- Send 12bits (1.5bytes) with tab-delimiter followed by 12bits (1.5bytes)
|--------- Channel 1 --------|--------- Channel 2 --------|
Hex-Hex-Hex	Hex-Hex-Hex	Hex-Hex-Hex	Hex-Hex-Hex
- Encoding from ADC is offset binary where K = 2^23 gives a dynamic range of  -8388608 (all off) to 8322607 (all on)
- Each channel is converted to 2 tab-delimited, ASCII-Hex encoded, 12bit words to be interpreted as a continuous 24bit offset binary word
*/
void HEXprint24bit (byte inVal1, byte inVal2, byte inVal3, byte inVal4, byte inVal5, byte inVal6)
{
  HEXprint (inVal1);
  if (inVal2 >= 16) {
    HWSERIAL.print (inVal2 / 16, HEX); // 4 most significant bits
    HWSERIAL.write (0x09);
    HWSERIAL.print (inVal2 - ((inVal2 / 16) * 16), HEX); // 4 least significant bits
  } else {
    HWSERIAL.print (0);
    HWSERIAL.write (0x09);
    HWSERIAL.print (inVal2, HEX);
  }
  HEXprint (inVal3);

  HWSERIAL.write (0x09);

  HEXprint (inVal4);
  /*
  if (inVal4 >= 16) {
    HWSERIAL.print (inVal4, HEX);
  } else {
    HWSERIAL.print (0);
    HWSERIAL.print (inVal4, HEX);
  }
  */
  if (inVal5 >= 16) {
    HWSERIAL.print (inVal5 / 16, HEX);
    HWSERIAL.write (0x09);
    HWSERIAL.print (inVal5 - ((inVal5 / 16) * 16), HEX);
  } else {
    HWSERIAL.print (0);
    HWSERIAL.write (0x09);
    HWSERIAL.print (inVal5, HEX);
  }
  HEXprint (inVal6);
}

void SimpleHEXprint24bit (byte inVal1, byte inVal2, byte inVal3, byte inVal4, byte inVal5, byte inVal6)
{
  HEXprint (inVal1);
  HEXprint (inVal2);
  HEXprint (inVal3);
  HWSERIAL.write(0x09);
  HEXprint (inVal4);
  HEXprint (inVal5);
  HEXprint (inVal6);
}

void HEXprint24bit4ch (byte inVal1, byte inVal2, byte inVal3, byte inVal4,
		  byte inVal5, byte inVal6, byte inVal7, byte inVal8,
		  byte inVal9, byte inVal10, byte inVal11, byte inVal12)
{
  byte outStream[16];

  outStream[0] = inVal1;
  outStream[1] = inVal2;
  outStream[2] = inVal3;

  outStream[3] = inVal4;
  outStream[4] = inVal5;
  outStream[5] = inVal6;

  outStream[6] = inVal7;
  outStream[7] = inVal8;
  outStream[8] = inVal9;

  outStream[9] = inVal10;
  outStream[10] = inVal11;
  outStream[11] = inVal12;

  // time stamp, etc.
  outStream[12] = 0x55;
  outStream[13] = 0x44;
  outStream[14] = 0x33;

  outStream[15] = 0x0D;		// return

  HWSERIAL.write (outStream, 16);
}


void HEXprint24bit5ch (byte inVal1, byte inVal2, byte inVal3, byte inVal4,
		  byte inVal5, byte inVal6, byte inVal7, byte inVal8,
		  byte inVal9, byte inVal10, byte inVal11, byte inVal12,
		  byte inVal13, byte inVal14, byte inVal15)
{
  byte outStream[16];

  outStream[0] = inVal1;
  outStream[1] = inVal2;
  outStream[2] = inVal3;

  outStream[3] = inVal4;
  outStream[4] = inVal5;
  outStream[5] = inVal6;

  outStream[6] = inVal7;
  outStream[7] = inVal8;
  outStream[8] = inVal9;

  outStream[9] = inVal10;
  outStream[10] = inVal11;
  outStream[11] = inVal12;

  outStream[12] = inVal13;
  outStream[13] = inVal14;
  outStream[14] = inVal15;

  outStream[15] = 0x0D;		// return

  HWSERIAL.write (outStream, 16);
}

void sendData (int NchannelX)
{
  // 24 bit 5 ch
  if (NchannelX == 5) {
    HEXprint24bit5ch (value0[0], value0[1], value0[2], value1[0], value1[1],
    value1[2], value2[0], value2[1], value2[2], value3[0],
    value3[1], value3[2], value4[0], value4[1],
    value4[2]);
  } else if (NchannelX == 4) {
    // 24 bit 4 ch 
    HEXprint24bit4ch (value0[0], value0[1], value0[2], value1[0], value1[1],
    value1[2], value2[0], value2[1], value2[2], value3[0],
    value3[1], value3[2]);
  } else if (NchannelX == 2) {
    // 24 bit 2ch to 12 bit 4 ch
    SimpleHEXprint24bit (value0[0], value0[1], value0[2], value1[0], value1[1], value1[2]);
    HWSERIAL.print ("\r");
    // HEXprint24bit (value0[0], value0[1], value0[2], value1[0], value1[1], value1[2]);
    // HWSERIAL.print ("\r");
  }
}

void getData ()
{
  byte value[3];

  if (DATA_READY)
    {
      AD7173.get_data (value);

      countD = countD + 1;

      if (countD == 1)
	{
	  value0[0] = value[0];
	  value0[1] = value[1];
	  value0[2] = value[2];
	}

      if (countD == 2)
	{
	  value1[0] = value[0];
	  value1[1] = value[1];
	  value1[2] = value[2];
	}

      if (countD == 3)
	{
	  value2[0] = value[0];
	  value2[1] = value[1];
	  value2[2] = value[2];
	}

      if (countD == 4)
	{
	  value3[0] = value[0];
	  value3[1] = value[1];
	  value3[2] = value[2];
	}

      if (countD == 5)
	{
	  value4[0] = value[0];
	  value4[1] = value[1];
	  value4[2] = value[2];
	}


      if (countD >= Nchannel)
	{
	  countD = 0;
	  sendData (Nchannel);
	}

    }
}



/*
  use for high cut filter frequency and/or direct : OFF LED
  high speed blinking is high cutoff frequency
  low speed blinking is low cutoff frequency
*/
void callbackX (void)
{
  if (ledON2 == 1) {
    ledON2 = 0;
    digitalWrite (ledPin2, LOW);
  } else if (ledON2 == 0) {
    ledON2 = 1;
    digitalWrite (ledPin2, HIGH);
  }
}

/*
  this clock use for digital filter high cut frequency setup
  39965.625 Hz -> 20000 -->> 200Hz
*/
void callback (void)
{
  if (ledON == 1) {
    ledON = 0;
    digitalWrite (ledPin, LOW);
  } else {
    ledON = 1;
    digitalWrite (ledPin, HIGH);
  }
}


/* 
SPS_381, SPS_503, SPS_1007, SPS_2597, SPS_5208, SPS_10417, SPS_15625, SPS_31250
*/
void setupSampleF (int maxF)
{
  // 0 is SPS_1007 , 1 is SPS_2597
  if (maxF == 0) {
    AD7173.set_filter_config (FILTER0, SPS_503);
    AD7173.set_filter_config (FILTER1, SPS_503);
    AD7173.set_filter_config (FILTER2, SPS_503);
    AD7173.set_filter_config (FILTER3, SPS_503);
  } else if (maxF == 1) {
    AD7173.set_filter_config (FILTER0, SPS_1007);
    AD7173.set_filter_config (FILTER1, SPS_1007);
    AD7173.set_filter_config (FILTER2, SPS_1007);
    AD7173.set_filter_config (FILTER3, SPS_1007);
  } else if (maxF == 2) {
    AD7173.set_filter_config (FILTER0, SPS_381);
    AD7173.set_filter_config (FILTER1, SPS_381);
    AD7173.set_filter_config (FILTER2, SPS_381);
    AD7173.set_filter_config (FILTER3, SPS_381);
  }

  countD = 0;


  AD7173.set_offset_config (OFFSET0, CH0offset);
  AD7173.set_offset_config (OFFSET1, CH1offset);
  AD7173.set_offset_config (OFFSET2, CH2offset);
  AD7173.set_offset_config (OFFSET3, CH3offset);

  /*
    set the ADC data and clock mode
    set the ADC data and clock mode
    CONTINUOUS_CONVERSION_MODE, SINGLE_CONVERSION_MODE
    in SINGLE_CONVERSION_MODE after all setup channels are sampled the ADC goes into STANDBY_MODE
    to exit STANDBY_MODE use this same function to go into CONTINUOUS or SINGLE_CONVERSION_MODE
    INTERNAL_CLOCK, INTERNAL_CLOCK_OUTPUT, EXTERNAL_CLOCK_INPUT, EXTERNAL_CRYSTAL
  */

  AD7173.set_adc_mode_config (CONTINUOUS_CONVERSION_MODE, EXTERNAL_CRYSTAL);
  // enable or disable CONTINUOUS_READ_MODE, to exit use AD7173.reset();
  // AD7173.reset(); return all registers to default state, so everything has to be setup again
  AD7173.set_interface_mode_config (false);
}



void impedanceCheck (int CH)
{

  digitalWrite (ledPin4, LOW);
  digitalWrite (ledPin5, LOW);
  digitalWrite (ledPin6, LOW);
  digitalWrite (ledPin7, LOW);

  switch (CH)
    {
    case 0:
      sine1.amplitude (sineP);
      digitalWrite (ledPin4, HIGH);
      break;
    case 1:
      sine1.amplitude (sineP);
      digitalWrite (ledPin5, HIGH);
      break;
    case 2:
      sine1.amplitude (sineP);
      digitalWrite (ledPin6, HIGH);
      break;
    case 3:
      sine1.amplitude (sineP);
      digitalWrite (ledPin7, HIGH);
      break;
    case 4:
      sine1.amplitude (0.0);
      digitalWrite (ledPin4, LOW);
      digitalWrite (ledPin5, LOW);
      digitalWrite (ledPin6, LOW);
      digitalWrite (ledPin7, LOW);
      break;

    default:
      digitalWrite (ledPin4, LOW);
      digitalWrite (ledPin5, LOW);
      digitalWrite (ledPin6, LOW);
      digitalWrite (ledPin7, LOW);
    }

  setupInput (InputX);

}



int nothings;


void loop ()
{
  char inChar = '0';

  getData ();

  if (HWSERIAL.available ()) {
    // get the new byte:
    inChar = (char) HWSERIAL.read ();

    // if the incoming character is a newline, set a flag
    // so the main loop can do something about it:

    // B is problem (-CH1 ) -->> T

    if ((inChar == 'S') || (inChar == 'I') || (inChar == 'P')
	|| (inChar == 'O') || (inChar == 'Q') || (inChar == 'Z')
	|| (inChar == 'M') || (inChar == 'V') || (inChar == 'W')
	|| (inChar == 'A') || (inChar == 'T') || (inChar == 'C')
	|| (inChar == 'D') || (inChar == 'E') || (inChar == 'F')
	|| (inChar == 'G') || (inChar == 'H') || (inChar == 'J')
	|| (inChar == 'K') || (inChar == 'X') || (inChar == 'Y')
	|| (inChar == 'U') || (inChar == 'B')) {
      stringComplete = true;
      Serial.println ("command input ");
    }

  }

  // print the message and a notice if it's changed:
  if (stringComplete)
    {

      switch (inChar)
	{
	case 'S':
	  // 32 Hz cut // in IBVA :0

	  HWSERIAL.print ("sssssssssssssss");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "sssssssssssssss" );
	  setupInputX (0);	// fiter frequency : 0 is 32Hz
	  break;
	case 'I':
	  // direct : go down from 200, 300, 400, 500, 600 // in IBVA :6

	  HWSERIAL.print ("iiiiiiiiiiiiiii");
	  HWSERIAL.print ("\r");
	  //     Serial.println( "iiiiiiiiiiiiiii" );
	  setupInputX (6);	// fiter frequency : 1 is 200Hz 
	  break;
	case 'P':
	  // 200 Hz cut // in IBVA :5

	  HWSERIAL.print ("ppppppppppppppp");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "ppppppppppppppp" );
	  setupInputX (1);	// fiter frequency : 1 is 200Hz
	  break;
	case 'O':
	  // 150 Hz cut // in IBVA :4

	  HWSERIAL.print ("ooooooooooooooo");
	  HWSERIAL.print ("\r");
	  //      Serial.println( "ooooooooooooooo" );
	  setupInputX (2);	// fiter frequency : 2 is 150Hz
	  break;
	case 'Q':
	  // 100 Hz cut // in IBVA :3

	  HWSERIAL.print ("qqqqqqqqqqqqqqq");
	  HWSERIAL.print ("\r");
	  //    Serial.println( "qqqqqqqqqqqqqqq" );
	  setupInputX (3);	// fiter frequency : 3 is 100Hz
	  break;
	case 'Z':
	  // 50 Hz cut // in IBVA :2

	  HWSERIAL.print ("zzzzzzzzzzzzzzz");
	  HWSERIAL.print ("\r");
	  //    Serial.println( "zzzzzzzzzzzzzzz" );
	  setupInputX (4);	// fiter frequency : 4 is 50Hz     
	  break;
	case 'M':
	  // 40 Hz cut // in IBVA :1

	  HWSERIAL.print ("mmmmmmmmmmmmmmm");
	  HWSERIAL.print ("\r");
	  //    Serial.println( "mmmmmmmmmmmmmmm" );
	  setupInputX (5);	// fiter frequency : 5 is 40Hz
	  break;
	case 'V':
	  // setup 500 Hz sampling
	  setupSampleF (0);

	  HWSERIAL.print ("vvvvvvvvvvvvvvv");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "vvvvvvvvvvvvvvv" );

	  break;
	case 'W':
	  // setup 1.3 KHz sampling  -->> 250 Hz
	  setupSampleF (1);

	  HWSERIAL.print ("wwwwwwwwwwwwwww");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "wwwwwwwwwwwwwww" );

	  break;
	case 'B':
	  // setup 1.3 KHz sampling  -->> 250 Hz
	  setupSampleF (2);

	  HWSERIAL.print ("bbbbbbbbbbbbbbb");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "bbbbbbbbbbbbbbb" );

	  break;
	case 'A':
	  // setup impedance check +1ch

	  HWSERIAL.print ("aaaaaaaaaaaaaaa");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "aaaaaaaaaaaaaaa" );
	  impedanceCheck (0);
	  break;
	case 'T':
	  // setup impedance check -1ch

	  HWSERIAL.print ("ttttttttttttttt");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "ttttttttttttttt" );
	  impedanceCheck (1);
	  break;
	case 'C':
	  // setup impedance check +2ch

	  HWSERIAL.print ("ccccccccccccccc");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "ccccccccccccccc" );
	  impedanceCheck (2);
	  break;
	case 'D':
	  // setup impedance check -2ch

	  HWSERIAL.print ("ddddddddddddddd");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "ddddddddddddddd" );
	  impedanceCheck (3);
	  break;
	case 'E':
	  // setup impedance check OFF

	  HWSERIAL.print ("eeeeeeeeeeeeeee");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "eeeeeeeeeeeeeee" );
	  impedanceCheck (4);
	  break;
	case 'F':
	  // testOSC 10Hz -> 12

	  sine1.frequency (12);

	  HWSERIAL.print ("fffffffffffffff");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "fffffffffffffff" );
	  setupInput (InputX);
	  break;
	case 'G':
	  // testOSC 30Hz -> 28

	  sine1.frequency (28);

	  HWSERIAL.print ("ggggggggggggggg");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "ggggggggggggggg" );
	  setupInput (InputX);
	  break;
	case 'H':

	  // testOSC 100Hz  -> 111

	  sine1.frequency (111);

	  HWSERIAL.print ("hhhhhhhhhhhhhhh");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "hhhhhhhhhhhhhhh" );
	  setupInput (InputX);
	  break;
	case 'J':
	  // power 0.1
	  sineP = 0.1;
	  sine1.amplitude (sineP);	// max 1

	  HWSERIAL.print ("jjjjjjjjjjjjjjj");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "jjjjjjjjjjjjjjj" );
	  setupInput (InputX);
	  break;
	case 'K':
	  // power 1.0
	  sineP = 1.0;
	  sine1.amplitude (sineP);	// max 1

	  HWSERIAL.print ("kkkkkkkkkkkkkkk");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "kkkkkkkkkkkkkkk" );
	  setupInput (InputX);
	  break;
	case 'Y':
	  Nchannel = 4;		// set 4 channel mode
	  setupSampleF (1);

	  HWSERIAL.print ("yyyyyyyyyyyyyyy");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "yyyyyyyyyyyyyyy" );
	  setupInput (InputX);
	  break;
	case 'X':
	  Nchannel = 2;		// set 2 channel mode, 500Hz
	  setupSampleF (1);
	  HWSERIAL.print ("xxxxxxxxxxxxxxx");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "xxxxxxxxxxxxxxx" );
	  setupInput (InputX);
	  break;
	case 'U':
	  Nchannel = 2;		// set 2 channel mode, 250HZ, Y is 4 channel mode, X is 2 channel mode
	  setupSampleF (0);
	  HWSERIAL.print ("uuuuuuuuuuuuuuu");
	  HWSERIAL.print ("\r");
	  //   Serial.println( "uuuuuuuuuuuuuuu" );
	  setupInput (InputX);
	  break;

	default:
	  nothings = 0;
	}

      stringComplete = false;

    }

}
