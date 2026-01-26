pipeline {
    agent any

    environment {
        APP_NAME = "gator"
        BUILD_DIR = "build"
        SECURITY_DIR = "security"
        SONAR_COMPOSE = "docker-compose.sonarqube.yml"
    }

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Verify formatting') {
            steps {
                sh '''
                if [ -n "$(gofmt -l .)" ]; then
                  echo "gofmt violations found"
                  gofmt -l .
                  exit 1
                fi
                '''
            }
        }

        stage('Static analysis') {
            steps {
                sh '''
                go install honnef.co/go/tools/cmd/staticcheck@latest
                staticcheck ./...
                '''
            }
        }

        stage('Unit tests') {
            steps {
                sh '''
                go test -coverprofile=coverage.out ./...
                '''
            }
        }

        stage('Start SonarQube (if needed)') {
            steps {
                sh '''
                if ! docker ps --format '{{.Names}}' | grep -q '^sonarqube$'; then
                  docker compose -f ${SONAR_COMPOSE} up -d
                fi
                '''
            }
        }

        stage('SonarQube scan') {
            steps {
                withSonarQubeEnv('sonarqube') {
                    sh '''
                    sonar-scanner \
                      -Dsonar.projectKey=${APP_NAME} \
                      -Dsonar.projectName=${APP_NAME}
                    '''
                }
            }
        }

        stage('Quality gate') {
            steps {
                timeout(time: 5, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Build binary') {
            steps {
                sh '''
                mkdir -p ${BUILD_DIR}
                CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
                  go build -o ${BUILD_DIR}/${APP_NAME}
                '''
            }
        }

        stage('Build Docker image') {
            steps {
                sh '''
                docker build -t ${APP_NAME}:${BUILD_NUMBER} .
                '''
            }
        }

        stage('Generate SBOMs') {
            steps {
                sh '''
                mkdir -p ${SECURITY_DIR}

                syft dir:. -o cyclonedx-json > ${SECURITY_DIR}/sbom-source.json
                syft file:${BUILD_DIR}/${APP_NAME} -o cyclonedx-json > ${SECURITY_DIR}/sbom-binary.json
                syft ${APP_NAME}:${BUILD_NUMBER} -o cyclonedx-json > ${SECURITY_DIR}/sbom-image.json
                '''
            }
        }

        stage('Vulnerability scan (CRA policy)') {
            steps {
                sh '''
                grype sbom:${SECURITY_DIR}/sbom-image.json \
                  --fail-on critical \
                  -o json > ${SECURITY_DIR}/vuln-report.json
                '''
            }
        }
    }

    post {

        always {
            archiveArtifacts artifacts: '''
                ${BUILD_DIR}/*
                ${SECURITY_DIR}/*.json
                coverage.out
            ''', fingerprint: true

            sh '''
            docker compose -f ${SONAR_COMPOSE} down || true
            '''
        }
    }
}

