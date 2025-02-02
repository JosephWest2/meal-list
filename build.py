#!/usr/bin/env python3
import subprocess
import shutil

if shutil.which("go") is None:
    print("Please install go")
    exit(1)

if shutil.which("templ") is None:
    print("Installing templ")
    subprocess.run(["go", "install", "github.com/a-h/templ/cmd/templ@latest"], shell=True)

print("Building")
subprocess.run(["templ", "generate"], shell=True)