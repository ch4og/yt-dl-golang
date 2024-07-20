#!/bin/bash
set -a        
source $PWD/.env
set +a
cd telegram-bot-api/workdir
../bin/telegram-bot-api