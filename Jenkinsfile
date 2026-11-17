// Jenkinsfile - Producto 2 (FP.050 DevOps y cloud computing)
// Cada push a GitHub lanza este pipeline:
//   1) Construye la imagen Docker de la web
//   2) La carga en Minikube
//   3) La despliega en Kubernetes con varias réplicas y un balanceador
//
// Rama main    -> entorno "produccion" (3 réplicas, puerto 30080)
// Rama develop -> entorno "pruebas"    (2 réplicas, puerto 30081)

pipeline {
    agent any

    environment {
        APP = 'web-uoc'
    }

    stages {
        stage('Entorno') {
            steps {
                script {
                    def rama = env.GIT_BRANCH ?: env.JOB_NAME
                    if (rama.contains('main') || rama.contains('produccion')) {
                        env.NS = 'produccion'
                        env.REPLICAS = '3'
                        env.PUERTO = '30080'
                    } else {
                        env.NS = 'pruebas'
                        env.REPLICAS = '2'
                        env.PUERTO = '30081'
                    }
                    env.TAG = "${env.NS}-${env.BUILD_NUMBER}"
                    echo "Rama: ${rama} -> entorno ${env.NS}"
                }
            }
        }

        stage('Build') {
            steps {
                sh 'docker build -t $APP:$TAG .'
            }
        }

        stage('Cargar imagen en Minikube') {
            steps {
                sh 'minikube image load $APP:$TAG'
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    kubectl create namespace $NS --dry-run=client -o yaml | kubectl apply -f -

                    sed -e "s|IMAGEN|$APP:$TAG|" \
                        -e "s|REPLICAS|$REPLICAS|" \
                        -e "s|ENTORNO_APP|$NS|" \
                        k8s/deployment.yaml | kubectl apply -n $NS -f -

                    sed "s|PUERTO_NODO|$PUERTO|" k8s/service.yaml | kubectl apply -n $NS -f -

                    kubectl rollout status deployment/$APP -n $NS --timeout=180s
                    kubectl get pods,svc -n $NS -o wide
                '''
            }
        }
    }
}
