#!/bin/bash

# IT'S A HELPER SCRIPT THAT I USE TO RUN LOCAL TELEGRAM API

set -a        
source $PWD/.env
set +a
cd telegram-bot-api/workdir
../bin/telegram-bot-api
