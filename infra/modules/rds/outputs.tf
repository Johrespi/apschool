output "db_endpoint" {
  value = aws_db_instance.db.endpoint
}

output "db_address" {
  value = aws_db_instance.db.address
}


output "db_port" {
  value = aws_db_instance.db.port
}


output "db_name" {
  value = aws_db_instance.db.db_name
}


output "security_group_id" {
  value = aws_security_group.db_sg.id
}


