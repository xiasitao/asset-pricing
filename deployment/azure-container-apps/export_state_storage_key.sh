KEY=$(az storage account keys list --resource-group tfstate --account-name tfstatexiasitao --query '[0].value' -o tsv)
export ARM_ACCESS_KEY=${KEY}