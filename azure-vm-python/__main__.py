"""An Azure RM Python Pulumi program"""

from arpeggio.cleanpeg import sequence
import pulumi
from pulumi_azure_native import storage
from pulumi_azure_native import resources
from pulumi_azure_native import network
from pulumi_azure_native import compute
from pulumi_azure_native.network.get_public_ip_address import get_public_ip_address
from pulumi_azure_native.network.network_security_group import NetworkSecurityGroup
from pulumi_azure_native.network.public_ip_address import PublicIPAddress
from pulumi_azure_native.network.security_rule import SecurityRule

# Create an Azure Resource Group
resource_group = resources.ResourceGroup("resourceGroup_py",
    location="eastus")

# network resouce group
network_security_group = network.NetworkSecurityGroup("networkSecurityGroup",
    location=resource_group.location,
    network_security_group_name="testnsg",
    resource_group_name = resource_group.name,
    security_rules=[network.SecurityRuleArgs(
        access="Allow",
        destination_address_prefix="*",
        destination_port_range="22",
        direction="Inbound",
        name="SSH",
        priority=1000,
        protocol="Tcp",
        source_address_prefix="*",
        source_port_range="*"
        )]

    )
    
    

# publicIPAddressResource 
public_ip_address_res = network.PublicIPAddress("publicIPAddress",
    dns_settings = network.PublicIPAddressDnsSettingsArgs(
        domain_name_label="dnslbl-anirudh",
    ),
    location=resource_group.location,
    public_ip_address_name="test-ip",
    resource_group_name = resource_group.name,
    idle_timeout_in_minutes=4,
    public_ip_address_version="IPv4",
    public_ip_allocation_method="Dynamic",
    sku= network.PublicIPAddressSkuArgs(
        name="Basic"
    )
    )

# virtualNetworkResource
virtual_network = network.VirtualNetwork("virtualNetwork",
    address_space = network.AddressSpaceArgs(
        address_prefixes=["10.0.0.0/16"],
    ),
    virtual_network_name="test-vnet",
    location=resource_group.location,
    resource_group_name = resource_group.name)

# subnetResource
subnet = network.Subnet("subnet",
    address_prefix="10.0.0.0/16",
    private_endpoint_network_policies="Enabled",
    private_link_service_network_policies="Enabled",
    resource_group_name=resource_group.name,
    subnet_name="subnet1",
    virtual_network_name= virtual_network.name)

# networkInterfaceResource
network_interface = network.NetworkInterface("networkInterface",
    ip_configurations=[network.NetworkInterfaceIPConfigurationArgs(
        name="ipconfig1",
        public_ip_address=network.PublicIPAddressArgs(
            id=public_ip_address_res.id,
        ),
        subnet=network.SubnetArgs(
            id=subnet.id,
        ),
    )],
    location=resource_group.location,
    network_interface_name="test-nic",
    network_security_group= network.NetworkSecurityGroupArgs(
        id=network_security_group.id
    ) ,
    resource_group_name=resource_group.name)

# virtal machine 

virtual_machine = compute.VirtualMachine("virtualMachine",
    hardware_profile=compute.HardwareProfileArgs(
        vm_size="Standard_D2s_v3",
    ),
    location=resource_group.location,
    network_profile=compute.NetworkProfileArgs(
        network_interfaces=[compute.NetworkInterfaceReferenceArgs(
            id=network_interface.id,
            primary=True,
        )],
    ),
    os_profile=compute.OSProfileArgs(
        admin_password="Unif!12#",
        admin_username="admin123",
        computer_name="myVM",
        linux_configuration= compute.LinuxConfigurationArgs(
            patch_settings= compute.LinuxPatchSettingsArgs(
                assessment_mode="ImageDefault",
            ),
            provision_vm_agent=True,
        ),
    ),
    resource_group_name=resource_group.name,
    storage_profile= compute.StorageProfileArgs(
        image_reference= compute.ImageReferenceArgs(
            offer="UbuntuServer",
            publisher="Canonical",
            sku="16.04-LTS",
            version="latest",
        ),
        os_disk= compute.OSDiskArgs(
            caching="ReadWrite",
            create_option="FromImage",
            managed_disk=compute.ManagedDiskParametersArgs(
                storage_account_type="Premium_LRS",
            ),
            name="myVMosdisk",
        ),
    ),
    vm_name="myVM")


    
# example_public_ip = pulumi.Output.all(public_ip_address_res.name, resource_group.name).apply(lambda name, resource_group_name: network.get_public_ip(name=name,
#             resource_group_name=resource_group_name))
# pulumi.export("publicIpAddress", example_public_ip.ip_address)


# publicIPofmachine = get_public_ip_address(
#     public_ip_address_name=public_ip_address_res.name,
#     resource_group_name=resource_group.name
# )
# pulumi.export("Public IP of the machine:" , publicIPofmachine.ip_address)
# pulumi.export("Domain Name Label:", publicIPofmachine.dns_settings.domain_name_label)


