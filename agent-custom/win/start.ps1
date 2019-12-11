New-Item "\azp\agent" -ItemType directory | Out-Null

# Let the agent ignore the token env variables
$Env:VSO_AGENT_IGNORE = "AZP_TOKEN,AZP_TOKEN_FILE"

Set-Location agent
if (-not (Test-Path Env:AZP_URL)) {
 $Env:AZP_DOWNLOAD_URL = "$(cat \vsts\agent\.url)"
}


Write-Host $Env:AZP_DOWNLOAD_URL

$packageUrl = $Env:AZP_DOWNLOAD_URL

Write-Host "Downloading and installing Azure Pipelines agent..." -ForegroundColor Cyan

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($packageUrl, "$(Get-Location)\agent.zip")

Expand-Archive -Path "agent.zip" -DestinationPath "\azp\agent"
try
{
  Copy-Item \vsts\agent\.agent -Destination \azp\agent\.agent
  Copy-Item \vsts\agent\.credentials -Destination \azp\agent\.credentials
  Write-Host "Running Azure Pipelines agent..." -ForegroundColor Cyan

  .\run.cmd
}
catch
{
    write-host "Caught an exception:" -ForegroundColor Red
    write-host "Exception Type: $($_.Exception.GetType().FullName)" -ForegroundColor Red
    write-host "Exception Message: $($_.Exception.Message)" -ForegroundColor Red
}
finally
{
  Write-Host "Cleanup. Removing Azure Pipelines agent..." -ForegroundColor Cyan

  cat C:\azp\agent\_diag\*.log
}