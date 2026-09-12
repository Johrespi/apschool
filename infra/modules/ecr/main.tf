resource "aws_ecr_repository" "this" {
  name                 = "${var.project_name}-${var.environment}-backend"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.this.name
  policy = jsonencode({
    rules : [{
      rulePriority : 1
      selection : {
        tagStatus : "any"
        countType : "imageCountMoreThan"
        countNumber : 3
      }
      action : { type : "expire" }
  }] })
}
