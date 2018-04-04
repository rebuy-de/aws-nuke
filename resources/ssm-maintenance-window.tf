resource "aws_ssm_maintenance_window" "window" {
  name     = "AWSNukeTest"
  schedule = "cron(0 16 ? * TUE *)"
  duration = 3
  cutoff   = 1
}
