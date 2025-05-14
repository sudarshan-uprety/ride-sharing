pipeline {
    agent any

    environment {
        // Docker image settings
        DOCKER_IMAGE_NAME = "app"
        DOCKER_IMAGE_TAG = "${env.BRANCH_NAME}"
        DOCKER_TAR_FILE = "${DOCKER_IMAGE_NAME}-${DOCKER_IMAGE_TAG}.tar"
    }

    stages {
        stage('Setup Environment') {
            steps {
                script {
                    // Simple environment setup based on branch
                    envName = env.BRANCH_NAME.toUpperCase()
                    IS_PROD = (env.BRANCH_NAME == "prod")
                    
                    // Get the environment file credential ID based on branch
                    envFileCredentialId = "RIDE_AUTH_${envName}_ENV"
                    
                    echo "Using environment file from credential: ${envFileCredentialId}"
                }
            }
        }

        stage('Checkout Code') {
            steps { 
                checkout scm 
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // Fetch the environment file from Jenkins credentials
                    withCredentials([file(credentialsId: envFileCredentialId, variable: 'ENV_FILE')]) {
                        // Copy the environment file to the workspace
                        sh "cp \"${ENV_FILE}\" .env"
                        
                        sh "echo 'Env file exists:' && ls -la .env"

                        sh "cat .env"
                        
                        // Build the Docker image with the environment file
                        sh """
                            export APP_ENV=${IS_PROD ? 'production' : 'development'}
                            docker-compose build app
                        """
                        
                        // Save the Docker image to a tar file
                        sh "docker save -o ${DOCKER_TAR_FILE} ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"
                    }
                }
            }
        }

        stage('Transfer Docker Image') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'HOST_IP', variable: 'SERVER_HOST'),
                        string(credentialsId: 'SERVER_USER', variable: 'SERVER_USER'),
                        file(credentialsId: 'SERVER_KEY', variable: 'SSH_KEY_PATH'),
                        string(credentialsId: 'SERVER_PORT', variable: 'SSH_PORT')
                    ]) {
                        // Transfer Docker image using SCP
                        sh """
                            scp -i "${SSH_KEY_PATH}" -P "${SSH_PORT}" -o StrictHostKeyChecking=no \
                                "${DOCKER_TAR_FILE}" \
                                "${SERVER_USER}@${SERVER_HOST}:/home/ubuntu/sudarshan/ride-sharing-auth/${env.BRANCH_NAME}"
                        """
                    }
                }
            }
        }

        stage('Deploy on Remote Server') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'HOST_IP', variable: 'SERVER_HOST'),
                        string(credentialsId: 'SERVER_USER', variable: 'SERVER_USER'),
                        file(credentialsId: 'SERVER_KEY', variable: 'SSH_KEY_PATH'),
                        string(credentialsId: 'SERVER_PORT', variable: 'SSH_PORT')
                    ]) {
                        // Deploy on remote server
                        sh """
                            ssh -i "${SSH_KEY_PATH}" -p "${SSH_PORT}" -o StrictHostKeyChecking=no \
                                "${SERVER_USER}@${SERVER_HOST}" << EOF
                                # Load the Docker image
                                sudo docker load -i /home/ubuntu/sudarshan/ride-sharing-auth/${env.BRANCH_NAME}/${DOCKER_TAR_FILE}
                                
                                # Stop any existing container
                                sudo docker stop ride-sharing || true
                                sudo docker rm ride-sharing || true
                                
                                # Run the new container
                                sudo docker run -d --name ride-sharing --restart unless-stopped -p 8080:8080 ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}
                                
                                # Clean up
                                rm -f /tmp/${DOCKER_TAR_FILE}
                                sudo docker image prune -f
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            script {
                def commiterEmail = sh(script: "git show -s --format='%ae'", returnStdout: true).trim()
                
                // Clean up workspace
                sh "rm -f ${DOCKER_TAR_FILE} || true"
                sh "docker rmi ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} || true"
                
                cleanWs()
            }
        }
        failure {
            script {
                def commiterEmail = sh(script: "git show -s --format='%ae'", returnStdout: true).trim()
                emailext body: '${DEFAULT_CONTENT}',
                    to: commiterEmail, 
                    subject: '${DEFAULT_SUBJECT}', 
                    saveOutput: false
            }
        }
    }
}