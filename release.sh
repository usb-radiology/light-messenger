#!/bin/bash

usage () {
    cat <<EOM
Usage:
$(basename $0) [tag-name] [deploy]

EOM
#    exit 0
}

print_command_header () {
  printf "\n##\n# LIGHT-MESSENGER: $1\n##\n"
}

check_command_success () {
  if [ $? -ne 0 ]
  then
    print_command_header "$1 failed"
    exit 1
  fi
}

release_tag () {
  tag=$1

  print_command_header "push changes"
  git push origin
  check_command_success "push"
  
  print_command_header "tagging $tag"
  git tag $tag
  check_command_success "creating git tag"
  
  print_command_header "pushing git tag"
  git push origin $tag
  check_command_success "pushing git tag"
}

release_deploy () {
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
}

##
# - check for any modified / untracked files
# - check that all tests pass
# - check release tag is provided and tag release if yes
# - check if deploy tag is provided and deploy if yes
##

if [[ $(git diff --stat) != '' ]]; then
  printf "modified / untracked files\n"
  exit 1
fi

print_command_header "run tests"
make test
check_command_success "tests"

arg_tag=$1
arg_deploy=$2

if [[ $arg_tag ]]
then
  release_tag "$arg_tag"

  print_command_header "build"
  make build
  check_command_success "build"

  if [[ $arg_deploy ]]
  then
    release_deploy
  else
    printf "No deploy tag provided\n"
  fi

else
  printf "require tag name\n"
fi

printf "done\n"
