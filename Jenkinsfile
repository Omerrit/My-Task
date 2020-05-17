def output
def errmsg

pipeline {
    agent any
    environment {
    SERVER = 'gerrit.lan:31123'
    CONTAINER_IMAGE = 'gerrit.lan:31123/go'
    IMAGE_TAG = "$CONTAINER_IMAGE"+":$GIT_COMMIT"
    IMAGE_LATEST_TAG = "$CONTAINER_IMAGE"+":latest"
    IMAGE_TESTING_TAG = "$CONTAINER_IMAGE"+":testing"
  }
    stages {
        stage('Build image') {
		when {
         	branch '**'
          	}
            steps {
                echo "Builder Started ..."
		sh 'ls'
            script {
                def outfile = 'stdout.out'
                def stderrfile = 'stderr.out'
                def status = sh(
                script: "docker build --no-cache --tag=$IMAGE_TAG . >${outfile} 2>${stderrfile}",
                returnStatus: true
                )
                output = readFile(outfile).trim()
                errmsg = readFile(stderrfile).trim()
                if (status == 0) {
                println("STDOUT: - ${output}")
                } else {
                println("STDOUT: - ${output}")
                println("ERROR: - ${errmsg}")
                }
        }
        }
        
    
    post {
        success { 
            sh 'docker tag $IMAGE_TAG $IMAGE_TESTING_TAG'
            gerritReview score:1 }
	    unstable { gerritReview labels: [Verified: 0], message: 'Build is unstable' }
        failure {
//            sh 'docker rm `docker ps -aq -f status=exited` '
//            sh 'docker rmi `docker image ls -q -f dangling=true` '
            gerritReview labels: [Verified: -1], message: "ERROR: - ${errmsg}"

        }
        }
    }
    
	
    stage('Push image to Registry') {
		when {
		branch 'master'
		}
	    steps {
		echo "uploading container to registry ..."
		sh 'docker login -u truenameless -p Tr21^a58G21bLx7bR@7Grl54^ $SERVER'
		sh 'docker tag $IMAGE_TAG $IMAGE_LATEST_TAG'
		sh 'docker push $IMAGE_TAG'
		sh 'docker push $IMAGE_LATEST_TAG'
		sh 'docker logout $SERVER'
		sh 'docker image ls $CONTAINER_IMAGE'
		}
	}
}
}
