# context

Package context is a proof of concept implementation of **scoped context**,
proposed in [this blog post](https://posener.github.io/goroutine-scoped-context).

This library should not be used for production code.

#### Usage

The context package should be imported from `github.com/posener/context`.

```diff
 import (
-   "context"
+   "github.com/posener/context"
 )
```

Since this implementation does not involve changes to the runtime,
the goroutine context must be initialized.

```diff
 func main() {
+    context.Init()
 	// Go code goes here.
 }
```

Functions should not anymore receive the context in the first argument.
They should get it from the goroutine scope.

```diff
-func foo(ctx context.Context) {
+func foo() {
+	ctx := context.Get()
 	// Use context.
 }
```

Applying context to a scope:

```go
unset := context.Set(ctx)
// ctx is applied until unset is called, or a deeper `Set` call.
unset()
```

Or:

```go
defer context.Set(ctx)()
// ctx is applied until the end of the function or a deeper `Set` call.
```

Invoking goroutines should be done with `context.Go` or `context.GoCtx`

Running a new goroutine with the current stored context:

```diff
-go foo()
+context.Go(foo)
```

More complected functions:

```diff
-go foo(1, "hello")
+context.Go(func() { foo(1, "hello") })
```

Running a goroutine with a new context:

```go
// `ctx` is the context that we want to have in the invoked goroutine
context.GoCtx(ctx, foo)
```

`context.TODO` should not be used anymore:

```diff
-f(context.TODO())
+f(context.Get())```

## Sub Packages

* [runtime](./runtime)


---

Created by [goreadme](https://github.com/apps/goreadme)
