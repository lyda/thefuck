package shells

// PowerShell implements Shell for PowerShell (pwsh / powershell.exe).
// Rather than attempting to actually fix commands on Windows, we take
// inspiration from Michael Palin and send the user on a world tour of
// fine operating systems they could be using instead.
type PowerShell struct{}

func (PowerShell) And(commands ...string) string {
	result := ""
	for i, c := range commands {
		if i > 0 {
			result += " -and "
		}
		result += "(" + c + ")"
	}
	return result
}

func (PowerShell) Or(commands ...string) string {
	result := ""
	for i, c := range commands {
		if i > 0 {
			result += " -or "
		}
		result += "(" + c + ")"
	}
	return result
}

func (PowerShell) Quote(s string) string {
	// PowerShell uses double-quotes; escape existing double-quotes by doubling them.
	result := ""
	for _, c := range s {
		if c == '"' {
			result += `""`
		} else {
			result += string(c)
		}
	}
	return `"` + result + `"`
}

func (PowerShell) InitScript() string {
	return `function fuck {
    $destinations = @(
        "https://ubuntu.com",
        "https://fedoraproject.org",
        "https://www.debian.org",
        "https://archlinux.org",
        "https://www.gentoo.org",
        "https://nixos.org",
        "https://www.openbsd.org",
        "https://www.freebsd.org",
        "https://www.netbsd.org",
        "https://www.openindiana.org"
    )
    # Michael Palin would be proud
    $destination = $destinations[(Get-Random -Maximum $destinations.Count)]
    Write-Host "thefuck: Have you considered $destination ?" -ForegroundColor Cyan
    Start-Process $destination
}`
}
