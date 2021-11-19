import * as resources from "@pulumi/azure-native/resources";
import * as network from "@pulumi/azure-native/network";
import * as compute from "@pulumi/azure-native/compute";
import * as pulumi from "@pulumi/pulumi";



// Create an Azure Resource Group
const resourceGroup = new resources.ResourceGroup("resourceGroup_ts",{
    resourceGroupName:"myResourcegroup_ts",
    location : "westus"

});

// Network Security Group

const networkSecurityGroup = new network.NetworkSecurityGroup("networkSecurityGroup", {
    location: resourceGroup.location,
    networkSecurityGroupName: "testnsg",
    resourceGroupName: resourceGroup.name,
    securityRules: [{
        access: "Allow",
        destinationAddressPrefix: "*",
        destinationPortRange: "22",
        direction: "Inbound",
        name: "SSH",
        priority: 1000,
        protocol: "Tcp",
        sourceAddressPrefix: "*",
        sourcePortRange: "*",
    }],
});

// Public IP Address Resource

const publicIPAddress_res = new network.PublicIPAddress("publicIPAddress", {
    dnsSettings: {
        domainNameLabel: "dnslbl-anirudh-ts",
    },
    location: resourceGroup.location,
    publicIpAddressName: "test-ip",
    resourceGroupName: resourceGroup.name,
    idleTimeoutInMinutes:4,
    publicIPAddressVersion:"IPv4",
    publicIPAllocationMethod:"Dynamic",
});

// Virtual Network Resource

const virtualNetwork = new network.VirtualNetwork("virtualNetwork", {
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    location: resourceGroup.location,
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: "test-vnet",
});

// Subnet Resource 

const subnet = new network.Subnet("subnet", {
    addressPrefix: "10.0.0.0/16",
    privateEndpointNetworkPolicies : "Enabled",
    privateLinkServiceNetworkPolicies : "Enabled",
    resourceGroupName: resourceGroup.name,
    subnetName: "subnet1",
    virtualNetworkName: virtualNetwork.name,
});


// Network Interface Resource 

const networkInterface = new network.NetworkInterface("networkInterface", {
    enableAcceleratedNetworking: true,
    ipConfigurations: [{
        name: "ipconfig1",
        publicIPAddress: {
            id:publicIPAddress_res.id,
        },
        subnet: {
            id: subnet.id,
        },
    }],
    location: resourceGroup.location,
    networkInterfaceName: "test-nic",
    resourceGroupName: resourceGroup.name,
    networkSecurityGroup : {
        id : networkSecurityGroup.id
    }
        
});

// Virtual Machine 


const virtualMachine = new compute.VirtualMachine("virtualMachine", {
    hardwareProfile: {
        vmSize: "Standard_D2s_v3",
    },
    location: resourceGroup.location,
    networkProfile: {
        networkInterfaces: [{
            id: networkInterface.id,
            primary: true,
        }],
    },
    osProfile: {
        adminPassword: "Unif!12#",
        adminUsername: "admin123",
        computerName: "myVM",
        linuxConfiguration: {
            patchSettings: {
                assessmentMode: "ImageDefault",
            },
            provisionVMAgent: true,
        },
    },
    resourceGroupName: resourceGroup.name,
    storageProfile: {
        imageReference: {
            offer: "UbuntuServer",
            publisher: "Canonical",
            sku: "16.04-LTS",
            version: "latest",
        },
        osDisk: {
            caching: "ReadWrite",
            createOption: "FromImage",
            managedDisk: {
                storageAccountType: "Premium_LRS",
            },
            name: "myVMosdisk",
        },
    },
    vmName: "myVM",
});


// Print Public IP 

const examplePublicIP = pulumi.all([publicIPAddress_res.name, resourceGroup.name]).apply(([name, resourceGroupName]) => network.getPublicIPAddress({
    publicIpAddressName: name,
    resourceGroupName: resourceGroupName,
}));
export const publicIpAddress = examplePublicIP.ipAddress;



// const exportipadd = network.getPublicIPAddress({
//     publicIpAddressName : "test-ip",
//     resourceGroupName: "myResourcegroup_ts",
// });

// export const publicIpAddress = exportipadd.then(example => example.ipAddress);

// const done = pulumi.all({ : virtualMachine.id, name: publicIPAddress_res.name, resourceGroup: resourceGroup });

// export const ipAddress = done.apply(async (d) => {
//     return network.getPublicIPAddress({
//         resourceGroupName: d.resourceGroupName,
//         publicIpAddressName: d.name,
//     }).then(ip => ip.ipAddress);
// });
