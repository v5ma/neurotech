import serial


def offsetBinaryToInt(hexstr, offset):
    # Perhaps we should test that the offset makes sense give the length of the hex string?
    if len(hexstr) != 6:
        return -1
    not_encoded = int(hexstr, 16) # Make a python integer out of the hex string
    if not_encoded & 1 << offset: # Test if the most significant bit is on == this is a positive number
        encoded = not_encoded & ~(1 << offset)
    else:
        encoded = not_encoded - 2**offset
    return encoded

def main():
    """
    offsetBinaryToInt("800001", 23) # 1
    offsetBinaryToInt("7fffff", 23) # -1
    offsetBinaryToInt("800000", 23) # 0
    offsetBinaryToInt("ffffff", 23) # max
    offsetBinaryToInt("000000", 23) # min
    """
    chan1, chan2 = '', ''
    c2rdy = False
    while True:
        data = s.read(32)
        for d in data:
            if d == "\t":
                c2rdy = True
                continue
            elif d == "\r":
                print("channel 1 = %s \t channel 2 = %s" % (offsetBinaryToInt(chan1, 23) * 5. / 2.**23, offsetBinaryToInt(chan2, 23) * 5. / 2.**23))
                chan1, chan2 = '', ''
                c2rdy = False
                continue
            if c2rdy:
                chan2 += d
            else:
                chan1 += d

if __name__ == '__main__':
    s = serial.Serial("/dev/ttyACM0", baudrate=230400)
    try:
        main()
    except Exception as e:
        print(e)
        s.close()
