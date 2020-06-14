pipeline {
    agent any
    // tools { 
    //     maven 'Maven 3.5.2' 
    //     jdk 'jdk-13.0.1+9' 
    // }
    environment {
        GO111MODULE = 'on'
    }
    stages {
        stage ('Compile') {
            steps {
                echo "Branch is ${BRANCH_NAME} ..."
                def root = tool name: 'Go Version 1.14.4', type: 'go'
                sh '''
                echo "PATH = ${PATH}"
                echo "BRANCH_NAME = ${BRANCH_NAME}"
                go build
                '''
                // withNPM(npmrcConfig:'my-custom-npmrc') {
                //     echo "Performing npm build..."
                //     sh 'npm install'
                // }
            }
        }
    }
}