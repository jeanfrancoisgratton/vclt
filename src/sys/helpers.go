// vclt
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/sys/helpers.go
// Creation date: 2025/03/21 18:33

package sys

// Helper function to ensure that the token is globally set for further use
func _setAuthTkn(tkn string) {
	AuthToken = tkn
}

// Helper function to advertise that we are logged in
func _setLoggedIn() {
	IsLoggedIn = true
}
