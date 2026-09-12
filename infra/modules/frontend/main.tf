resource "aws_s3_bucket" "this" {
  bucket_prefix = "${var.project_name}-${var.environment}-frontend-"
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket                  = aws_s3_bucket.this.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

}

resource "aws_cloudfront_origin_access_control" "this" {
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
  name                              = "${var.project_name}-${var.environment}-oac"
  origin_access_control_origin_type = "s3"
}

resource "aws_cloudfront_distribution" "this" {
  enabled             = true
  default_root_object = "index.html"
  origin {
    origin_id                = "S3-apschool-frontend"
    domain_name              = aws_s3_bucket.this.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.this.id
  }

  origin {
    origin_id   = "ALB-backend"
    domain_name = var.alb_dns_name
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }
  default_cache_behavior {
    target_origin_id       = "S3-apschool-frontend"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    cache_policy_id        = "658327ea-f89d-4fab-a63d-7e88639e58f6"

  }

  ordered_cache_behavior {
    path_pattern             = "/api/*"
    target_origin_id         = "ALB-backend"
    viewer_protocol_policy   = "redirect-to-https"
    allowed_methods          = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods           = ["GET", "HEAD"]
    cache_policy_id          = "4135ea2d-6df8-44a3-9df3-4b5a84be39ad"
    origin_request_policy_id = "216adef6-5c7f-47e4-b989-5492eafa07d3"
  }

  custom_error_response {
    error_code            = 403
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 10
  }

  custom_error_response {
    error_code            = 404
    response_code         = 200
    response_page_path    = "/index.html"
    error_caching_min_ttl = 10
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
}

resource "aws_s3_bucket_policy" "this" {
  bucket = aws_s3_bucket.this.id
  policy = jsonencode({
    Version : "2012-10-17"
    Statement : [{
      Effect : "Allow"
      Action : "s3:GetObject"
      Principal : { Service : "cloudfront.amazonaws.com" }
      Resource : "${aws_s3_bucket.this.arn}/*"
      Condition : {
        StringEquals : {
          "AWS:SourceArn" : "${aws_cloudfront_distribution.this.arn}"
        }
      }
      }
    ]
  })
}
