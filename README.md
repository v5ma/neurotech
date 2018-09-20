Wiki
====

https://www.noisebridge.net/wiki/DreamTeam/Brainduinov2

Commands
========

Commands are sent by write ASCII characters to the bluetooth serial device. The brainduino echos back the lower-case compliment of the command sent 15 times followed by a return carridge. For example:

```
Client writes 'S' to bluetooth serial device
If the command is recieved by the brainduino, the brainduino writes 'sssssssssssssss' + '\r' back over bluetooth

```
ASCII | Result
---   | ---
I     | Bypass the MAX7480 (8th-order butterworth filter)
P     | Set 200HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
O     | Set 150HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
Q     | Set 100HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
Z     | Set 50HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
M     | Set 40HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
S     | Set 32HZ low-pass corner frequency of MAX7480 (8th-order butterworth filter)
V     | Set sample rate to 500HZ
W     | Set sample rate to 250HZ
B     | Set sample rate to 190HZ
A     | Impedance check on with +1ch
T     | Impedance check on with -1ch
C     | Impedance check on with +2ch
D     | Impedance check on with -2ch
E     | Impedance check off
F     | Generate test signal with frequency = 12HZ
G     | Generate test signal with frequency = 28HZ
H     | Generate test signal with frequency = 111HZ
J     | Generate test signal with power = 0.1
K     | Generate test signal with power = 1
Y     | Set 4 channel mode
X     | Set 2 channel mode with 500HZ sample rate
U     | Set 2 channel mode with 250HZ sample rate
