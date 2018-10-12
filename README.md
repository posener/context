# context

A proof of concept implementation of **goroutine scoped context**,
proposed in [this blog post](https://posener.github.io/goroutine-scoped-context).

This library should not be used for production code.

## Usage

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
+   ctx := context.Get()
    // Use ctx...
}
```

Running the previously defined `foo`, with the context:

```diff
-foo(ctx)
+context.RunCtx(ctx, foo)
```

Invoking goroutines should be done with `context.Go` or `context.GoCtx`

Running a new goroutine with the current stored context:

```diff
-go foo()
+context.Go(foo)
```

More complected function:

```diff
-func bar(ctx context.Context, i int, s string) {
+func bar(i int, s string) {
+   ctx := context.Get()
    // Use ctx...
}
```

Should be wrapped with an empty function:

```diff
-bar(ctx, 1, "hello")
+context.RunCtx(ctx, func() { bar(1, "hello") })
```

Or fo goroutines:

```diff
-go bar(ctx, 1, "hello")
+context.GoCtx(ctx, func() { bar(1, "hello") })
```

Running a goroutine with a new context:

```go
// `ctx` is the context that we want to have in the invoked goroutine
context.GoCtx(ctx, foo)
```

`context.TODO` and should not be used anymore:

```diff
-foo(context.TODO())
+foo(context.Get())
```