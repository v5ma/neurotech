import argparse
import asyncio
import http.server
import os
import socketserver
import serial
import websockets


def offsetBinaryToInt(hexstr, offset):
    if len(hexstr) != 6:
        return -1
    not_encoded = int(hexstr, 16)  # Make a python integer out of the hex string
    if not_encoded & 1 << offset:  # Test if the most significant bit is on == this is a positive number
        encoded = not_encoded & ~(1 << offset)
    else:
        encoded = not_encoded - 2**offset
    return encoded


def adc2volts(i):
    return i * 5 / 2**23


def findPath():
    if args.path:
        return args.path
    basepath = "/dev/rfcomm"
    for i in range(5):
        testpath = basepath + str(i)
        if os.path.exists(testpath):
            return testpath


async def brains(websocket, path):
    """
    Example Brainduino serial stream:

    F80F19\t7801F4\rF80F19\t7801F4\rF80F19\t7801F4\r
    """
    print("brains")
    channels = (bytearray(6), bytearray(6))
    ctr = 0
    while True:
        bs = s.read(64)
        for b in bs:
            if b == ord('\r'):
                try:
                    c1, c2 = offsetBinaryToInt(channels[0], 23), offsetBinaryToInt(channels[1], 23)
                    c1, c2 = adc2volts(c1), adc2volts(c2)
                    await websocket.send("%s %s" % (c1, c2))
                except Exception:
                    await websocket.send("%s %s" % channels)
                ctr = 0
            elif b == ord('\t'):
                continue
            elif ctr < 6:
                channels[0][ctr] = b
                ctr += 1
            elif ctr >= 6:
                channels[1][ctr % 6] = b
                ctr += 1


async def brainSocketClient():
    async with websockets.connect('ws://localhost:8765') as websocket:
        brain(websocket, 5678)
    asyncio.get_event_loop().run_forever()


def main():
    asyncio.get_event_loop().run_until_complete(brainSocketClient())


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--path", help="path to the brainduino device")
    args = parser.parse_args()
    brainduino_path = findPath()
    s = serial.Serial(brainduino_path, baudrate=230400)
    if s.isOpen():
        print("Connected to brainduino at path=[%s]" % brainduino_path)
        s.write(bytes("S", "utf-8"))
        s.write(bytes("U", "utf-8"))
        main()
    else:
        print("Failed to connect to brainduino. Exiting now.")
