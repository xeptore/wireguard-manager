package password

import "regexp"

var StrongPasswordRegExp = regexp.MustCompile(`^[A-Za-z0-9-_!\?@#\$%\^&\*+=~\\\|/"':;\.\{\}\(\)\[\],]{8,128}$`)
