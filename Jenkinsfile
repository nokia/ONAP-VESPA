properties([gitLabConnection('gitlabe1')])
pipeline {
    agent {dockerfile true}
    stages {
        stage("Prepare") {
            steps {
                updateGitlabCommitStatus(name: "Jenkins-Build", state: "running")
                sh "go version"
                sh "make tools"
                sh "make clean"
                sh "make sloc"
            }   
        }
        stage('Build') {
            environment {
                // Environment variables for makefile
                BRANCH = "${env.GIT_BRANCH}"
                BUILDID = "${env.BUILD_NUMBER}"
                VERSION = "${env.GIT_BRANCH.replaceAll(~/-+/, ".")}"
            }
            steps {
                sh "make build"
            }
        }
        stage('Static analysis') {
            steps {
                sh "make lint"
            }
        }
        stage('Test') {
            steps {
                sh "make test"
            }
        }
        stage('Package') {
            environment {
                // Environment variables for makefile
                BRANCH = "${env.GIT_BRANCH}"
                BUILDID = "${env.BUILD_NUMBER}"
                VERSION = "${env.GIT_BRANCH.replaceAll(~/-+/, ".")}"
            }
            steps {
                sh "make rpm"
                zip zipFile: "build/ves-tools-win.zip", dir: "build/windows"
                zip zipFile: "build/ves-tools-linux.zip", dir: "build/linux"
            }
        }
        
        stage('Publish') {
            when { // Publish to artifactory only when running on master branch
                anyOf {
                    branch 'master'
                    // tag pattern: "\\d+\\.\\d+\\.\\d+", comparator: "REGEXP"
                    expression { env.TAG_NAME ==~ /\d+\.\d+\.\d+/ }
                }
            }
            steps {
                script {
                    def server = Artifactory.server "artifactory-espoo1"
                    def uploadSpec = """{
                    "files": [
                        {
                        "pattern": "build/*.rpm",
                        "target": "sdmexpert-snapshots-local/ves-agent/"
                        }
                    ]
                    }"""
                    def buildInfo = server.upload(uploadSpec)
                    buildInfo.retention maxBuilds: 40, deleteBuildArtifacts: true
                    // buildInfo.env.collect()
                    server.publishBuildInfo(buildInfo)

                    def promotionConfig = [
                        // Mandatory parameters
                        'buildName'          : buildInfo.name,
                        'buildNumber'        : buildInfo.number,
                        'targetRepo'         : 'sdmexpert-local',
                    
                        // Optional parameters
                        'comment'            : 'Product validated by SyVe',
                        'sourceRepo'         : 'sdmexpert-snapshots-local',
                        'status'             : 'Released',
                        'includeDependencies': true,
                        'copy'               : false,
                        // 'failFast' is true by default.
                        // Set it to false, if you don't want the promotion to abort upon receiving the first error.
                        'failFast'           : true
                    ]
                    Artifactory.addInteractivePromotion server: server, promotionConfig: promotionConfig, displayName: "Promote this build"
                }
            }
        }
    }
    post {
        always {
            junit "build/testresults.xml"
            checkstyle canComputeNew: false, defaultEncoding: '', healthy: '', pattern: 'build/checkstyle.xml', unHealthy: ''
            cobertura coberturaReportFile: 'build/coverage.xml'
            sloccountPublish encoding: '', ignoreBuildFailure: true, pattern: 'build/sloccount.scc'
            archive "build/*.rpm,build/ves-tools-*.zip"
            deleteDir() // Delete workspace
        }
        aborted {
            updateGitlabCommitStatus(name: "Jenkins-Build", state: "canceled")
        }
        failure {
            updateGitlabCommitStatus(name: "Jenkins-Build", state: "failed")
        }
        success {
            updateGitlabCommitStatus(name: "Jenkins-Build", state: "success")
        }
        unstable {
            updateGitlabCommitStatus(name: "Jenkins-Build", state: "failed")
        }
    }
}