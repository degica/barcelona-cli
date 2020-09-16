package api

func ExampleDump_barcelonatoken() {

	dump([]byte("X-Barcelona-Token: abcdefg"), nil)

	// Output:
	// X-Barcelona-Token: [filtered]
}

func ExampleDump_githubtoken() {

	dump([]byte("X-Github-Token: abcdefg"), nil)

	// Output:
	// X-Github-Token: [filtered]
}

func ExampleDump_vaulttoken() {

	dump([]byte("X-Vault-Token: abcdefg"), nil)

	// Output:
	// X-Vault-Token: [filtered]
}

func ExampleDump_nottoken() {

	dump([]byte("String: abcdefg"), nil)

	// Output:
	// String: abcdefg
}
