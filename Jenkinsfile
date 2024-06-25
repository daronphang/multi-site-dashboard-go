pipeline {
    agent none
    options {
        disableConcurrentBuilds()
        skipStagesAfterUnstable()
        buildDiscarder(logRotator(numToKeepStr:'10'))
    }

    environment {
        // Branches.
        DEV_BRANCH = "dev"
        PROD_BRANCH = "main"

        // Application.
        PROJECT_NAME = "dashboard"
        SERVICE_NAME = "multi-site-dashboard-go"
        LANGUAGE = 'go' // angular, dotnet, nodejs, go, python
        K8S_MANIFESTS_URL = "http://tools.mf.platform/bitbucket/scm/dash/kubernetes-manifests.git"

        // Image.
        IMAGE_BUILD_TAG_PREFIX = "harbor.mf.platform/$PROJECT_NAME/$SERVICE_NAME"

        // Credentials.
        BITBUCKET_CREDENTIALS_ID = "bitbucket-daronphang"
        HARBOR_CREDENTIALS_ID = "harbor-dashboard"
        NEXUS_CREDENTIALS = credentials("bitbucket-tanggy")
        SONAR_LOGIN_HASH = "sqp_8677e78b37c03fc427db576acee53db78d1fe364"

        // Secret files to be injected by Jenkins when deploying to k8s.
        // Comment out if not required, else it will run the creating secrets stage.
        // Filename should have the environment as suffix e.g. orderapp-secret-env-dev, orderapp-secret-env-prod
        // SECRET_FILE_ENV_PREFIX = "$SERVICE_NAME-secret-env"

        // For Angular.
        BASEHREF = "/dashboard/"

        // DAST ZAP.
        // for backend api, give path to .json file eg http://dev.mf.platform/<SERVICE_NAME>/swagger/v1/swagger.json
        ZAP_TARGET_URL = "http://dev.mf.platform/dashboard" 
    }
    stages {
        stage('Continuous Integration') {
            agent { label 'k8s-development' }
            environment {
                ENVIRONMENT = "test"
                IMAGE_BUILD_TAG_LATEST = "$IMAGE_BUILD_TAG_PREFIX:t$BUILD_NUMBER"
                SONAR_URL = "http://tools.mf.platform/sonar/"
            }
            when {
                beforeAgent true
                expression {
                    return !(env.BRANCH_NAME == env.PROD_BRANCH || env.BRANCH_NAME == env.DEV_BRANCH)
                }
            }
            stages {
                stage('Test Image Build - ANGULAR') {
                    when {
                        environment(name: 'LANGUAGE', value: 'angular')
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build --build-arg CONFIG=testing \
                        --build-arg BASEHREF=$BASEHREF \
                        --build-arg DEPLOYMENT_IMAGE=build \
                        -t $IMAGE_BUILD_TAG_LATEST .
                        '''
                    }
                }

                stage('Test Image Build - OTHERS') {
                    when {
                        not {
                            environment(name: 'LANGUAGE', value: 'angular')
                        }
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build -t $IMAGE_BUILD_TAG_LATEST .
                        '''
                    }
                }

                stage('Run Tests in Image') {
                    steps {
                        // Customize the steps in this stage as necessary.
                        sh '''
                        chmod 766 $WORKSPACE
                        docker run -d -v $WORKSPACE:/coverage:rw \
                        --name test-$SERVICE_NAME \
                        $IMAGE_BUILD_TAG_LATEST \
                        sleep 5m
                        docker exec test-$SERVICE_NAME \
                        sh -c "npm run test:ci && cp /dashboard/coverage/lcov.info /coverage/lcov.info"
                        '''
                    }
                    post {
                        always {
                            sh '''
                            docker rm -f test-$SERVICE_NAME
                            chmod 755 $WORKSPACE
                            '''
                        }
                        failure {
                            sh 'docker rmi $IMAGE_BUILD_TAG_LATEST'
                        }
                    }
                }

                // stage('Image Vulnerability Scan- Anchore') {
                //     steps {
                //         sh 'echo "$IMAGE_BUILD_TAG_LATEST $WORKSPACE/Dockerfile " > anchore_images'
                //         anchore name: 'anchore_images', forceAnalyze: true,  bailOnFail: true, bailOnPluginFail: true
                //     }
                // }

                stage('SAST - SonarQube DOTNET') {
                    when {
                        environment(name: 'LANGUAGE', value: 'dotnet')
                    }
                    steps {
                        withSonarQubeEnv('SonarQube') {
                            sh '''
                            sudo dotnet sonarscanner begin /k:"${SERVICE_NAME}" \
                            /d:sonar.host.url="${SONAR_URL}" \
                            /d:sonar.login="${SONAR_LOGIN_HASH}"
                            sudo dotnet build
                            sudo dotnet sonarscanner end /d:sonar.login="${SONAR_LOGIN_HASH}"
                            '''
                        }
                    }
                    post {
                        always {
                            sh 'sudo rm -rf .sonarqube'
                        }
                    }
                }

                stage('SAST - SonarQube OTHERS') {
                    when {
                        not {
                            environment(name: 'LANGUAGE', value: 'dotnet')
                        }
                    }
                    steps {
                        withSonarQubeEnv('SonarQube') {
                            sh '''
                            sudo sonar-scanner -Dsonar.projectKey=${SERVICE_NAME} \
                            -Dsonar.projectName=${SERVICE_NAME} \
                            -Dsonar.host.url=${SONAR_URL} \
                            -Dsonar.login=${SONAR_LOGIN_HASH} \
                            -Dsonar.typescript.lcov.reportPaths=$WORKSPACE/lcov.info
                            '''
                        }
                    }
                    post {
                        always {
                            sh 'sudo rm -rf .sonarqube'
                        }
                    }
                }

                stage('Quality Gate') {
                    steps {
                        timeout(time: 1, unit: 'HOURS') {
                            // Parameter indicates whether to set pipeline to UNSTABLE if Quality Gate fails
                            // true = set pipeline to UNSTABLE, false = don't
                            waitForQualityGate abortPipeline: false
                        }
                    }
                }
            }
        }

        stage('Continuous Deployment - Development') {
            agent { label 'k8s-development' }
            environment {
                ENVIRONMENT = 'dev'
                IMAGE_BUILD_TAG = "$IMAGE_BUILD_TAG_PREFIX:d$BUILD_NUMBER"
                IMAGE_BUILD_TAG_LATEST = "$IMAGE_BUILD_TAG_PREFIX:$ENVIRONMENT"
                ZAP_REPORT_FILENAME_PREFIX = "zap-report-$SERVICE_NAME-$IMAGE_VERSION_NUMBER"
            }
            when {
                beforeAgent true
                branch "$DEV_BRANCH"
            }
            stages {
                stage('Image Build - ANGULAR') {
                    when {
                        environment(name: 'LANGUAGE', value: 'angular')
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build --build-arg CONFIG=development \
                        --build-arg BASEHREF=$BASEHREF \
                        -t $IMAGE_BUILD_TAG .
                        '''
                    }
                }

                stage('Image Build - OTHERS') {
                    when {
                        not {
                            environment(name: 'LANGUAGE', value: 'angular')
                        }
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build -t $IMAGE_BUILD_TAG .
                        '''
                    }
                }

                stage('Image Push') {
                    steps {
                        withCredentials([
                            usernamePassword(
                                credentialsId: "$HARBOR_CREDENTIALS_ID",
                                usernameVariable: 'HARBOR_USERNAME',
                                passwordVariable: 'HARBOR_PASSWORD'
                            ),
                        ]) {
                            sh '''
                            docker login harbor.mf.platform -u $HARBOR_USERNAME -p $HARBOR_PASSWORD
                            docker tag $IMAGE_BUILD_TAG $IMAGE_BUILD_TAG_LATEST
                            docker push $IMAGE_BUILD_TAG
                            docker push $IMAGE_BUILD_TAG_LATEST
                            '''
                        }
                    }
                    post {
                        always {
                            sh '''
                            docker rmi $IMAGE_BUILD_TAG
                            docker rmi $IMAGE_BUILD_TAG_LATEST
                            '''
                        }
                    }
                }

                stage('Create Secrets in k8s-dev Cluster') {
                    when {
                        expression {
                            return env.SECRET_FILE_ENV_PREFIX
                        }
                    }
                    steps {
                        withCredentials([
                            file(credentialsId: "$SECRET_FILE_ENV_PREFIX-$ENVIRONMENT", variable: 'SECRET_FILE')
                        ]) {
                            sh '''
                            chmod 744 /$SECRET_FILE
                            kubectl delete secret $SECRET_FILE_ENV_PREFIX -n $PROJECT_NAME
                            kubectl create secret generic $SECRET_FILE_ENV_PREFIX --from-env-file=/$SECRET_FILE -n $PROJECT_NAME
                            '''
                        }
                    }
                }

                stage('Deploy to k8s-dev Cluster') {
                    steps {
                        git branch: "$PROD_BRANCH", credentialsId: "$BITBUCKET_CREDENTIALS_ID", url: "$K8S_MANIFESTS_URL"
                        sh '''
                        kubectl apply -f $SERVICE_NAME/kubernetes-$ENVIRONMENT.yaml -n $PROJECT_NAME
                        kubectl set image deployment/$SERVICE_NAME $SERVICE_NAME=$IMAGE_BUILD_TAG -n $PROJECT_NAME
                        '''
                    }
                }

                stage('[DAST] Run OWASP ZAP Analysis - ANGULAR') {
                    when {
                        allOf {
                            expression {
                                return env.ZAP_TARGET_URL
                            };
                            environment(name: 'LANGUAGE', value: 'angular')
                        }
                    }
                    steps {
                        sh '''
                        chmod 777 $WORKSPACE
                        docker run --name owasp-$SERVICE_NAME-$BUILD_NUMBER \
                        -v $WORKSPACE:/zap/wrk/:rw \
                        -t softwaresecurityproject/zap-stable \
                        zap-baseline.py -t $ZAP_TARGET_URL \
                        -g gen.conf \
                        -r $ZAP_REPORT_FILENAME_PREFIX.html \
                        -J $ZAP_REPORT_FILENAME_PREFIX.json \
                        -I
                        '''
                    }
                    post {
                        always {    
                            sh '''
                            docker rm owasp-$SERVICE_NAME-$BUILD_NUMBER
                            chmod 755 $WORKSPACE
                            '''
                        }
                    }            
                }

                stage('[DAST] Run OWASP ZAP Analysis - OTHERS') {
                    when {
                        allOf {
                            expression {
                                return env.ZAP_TARGET_URL
                            };
                            not {
                                environment(name: 'LANGUAGE', value: 'angular')
                            }
                        }
                    }
                    steps {
                        sh '''
                        chmod 777 $WORKSPACE
                        docker run --name owasp-$SERVICE_NAME-$BUILD_NUMBER \
                        -v $WORKSPACE:/zap/wrk/:rw \
                        -t softwaresecurityproject/zap-stable \
                        zap-api-scan.py -t $ZAP_TARGET_URL \
                        -f openapi \
                        -g backend.conf \
                        -r $ZAP_REPORT_FILENAME_PREFIX.html \
                        -J $ZAP_REPORT_FILENAME_PREFIX.json \
                        -I
                        '''
                    }
                    post {
                        always {
                            sh '''
                            docker rm owasp-$SERVICE_NAME-$BUILD_NUMBER
                            chmod 755 $WORKSPACE
                            '''
                        }
                    }            
                }

                stage('[DAST] Publish OWASP ZAP Report') {
                    steps {
                        sh '''
                        sudo curl -v -u $NEXUS_CREDENTIALS --upload-file $WORKSPACE/$ZAP_REPORT_FILENAME_PREFIX.html https://tools.mf.platform/nexus/repository/$PROJECT_NAME/$SERVICE_NAME/$ZAP_REPORT_FILENAME_PREFIX.html
                        sudo curl -v -u $NEXUS_CREDENTIALS --upload-file $WORKSPACE/$ZAP_REPORT_FILENAME_PREFIX.json https://tools.mf.platform/nexus/repository/$PROJECT_NAME/$SERVICE_NAME/$ZAP_REPORT_FILENAME_PREFIX.json
                        '''
                        publishHTML([
                            allowMissing: false,
                            alwaysLinkToLastBuild: false,
                            keepAll: true,
                            reportDir: '',
                            reportFiles: '$ZAP_REPORT_FILENAME_PREFIX.html',
                            reportName: 'OWASP HTML Report',
                            useWrapperFileDirectly: true,
                        ]) 
                    }
                    post {
                        always {         
                            sh'''
                            sudo rm $WORKSPACE/$ZAP_REPORT_FILENAME_PREFIX.html
                            sudo rm $WORKSPACE/$ZAP_REPORT_FILENAME_PREFIX.json
                            '''
                        }
                    }            
                }
            }
            post {
                failure {
                    cleanWs(
                        cleanWhenNotBuilt: false,
                        deleteDirs: true,
                        notFailBuild: true
                    )
                }
            }
        }

        stage('Continuous Deployment - Production') {
            agent { label 'ntx-production' }
            environment {
                ENVIRONMENT = 'prod'
                IMAGE_BUILD_TAG = "$IMAGE_BUILD_TAG_PREFIX:p$BUILD_NUMBER"
                IMAGE_BUILD_TAG_LATEST = "$IMAGE_BUILD_TAG_PREFIX:$ENVIRONMENT"
            }
            when {
                beforeAgent true
                branch "$PROD_BRANCH"
            }
            stages {
                stage('Image Build - ANGULAR') {
                    when {
                        environment(name: 'LANGUAGE', value: 'angular')
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build --build-arg CONFIG=production \
                        --build-arg BASEHREF=$BASEHREF \
                        -t $IMAGE_BUILD_TAG .
                        '''
                    }
                }

                stage('Image Build - OTHERS') {
                    when {
                        not {
                            environment(name: 'LANGUAGE', value: 'angular')
                        }
                    }
                    steps {
                        // Customize the build arguments as necessary.
                        sh '''
                        docker build -t $IMAGE_BUILD_TAG .
                        '''
                    }
                }

                stage('Image Push') {
                    steps {
                        withCredentials([
                            usernamePassword(credentialsId: "$HARBOR_CREDENTIALS_ID", usernameVariable: 'HARBOR_USERNAME', passwordVariable: 'HARBOR_PASSWORD'),
                        ]) {
                            sh '''
                            docker login harbor.mf.platform -u $HARBOR_USERNAME -p $HARBOR_PASSWORD
                            docker tag $IMAGE_BUILD_TAG $IMAGE_BUILD_TAG_LATEST
                            docker push $IMAGE_BUILD_TAG
                            docker push $IMAGE_BUILD_TAG_LATEST
                            '''
                        }
                    }

                    post {
                        always {
                            sh '''
                            docker rmi $IMAGE_BUILD_TAG
                            docker rmi $IMAGE_BUILD_TAG_LATEST
                            '''
                        }
                    }
                }

                stage('Create Secrets in ntx-prod Cluster') {
                    when {
                        expression {
                            return env.SECRET_FILE_ENV_PREFIX
                        }
                    }
                    steps {
                        withCredentials([
                            file(credentialsId: "$SECRET_FILE_ENV_PREFIX-$ENVIRONMENT", variable: 'SECRET_FILE')
                        ]) {
                            sh '''
                            chmod 744 /$SECRET_FILE
                            kubectl delete secret $SECRET_FILE_ENV_PREFIX -n $PROJECT_NAME
                            kubectl create secret generic $SECRET_FILE_ENV_PREFIX --from-env-file=/$SECRET_FILE -n $PROJECT_NAME
                            '''
                        }
                    }
                }

                stage('Deploy to ntx-prod Cluster') {
                    steps {
                        git branch: "$PROD_BRANCH", credentialsId: "$BITBUCKET_CREDENTIALS_ID", url: "$K8S_MANIFESTS_URL"
                        sh '''
                        kubectl apply -f $SERVICE_NAME/kubernetes-$ENVIRONMENT.yaml -n $PROJECT_NAME
                        kubectl set image deployment/$SERVICE_NAME $SERVICE_NAME=$IMAGE_BUILD_TAG -n $PROJECT_NAME
                        '''
                    }
                }
            }
            post {
                // Clean after build
                failure {
                    cleanWs(
                        cleanWhenNotBuilt: false,
                        deleteDirs: true,
                        notFailBuild: true
                    )
                }
            }
        }
    }
}
