# This script is used for creating new agent pool configured as kubernetes poolprovider

$URI = Read-Host 'Enter URI of the account to be configured for poolprovider'
$PATToken = Read-Host 'Enter PAT Token'
$poolname = Read-Host 'Enter poolname'
$dnsname = Read-Host 'Enter dnsname'
$sharedSecret = Read-Host 'Enter shared secret (atleast 16 characters required)'
$targetSize = Read-Host 'Enter targetSize'

if($URI -eq $null -or $URI -eq "")
{
	Write-Error "Provide a valid URI"
}
elseif($poolname -eq $null -or $poolname -eq "")
{
	Write-Error "Provide a valid name for Pool"
}
elseif($dnsname -eq $null -or $dnsname -eq "")
{
	Write-Error "Provide a valid value for dns name"
}
elseif($sharedSecret -eq $null -or $sharedSecret -eq "" -or $sharedSecret.length -lt 16)
{
	Write-Error "Provide a valid value for sharedSecret"
}
elseif($targetSize -eq $null -or $targetSize -eq "")
{
	Write-Error "Provide a valid value for targetSize"
}
elseif($PATToken -eq $null -or $PATToken -eq "")
{
	Write-Error "Provide a valid value for PAT token"
}
else
{
	$body = @{
      "name"="$poolname"
      "type"="Ignore"
      "acquireAgentEndpoint"="https://$($dnsname)/acquire"
      "releaseAgentEndpoint"="https://$($dnsname)/release"
      "sharedSecret"="$sharedSecret"
    } | ConvertTo-Json

    $header = @{
     "Accept"="application/json"
	 "Authorization"="Basic "+[System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes(":$($PATToken)" ))
     "Content-Type"="application/json"
    } 

    
    $response = Invoke-WebRequest -Uri "$URI/_apis/distributedtask/agentclouds?api-version=5.0-preview" -Method 'Post' -Body $body -Headers $header
    Write-Host $response
    $jsonObj = ConvertFrom-Json $([String]::new($response.Content))

    Write-Host $jsonObj.agentCloudId

    $body1 = @{
      "name"="$poolname"
      "agentCloudId"=$jsonObj.agentCloudId
	  "targetSize"=$targetSize
    } | ConvertTo-Json

    $response1 = Invoke-WebRequest -Uri "$URI/_apis/distributedtask/pools?api-version=5.0-preview" -Method 'Post' -Body $body1 -Headers $header
    Write-Host $response1
    $jsonObj1 = ConvertFrom-Json $([String]::new($response1.Content))
}