pipeline {
    agent any

    environment {
        // Docker image settings
        DOCKER_IMAGE_NAME = "ride-sharing"
        DOCKER_IMAGE_TAG = "${env.BRANCH_NAME}-${env.GIT_COMMIT_SHORT}"
        DOCKER_TAR_FILE = "${DOCKER_IMAGE_NAME}-${DOCKER_IMAGE_TAG}.tar"

        // Dynamic environment detection
        IS_PROD = (env.BRANCH_NAME == "prod")
        COMPOSE_FILE = IS_PROD ? "docker-compose-prod.yml" : "docker-compose.yml"
    }

    stages {
        stage('Checkout Code') {
            steps { checkout scm }
        }

        stage('Build Docker Image') {
            steps {
                sh """
                    docker build \
                        --build-arg APP_ENV=${IS_PROD ? 'production' : 'development'} \
                        -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} \
                        .
                """
            }
        }

        stage('Save and Transfer Docker Image') {
            steps {
                script {
                    // 1. Save Docker image as .tar file
                    sh "docker save -o ${DOCKER_TAR_FILE} ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}"

                    // 2. SCP the image to the target server
                    withCredentials([
                        string(credentialsId: "HOST_IP_${env.BRANCH_NAME.toUpperCase()}", variable: 'SERVER_HOST'),
                        string(credentialsId: "SERVER_USER_${env.BRANCH_NAME.toUpperCase()}", variable: 'SERVER_USER'),
                        file(credentialsId: "SERVER_KEY_${env.BRANCH_NAME.toUpperCase()}", variable: 'SSH_KEY_PATH'),
                        string(credentialsId: "SERVER_PORT_${env.BRANCH_NAME.toUpperCase()}", variable: 'SSH_PORT')
                    ]) {
                        sh """
                            scp -i ${SSH_KEY_PATH} -P ${SSH_PORT} \
                                -o StrictHostKeyChecking=no \
                                ${DOCKER_TAR_FILE} \
                                ${SERVER_USER}@${SERVER_HOST}:/tmp/
                        """
                    }
                }
            }
        }

        stage('Deploy on Remote Server') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: "HOST_IP_${env.BRANCH_NAME.toUpperCase()}", variable: 'SERVER_HOST'),
                        string(credentialsId: "SERVER_USER_${env.BRANCH_NAME.toUpperCase()}", variable: 'SERVER_USER'),
                        file(credentialsId: "SERVER_KEY_${env.BRANCH_NAME.toUpperCase()}", variable: 'SSH_KEY_PATH'),
                        string(credentialsId: "SERVER_PORT_${env.BRANCH_NAME.toUpperCase()}", variable: 'SSH_PORT'),
                        file(credentialsId: "ENV_FILE_${env.BRANCH_NAME.toUpperCase()}", variable: 'ENV_FILE')
                    ]) {
                        // 1. Copy compose file and .env
                        sh """
                            rsync -avz -e "ssh -i ${SSH_KEY_PATH} -P ${SSH_PORT} -o StrictHostKeyChecking=no" \
                                ${COMPOSE_FILE} .env \
                                ${SERVER_USER}@${SERVER_HOST}:/home/${SERVER_USER}/ride-sharing/
                        """

                        // 2. SSH commands to load image and start containers
                        sh """
                            ssh -i ${SSH_KEY_PATH} -p ${SSH_PORT} -o StrictHostKeyChecking=no \
                                ${SERVER_USER}@${SERVER_HOST} << 'EOF'
                                cd /home/${SERVER_USER}/ride-sharing

                                # Load the Docker image
                                docker load -i /tmp/${DOCKER_TAR_FILE}

                                # Deploy with docker-compose
                                docker-compose -f ${COMPOSE_FILE} up -d

                                # Cleanup
                                rm /tmp/${DOCKER_TAR_FILE}
                                docker image prune -af
                            EOF
                        """
                    }
                }
            }
        }
    }

    post {
        always {
            // Cleanup local .tar file
            sh "rm -f ${DOCKER_TAR_FILE} || true"
            cleanWs()
        }
        failure {
            emailext body: 'Build failed! Check Jenkins for details.',
                to: 'team@example.com',
                subject: '🚨 Deployment Failed: ${env.JOB_NAME}'
        }
    }
}