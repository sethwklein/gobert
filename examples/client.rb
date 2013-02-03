require 'bertrpc'

# http://github.com/mojombo/bert
# http://github.com/mojombo/bertrpc

svc = BERTRPC::Service.new('localhost', 8000)
puts "Sending request..."
result = svc.call.calc.fib(42)
puts "Received #{result}"
