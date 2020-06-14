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
        stage('Compile') {
            steps {
                script {
                    /**
                     * To be able to access this Jenkins `tool` the https://wiki.jenkins.io/display/JENKINS/Go+Plugin plugin is needed.
                     * With more recent versions of Jenkins the documentation for adding a `go` installation is out of date. To properly
                     * configure a go installation go to the Jenkins tools configuration (Manage Jenkins -> Global Tool Configuration)
                     * find the "Go" and "Go installations" section and click "Add Go". The `name` specified below should
                     * line up with the "Go installation" to be used.
                     */
                    
                    def root = tool name: 'Go Version 1.14.4', type: 'go'

                    /**
                     * Add in GOPATH, GOROOT, GOBIN to the environment and add go to the path for jenkins.
                     * The environment variables are needed for glide to properly work and adding go to the path is required to
                     */
                    withEnv(["GOPATH=${env.WORKSPACE}/go", "GOROOT=${root}", "GOBIN=${root}/bin", "PATH+GO=${root}/bin"]) {
                        sh "mkdir -p ${env.WORKSPACE}/go/src"

                        echo "Branch is ${BRANCH_NAME} ..."
                        
                        sh '''
                        echo "PATH = ${PATH}"
                        echo "BRANCH_NAME = ${BRANCH_NAME}"
                        go build
                        '''

                    }
                }
            }
        }
        stage ('Package') {
			when {
			    not {
			        branch 'master'
			    }
			}
            steps {
                sh '''
                jar -cvf AvWxServer-${BRANCH_NAME}.jar application.properties AvWxServer
                '''
            }
		}

    }
}