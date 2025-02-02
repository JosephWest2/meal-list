#!/usr/bin/env python3
import subprocess
import shutil

if shutil.which("go") is None:
    print("Please install go")
    exit(1)

if shutil.which("air") is None:
    print("Installing air")
    subprocess.run(["go", "install", "github.com/air-verse/air@latest"], shell=True)

print("Running air")
subprocess.run(["air", "-c", ".air.toml"], shell=True)