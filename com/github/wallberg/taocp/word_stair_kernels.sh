#!/bin/bash

# Exercise 7.2.2.1-91: WORDS(1000) right word stairs
head -1000 ~/git/sandbox/taocp/assets/sgb-words.txt \
    | taocp ws kernels --right \
    | taocp xc --compact \
    | run word_stair_kernels.py --right

# Exercise 7.2.2.1-91: WORDS(500) left word stairs
head -500 ~/git/sandbox/taocp/assets/sgb-words.txt \
    | taocp ws kernels \
    | taocp xc --compact \
    | run word_stair_kernels.py

