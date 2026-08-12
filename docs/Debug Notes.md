# How to dump headers to the debug console

```
dump, err := httputil.DumpRequestOut(req, false)
	if err == nil {
		fmt.Printf("\n--- DEBUG OUTGOING HEADERS ---\n%s\n-------------------------\n", string(dump))
	}
```
