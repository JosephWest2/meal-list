#!/usr/bin/env python3
import subprocess

subprocess.run(["python", "build.py"], shell=True)
subprocess.run(["go", "run", "main.go"], shell=True)