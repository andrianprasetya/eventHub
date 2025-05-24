pipeline {
    agent any

    triggers {
        githubPush()  // <- ini trigger otomatis dari push GitHub
    }

    stages {
        stage('Checkout') {
            steps {
                git credentialsId: 'github-creds', url: 'https://github.com/andrianprasetya/eventHub.git', branch: 'master'
            }
        }

        // Tambahkan stage-stage berikutnya seperti build, ssh, deploy
    }
}