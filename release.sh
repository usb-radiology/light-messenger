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
  print_command_header "push changes"
  git push origin
  check_command_success "push"
  
  tag=$tag_arg
  print_command_header "tagging $tag"
  git tag $tag
  check_command_success "creating git tag"
  
  print_command_header "pushing git tag"
  git push origin $tag
  check_command_success "pushing git tag"
  
  print_command_header "run tests"
  make test
  check_command_success "tests"
  
  print_command_header "build"
  make build
  check_command_success "build"
  
  print_command_header "stop service"
  ssh nofasy@radmon 'sudo systemctl stop light-messenger'
  check_command_success "stop current service"
  
  print_command_header "copy binary"
  scp ./light-messenger.exec nofasy@radmon:~/light-messenger-dist/
  check_command_success "scp"
  
  print_command_header "start service"
  ssh nofasy@radmon 'sudo systemctl start light-messenger'
  check_command_success "start service"
  
  print_command_header "status service"
  ssh nofasy@radmon 'sudo systemctl status light-messenger'
  check_command_success "status service"
  
else
  printf "require tag name\n"
fi
