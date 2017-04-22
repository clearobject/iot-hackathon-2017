DEPLOY_DIR  = '/home/root/'
PROJECT = 'event'
BUILD_TARGET = 'run.go'

task :build do
  status = system("GOARCH=386 GOOS=linux go build #{BUILD_TARGET}")
  puts "Build #{status ? 'SUCCESS' : 'FAILED'}"
end

task :deploy => :build do
  addresses = ENV['DEVICE_ADDRS'].split(',')
  addresses.each do |addr|
    puts "Deploying #{PROJECT} via scp to #{addr}..."
    `scp -r run #{addr}:#{DEPLOY_DIR}/`
  end
  File.delete('run')
end