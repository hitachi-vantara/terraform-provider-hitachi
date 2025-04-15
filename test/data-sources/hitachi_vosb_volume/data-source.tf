data "hitachi_vosb_volume" "vosbvolumes" {
  vosb_address = ""
  volume_name       = "Mongonode3_vol4dd"

}

output "volumeoutput" {
  value = data.hitachi_vosb_volume.vosbvolumes
}
