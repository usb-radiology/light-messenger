pipeline {
    agent any
    
    /*
    environment {
        CI='true'
    }
    */

    tools {
        go 'go-1.12.1'
    }

    options {
      disableConcurrentBuilds()
    }

    stages {
      stage('Test') {
        steps {
          // for rice
          withEnv(["PATH+EXTRA=${HOME}/go/bin"]){
            sh "/usr/bin/docker-compose -f docker-compose.yml up -d --force-recreate"
            sh 'until nc -z localhost 3311; do sleep 1; echo "Waiting for DB to come up..."; done'
            sh 'sleep 10'
            sh 'cp config-sample.json config.json'
            sh 'go get github.com/GeertJohan/go.rice/rice'
            sh 'make build'
            sh './light-messenger.exec db-exec --script-path ./res/drop_tables.sql'
            sh './light-messenger.exec db-exec --script-path ./res/create_tables.sql'
            sh 'make test'
            // sh "/usr/bin/docker-compose -f docker-compose.yml down -v"
          }
        }
      }
    }

    post {
      always {
        sh '/usr/bin/docker-compose rm -f -s'
        sh '/usr/bin/docker-compose down --rmi local --remove-orphans'
        
        echo 'email job results'
        emailext (
          body: "${currentBuild.currentResult}: Job ${env.JOB_NAME} build ${env.BUILD_NUMBER}\n More info at: ${env.BUILD_URL}",
          recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'RequesterRecipientProvider']],
          subject: "Jenkins Build ${currentBuild.currentResult}: Job ${env.JOB_NAME}, Build ${env.BUILD_NUMBER}",
          to: "victor.parmar@usb.ch, joshy.cyriac@usb.ch"
        )
      }
    }
}
