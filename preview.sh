#!/bin/bash

FILENAME=$(cut -d" " -f1 <<< $1);
START_LINE=$(cut -d" " -f3 <<< $1);
END_LINE=$(cut -d" " -f4 <<< $1);
bat --style=numbers --color=always -r $START_LINE:$END_LINE $FILENAME
