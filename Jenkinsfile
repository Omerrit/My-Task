pipeline {
    agent any
    environment {
    SERVER = 'gerrit.lan:31123'
    CONTAINER_IMAGE = 'gerrit.lan:31123/go'
    IMAGE_TAG = "$CONTAINER_IMAGE"+":$GIT_COMMIT"
    IMAGE_LATEST_TAG = "$CONTAINER_IMAGE"+":latest"
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
                   try {
		sh 'docker build --no-cache --tag=$IMAGE_TAG .'
        } catch (err) {
                       echo "exception caught, going on"
                       println err
            }
        }
		sh 'docker tag $IMAGE_TAG $IMAGE_LATEST_TAG'
            }
        
    
    post {
        success { gerritReview score:1 }
	    unstable { gerritReview labels: [Verified: 0], message: 'Build is unstable' }
        failure { 
            echo "BUILD UNSTABLE"
            sh 'docker rm `docker ps -aq -f status=exited` '
            sh 'docker rmi `docker image ls -q -f dangling=true` '
            gerritReview score:-1, message: 
            println err
        
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
