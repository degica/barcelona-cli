package api

func Example_dump_barcelonatoken() {

	dump([]byte("X-Barcelona-Token: abcdefg"), nil)

	// Output:
	// X-Barcelona-Token: [filtered]
}

func Example_dump_githubtoken() {

	dump([]byte("X-Github-Token: abcdefg"), nil)

	// Output:
	// X-Github-Token: [filtered]
}

func Example_dump_vaulttoken() {

	dump([]byte("X-Vault-Token: abcdefg"), nil)

	// Output:
	// X-Vault-Token: [filtered]
}

func Example_dump_nottoken() {

	dump([]byte("String: abcdefg"), nil)

	// Output:
	// String: abcdefg
}
