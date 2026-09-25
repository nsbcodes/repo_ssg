if (!(Test-Path $PSScriptRoot/src/*) ) { Write-Host -ForegroundColor Red "$PSScriptRoot/src does not exist or is empty! Place some files for testing in that directory."; exit }
Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $PSScriptRoot/out/*
$sw = [Diagnostics.Stopwatch]::StartNew()
go run $PSScriptRoot/.. $PSScriptRoot/src $PSScriptRoot/out
$success = $?
$sw.Stop()
Write-Host -ForegroundColor Green "Static Site Compilation Time: $($sw.Elapsed)"
if ($success) { python -m http.server -d $PSScriptRoot/out }