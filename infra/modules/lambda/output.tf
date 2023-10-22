output "invoke_arn" {
  value = aws_lambda_function.main.invoke_arn
}

output "name" {
  value = aws_lambda_function.main.function_name
}
