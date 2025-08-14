# -------------------------------------------------------------------------
# ⚠️  WARNING: This file is maintained internally by the Hitachi Vantara team.
# ⚠️  Do NOT modify this file manually. Changes may be overwritten.
#
# 📘 This file is used to display the user consent message during Terraform runs.
#     It checks whether the user has provided consent and shows instructions if not.
# -------------------------------------------------------------------------


locals {
  consent_file_path  = "/opt/hitachi/terraform/user_consent.json"
  config_file_path   = "/opt/hitachi/terraform/bin/.internal_config"
  user_consent_given = can(file(local.consent_file_path))

  config_json = (!local.user_consent_given && can(file(local.config_file_path))) ? jsondecode(file(local.config_file_path)) : {}

  effective_consent_message = local.user_consent_given ? "" : (
    try(
      " \n${tostring(local.config_json.user_consent_message)}${tostring(local.config_json.run_consent_message)}\n\n ",
      " \n⚠️ ${tostring(local.config_json.run_consent_message)}\n\n "
    )
  )
}

output "consent_reminder" {
  value = local.effective_consent_message != "" ? local.effective_consent_message : null
}
