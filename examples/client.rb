require 'bertrpc'

svc = BERTRPC::Service.new('localhost', 8000)
result = svc.call.calc.fib(42)
puts "=> #{result}"
