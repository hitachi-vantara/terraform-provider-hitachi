resource "hitachi_vosb_compute_node" "mycompute" {
  vosb_address = ""
  compute_node_name = "ComputeNode-RESTAPI123"
  #name = "ComputeNode-RESTAPI-new"
  os_type = "VMware"
  # os_type = "Windows"
  # iscsi_connection {
  #   iscsi_initiator = "iqn.1998-01.com.vmware:node-06-0723aa94"
  #   port_names      = ["002-iSCSI-001"]
  #   #port_names = ["002-iSCSI-001","001-iSCSI-002"]
  # }

  # iscsi_connection {
  #   iscsi_initiator = "iqn.1998-01.com.vmware:node-06-0723aa94"
  #   port_names      = ["002-iSCSI-001", "001-iSCSI-002"]
  # }

  # iscsi_connection {
  #   iscsi_initiator = "iqn.1998-01.com.vmware:node-06-0723aa95"
  #   port_names      = ["002-iSCSI-001"]
  # }

  fc_connection {
    host_wwn = "60060e8107595326"

  }

  fc_connection {
    host_wwn = "90060e8107595325"
  }

  fc_connection {
    host_wwn = "90060e8107599325"
  }
}

output "computenodecreate" {
  value = resource.hitachi_vosb_compute_node.mycompute
}
