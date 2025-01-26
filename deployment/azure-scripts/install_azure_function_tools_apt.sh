wget -q https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb
apt install ./packages-microsoft-prod.deb
rm ./packages-microsoft-prod.deb
apt update
apt install azure-functions-core-tools-4

# See also: https://github.com/Azure/azure-functions-core-tools#linux