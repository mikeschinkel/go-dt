package dtx

// Noop is useful during development when you have yet used a variablebut want to
// run the debugger to see that value the variable will contain. Since Go does
// not compile if a variable is not used, passing it to dtx.Noop() will "use" it.
func Noop(args ...any) {}
