DEVICE_ADDR = ENV['DEVICE_ADDR']
DEPLOY_DIR  = '/home/root/'
PROJECT = 'demo'
BUILD_TARGET = "#{PROJECT}/run.go"

task :build do
  status = system("GOARCH=386 GOOS=linux go build #{BUILD_TARGET}")
  puts "Build #{status ? 'SUCCESS' : 'FAILED'}"
end

task :deploy => :build do
  puts "Deploying #{PROJECT} via scp to #{DEVICE_ADDR}..."
  status = system("scp -r run #{DEVICE_ADDR}:#{DEPLOY_DIR}/")
  File.delete('run')
  puts "Deployment #{status ? 'SUCCESS' : 'FAILED'}"
end