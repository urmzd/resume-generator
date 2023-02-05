#!/usr/bin/env bash

go build github.com/urmzd/generate-resumes
./generate-resumes config.m.toml > ./assets/templates/m/resume.tex
cd ./assets/templates/m
pdflatex resume.tex
gio open resume.pdf
