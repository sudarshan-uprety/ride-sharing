pipeline {
    agent any

    environment {
        DOCKER_IMAGE_NAME = "app"
        DOCKER_IMAGE_TAG = "${env.BRANCH_NAME ?: 'dev'}"
        DOCKER_TAR_FILE = "${DOCKER_IMAGE_NAME}-${DOCKER_IMAGE_TAG}.tar"
    }

    stages {
        stage('Setup Environment') {
            steps {
                script {
                    envName = (env.BRANCH_NAME ?: "dev").toUpperCase()
                    IS_PROD = (env.BRANCH_NAME == "prod")
                    envFileCredentialId = "RIDE_AUTH_${envName}_ENV"
                }
            }
        }

        stage('Checkout Code') {
            steps { 
                checkout scm 
            }
        }

        stage('Build and Save Docker Image') {
            steps {
                script {
                    withCredentials([file(credentialsId: envFileCredentialId, variable: 'ENV_FILE')]) {
                        sh '''#!/bin/bash
                            # Setup environment
                            cp "$ENV_FILE" .env
                            export APP_ENV=${IS_PROD ? 'production' : 'development'}
                            
                            # Build and save image
                            docker-compose build app
                            docker save -o "$DOCKER_TAR_FILE" "$DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
                            
                            # Verify files
                            echo "Build artifacts:"
                            ls -la "$DOCKER_TAR_FILE" docker-compose.yml .env
                        '''
                    }
                }
            }
        }

        stage('Deploy to Remote Server') {
            steps {
                script {
                    withCredentials([
                        string(credentialsId: 'HOST_IP', variable: 'SERVER_HOST'),
                        string(credentialsId: 'SERVER_USER', variable: 'SERVER_USER'),
                        file(credentialsId: 'SERVER_KEY', variable: 'SSH_KEY_PATH'),
                        string(credentialsId: 'SERVER_PORT', variable: 'SSH_PORT'),
                        file(credentialsId: envFileCredentialId, variable: 'ENV_FILE')
                    ]) {
                        def branchDir = "/home/ubuntu/sudarshan/ride-sharing-auth/${env.BRANCH_NAME}/"
                        
                        sh '''#!/bin/bash
                            # Transfer files securely
                            scp -i "$SSH_KEY_PATH" -P "$SSH_PORT" -o StrictHostKeyChecking=no \\
                                "$DOCKER_TAR_FILE" \\
                                "docker-compose.yml" \\
                                "$ENV_FILE" \\
                                "$SERVER_USER@$SERVER_HOST:$branchDir"
                            
                            # Deploy commands
                            ssh -i "$SSH_KEY_PATH" -p "$SSH_PORT" -o StrictHostKeyChecking=no \\
                                "$SERVER_USER@$SERVER_HOST" << 'EOF'
                                cd "$branchDir"
                                mv $(basename "$ENV_FILE") .env
                                sudo docker load -i "$DOCKER_TAR_FILE"
                                sudo docker-compose down || true
                                sudo docker-compose up -d
                                rm -f "$DOCKER_TAR_FILE"
                                sudo docker image prune -f
                            EOF
                        '''
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
