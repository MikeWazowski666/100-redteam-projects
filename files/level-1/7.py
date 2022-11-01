"""Simple ROT cipher."""

from string import ascii_lowercase
from sys import argv

def encode(message: str, shift: int = 13) -> str:
    """
    Encode a message using a ROT cipher.

    :param message: message to be encoded
    :param shift: shift for encoding
    :return: encoded message
    """
    enc_msg = []
    letters = []
    alphabet = list(ascii_lowercase)
    alphabet_len = len(alphabet)
    letters[:0] = message  # get letters in message

    for letter in letters:
        if letter not in alphabet:  # skip the letter if not in alphabet
            enc_msg.append(letter)
            continue
        shift_letter = (alphabet.index(letter) + shift) % alphabet_len  # find the place of the letter and add shift
        enc_msg.append(alphabet[shift_letter])
    enc_msg = ''.join(enc_msg)
    return enc_msg

if __name__ == "__main__":
    if len(argv) == 2:
        print(encode(argv[1]))
    elif len(argv) == 3:
        print(encode(argv[1], int(argv[2])))  
    else:
        print(f"usage: {argv[0]} message shift")