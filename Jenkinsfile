pipeline {
    agent any

    stages {
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }


        stage('Sync Deployments with rsync') {
            steps {
                script {
                    def envName = env.BRANCH_NAME.toUpperCase()
                    def envFileCredentialId = "ECOM_${envName}_ENV"
                    def composeUpCommand = "sudo docker-compose up --build -d ride_sharing_auth_db_${env.BRANCH_NAME} ride_sharing_auth_service_${env.BRANCH_NAME}"

                    withCredentials([ 
                        string(credentialsId: 'HOST_IP', variable: 'SERVER_HOST'),
                        string(credentialsId: 'SERVER_USER', variable: 'SERVER_USER'),
                        file(credentialsId: 'SERVER_KEY', variable: 'SSH_KEY_PATH'),
                        string(credentialsId: 'SERVER_PORT', variable: 'SSH_PORT'),
                        file(credentialsId: envFileCredentialId, variable: 'ENV_FILE')
                    ]) {
                        // Step 1: Create .env file
                        sh "ls -l \$ENV_FILE"
                        sh """
                            set -x
                            echo "Attempting to copy the env file..."
                            cp "$ENV_FILE" .env
                            echo "Successfully copied the env file."
                        """

                        // Step 2: Copy files to remote server using rsync
                        sh """
                            rsync -avz --delete --exclude='.git/' --exclude='.github/' -e "ssh -i \$SSH_KEY_PATH -p \$SSH_PORT -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" ./ .env \${SERVER_USER}@\${SERVER_HOST}:/home/\${SERVER_USER}/sudarshan/ride-sharing-auth/${env.BRANCH_NAME}/
                        """

                        // Step 3: SSH to start the docker containers
                        sh """
                            ssh -i \$SSH_KEY_PATH -p \$SSH_PORT -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \${SERVER_USER}@\${SERVER_HOST} "cd /home/\${SERVER_USER}/sudarshan/ride-sharing-auth/${env.BRANCH_NAME} && ${composeUpCommand}"
                        """

                        // Step 4: Prune unused docker images
                        sh """
                            ssh -i \$SSH_KEY_PATH -p \$SSH_PORT -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \${SERVER_USER}@\${SERVER_HOST} "sudo docker image prune -a --force"
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                commiterEmail = sh(script: "git show -s --format='%ae'", returnStdout: true).trim()
            }
            cleanWs()
        }
        failure {
            emailext body: '${DEFAULT_CONTENT}',
                to: commiterEmail, 
                subject: '${DEFAULT_SUBJECT}', 
                saveOutput: false
        }
    }
}
