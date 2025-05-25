pipeline {
    agent any

    environment {
        // Environment Variables
        GO_VERSION = "1.22.5"
        APP_NAME = "myapp"
        BUILD_DIR = "build"
        ENV = "production" // bisa diubah ke dev, staging, dll
          // Set PATH ke lokasi Go manual
        GOROOT = "/usr/local/go"
        GOPATH = "${env.WORKSPACE}/go"
        PATH = "${env.GOROOT}/bin:${env.GOPATH}/bin:${env.PATH}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

    stage('Verify Go') {
                steps {
                    sh 'go version'
                }
            }

        stage('Setup') {
            steps {
                sh 'go mod download'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./... -v -coverprofile=coverage.out'
            }
        }

        stage('Build') {
            steps {
                sh """
                    mkdir -p ${BUILD_DIR}
                    go build -o ${BUILD_DIR}/${APP_NAME} .
                """
            }
        }

        stage('Archive Build') {
            steps {
                archiveArtifacts artifacts: "${BUILD_DIR}/${APP_NAME}", fingerprint: true
            }
        }

        stage('Deploy') {
            when {
                branch 'master'
            }
            steps {
                echo "Deploying ${APP_NAME} to ${ENV} environment..."
                // Contoh deploy ke remote server via SSH
                sh """
                    scp ${BUILD_DIR}/${APP_NAME} root@165.22.63.86:/opt/myapp/${APP_NAME}
                    ssh root@165.22.63.86 'systemctl restart ${APP_NAME}'
                """
            }
        }
    }

    post {
        always {
            echo 'Cleaning up...'
            cleanWs()
        }
        success {
            echo '✅ Build and deployment succeeded!'
        }
        failure {
            echo '❌ Build or deployment failed!'
        }
    }
}