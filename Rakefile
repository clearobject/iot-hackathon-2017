DEPLOY_DIR  = '/home/root/'
PROJECT = 'event'
BUILD_TARGET = 'reflector.go'

task :build do
  status = system("GOARCH=386 GOOS=linux go build #{BUILD_TARGET}")
  puts "Build #{status ? 'SUCCESS' : 'FAILED'}"
end

task :deploy => :build do
  puts `ansible-playbook deploy/playbook.yml -i deploy/inventory`
end