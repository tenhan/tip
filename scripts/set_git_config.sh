#!/bin/sh
git config pull.ff only
git config core.safecrlf true
git config core.autocrlf input
git config commit.gpgsign true
git config --list --local