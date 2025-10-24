Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue
if (-not (Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue)) {
    Write-Output "Starting Docker Desktop (elevated)..."
    Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe" -Verb RunAs
    Start-Sleep -Seconds 5
}
$svc = Get-Service -Name com.docker.service -ErrorAction SilentlyContinue
if ($svc -and $svc.Status -ne 'Running') {
    Write-Output "Starting com.docker.service..."
    Start-Service -Name com.docker.service -ErrorAction SilentlyContinue
}
docker version
