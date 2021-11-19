package main

import (
	compute "github.com/pulumi/pulumi-azure-native/sdk/go/azure/compute"
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	resources "github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//variables
		addressPrefixVar := "10.1.0.0/16"
		adminPasswordOrKeyParam := "Password1234"
		adminUsernameParam := "MyAdminUser"
		publicIPAddressNameVar := "vmname"
		subnetNameParam := "Subnet"
		virtualNetworkNameParam := "vNet"
		subnetAddressPrefixVar := "10.1.0.0/24"
		ubuntuOSVersionParam := "18.04-LTS"
		vmNameParam := "simpleLinuxVM"
		//vmSizeParam := "Standard_B2s"
		networkInterfaceNameVar := "vmname"
		networkSecurityGroupNameParam := "SecGroupNet"

		// resources

		resourceGroupVar, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
			Location:          pulumi.String("eastus"),
			ResourceGroupName: pulumi.String("myResourceGroup"),
		})
		if err != nil {
			return err
		}

		//networkSecurityGroupResource
		networkSecurityGroupResource, err := network.NewNetworkSecurityGroup(ctx, "networkSecurityGroupResource", &network.NetworkSecurityGroupArgs{
			Location:                 resourceGroupVar.Location,
			NetworkSecurityGroupName: pulumi.String(networkSecurityGroupNameParam),
			ResourceGroupName:        resourceGroupVar.Name,
			SecurityRules: network.SecurityRuleTypeArray{
				&network.SecurityRuleTypeArgs{
					Access:                   pulumi.String("Allow"),
					DestinationAddressPrefix: pulumi.String("*"),
					DestinationPortRange:     pulumi.String("22"),
					Direction:                pulumi.String("Inbound"),
					Name:                     pulumi.String("SSH"),
					Priority:                 pulumi.Int(1000),
					Protocol:                 pulumi.String("Tcp"),
					SourceAddressPrefix:      pulumi.String("*"),
					SourcePortRange:          pulumi.String("*"),
				},
			},
		})
		if err != nil {
			return err
		}

		publicIPAddressResource, err := network.NewPublicIPAddress(ctx, "publicIPAddressResource", &network.PublicIPAddressArgs{
			DnsSettings: &network.PublicIPAddressDnsSettingsArgs{
				DomainNameLabel: pulumi.String("simplevmexample-anirudh"),
			},
			IdleTimeoutInMinutes:     pulumi.Int(4),
			Location:                 resourceGroupVar.Location,
			PublicIPAddressVersion:   pulumi.String("IPv4"),
			PublicIPAllocationMethod: pulumi.String("Dynamic"),
			PublicIpAddressName:      pulumi.String(publicIPAddressNameVar),
			ResourceGroupName:        resourceGroupVar.Name,
			Sku: &network.PublicIPAddressSkuArgs{
				Name: pulumi.String("Basic"),
			},
		})
		if err != nil {
			return err
		}
		//virtualNetworkResource
		virtualNetworkResource, err := network.NewVirtualNetwork(ctx, "virtualNetworkResource", &network.VirtualNetworkArgs{
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String(addressPrefixVar),
				},
			},
			Location:           resourceGroupVar.Location,
			ResourceGroupName:  resourceGroupVar.Name,
			VirtualNetworkName: pulumi.String(virtualNetworkNameParam),
		})
		if err != nil {
			return err
		}
		//subnetResource
		subnetResource, err := network.NewSubnet(ctx, "subnetResource", &network.SubnetArgs{
			AddressPrefix:                     pulumi.String(subnetAddressPrefixVar),
			PrivateEndpointNetworkPolicies:    pulumi.String("Enabled"),
			PrivateLinkServiceNetworkPolicies: pulumi.String("Enabled"),
			ResourceGroupName:                 resourceGroupVar.Name,
			SubnetName:                        pulumi.String(subnetNameParam),
			VirtualNetworkName:                virtualNetworkResource.Name,
		})
		if err != nil {
			return err
		}
		//networkInterfaceResource
		networkInterfaceResource, err := network.NewNetworkInterface(ctx, "networkInterfaceResource", &network.NetworkInterfaceArgs{
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				&network.NetworkInterfaceIPConfigurationArgs{
					Name:                      pulumi.String("ipconfig1"),
					PrivateIPAllocationMethod: pulumi.String("Dynamic"),
					PublicIPAddress: &network.PublicIPAddressTypeArgs{
						Id: publicIPAddressResource.ID(),
					},
					Subnet: &network.SubnetTypeArgs{
						Id: subnetResource.ID(),
					},
				},
			},
			Location:             resourceGroupVar.Location,
			NetworkInterfaceName: pulumi.String(networkInterfaceNameVar),
			NetworkSecurityGroup: &network.NetworkSecurityGroupTypeArgs{
				Id: networkSecurityGroupResource.ID(),
			},
			ResourceGroupName: resourceGroupVar.Name,
		})
		if err != nil {
			return err
		}

		Vmdetails, err := compute.NewVirtualMachine(ctx, "virtualMachine", &compute.VirtualMachineArgs{
			HardwareProfile: &compute.HardwareProfileArgs{
				VmSize: pulumi.String("Standard_D2s_v3"),
			},

			ResourceGroupName: resourceGroupVar.Name,
			Location:          resourceGroupVar.Location,
			NetworkProfile: &compute.NetworkProfileArgs{
				NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
					&compute.NetworkInterfaceReferenceArgs{
						Id: networkInterfaceResource.ID(),
					},
				},
			},
			OsProfile: &compute.OSProfileArgs{
				AdminPassword: pulumi.String(adminPasswordOrKeyParam),
				AdminUsername: pulumi.String(adminUsernameParam),
				ComputerName:  pulumi.String(vmNameParam),
			},
			StorageProfile: &compute.StorageProfileArgs{
				ImageReference: &compute.ImageReferenceArgs{
					Offer:     pulumi.String("UbuntuServer"),
					Publisher: pulumi.String("Canonical"),
					Sku:       pulumi.String(ubuntuOSVersionParam),
					Version:   pulumi.String("latest"),
				},
				OsDisk: &compute.OSDiskArgs{
					CreateOption: pulumi.String("FromImage"),
					DiskSizeGB:   pulumi.Int(30),
				},
			},
			VmName: pulumi.String(vmNameParam),
		})
		if err != nil {
			return err
		}
		ctx.Export("ip", Vmdetails.HardwareProfile)
		return nil
	})
}
