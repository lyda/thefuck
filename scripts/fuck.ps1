if ((Get-Command "fuck" -ErrorAction Ignore).CommandType -eq "Function") {
	fuck @args;
	[Console]::ResetColor()
	exit
}

"First time use of thefuck detected."

if ((Get-Content $PROFILE -Raw -ErrorAction Ignore) -like "*thefuck*") {
} else {
	"  - Adding thefuck initialisation to user `$PROFILE"
	$script = "`niex `"`$(thefuck init powershell)`"";
	Write-Output $script | Add-Content $PROFILE
}

"  - Adding fuck() function to current session..."
iex "$($(thefuck init powershell).Replace("function fuck", "function global:fuck"))"

"  - Invoking fuck()`n"
fuck @args;
[Console]::ResetColor()
