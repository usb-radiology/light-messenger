#!/bin/bash

usage () {
    cat <<EOM
Usage:
$(basename $0) [tag-name]

EOM
#    exit 0
}

print_command_header () {
  printf "\n##\n# RELEASE: $1\n##\n"
}

check_command_success () {
  if [ $? -ne 0 ]
  then
    print_command_header "$1 failed"
    exit 1
  fi
}

if [[ $(git diff --stat) != '' ]]; then
  echo "modified / untracked files\n\n"
  exit 1
fi

tag_arg=$1

if [[ $tag_arg ]]
then
  tag=$tag_arg
  print_command_header "tagging $tag"
  git tag $tag
  check_command_success "creating git tag"
  git push origin $tag
  check_command_success "pushing git tag"
  make test
  check_command_success "tests"
  make build
  check_command_success "build"
  ssh nofasy@radmon 'sudo systemctl stop light-messenger'
  check_command_success "stop current service"
  scp ./light-messenger.exec nofasy@radmon:~/light-messenger-dist/
  check_command_success "scp"
  ssh nofasy@radmon 'sudo systemctl start light-messenger'
  check_command_success "start service"
  ssh nofasy@radmon 'sudo systemctl status light-messenger'
else
  printf "require tag name\n"
fi
