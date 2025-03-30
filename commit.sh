#! /bin/bash

git add .

# take all args and pass them to git commit, fall back on "checkpoint" otherwise

message="checkpoint"
if [ -n "$*" ]; then
    message="$*"
fi

git commit -m "$message"

