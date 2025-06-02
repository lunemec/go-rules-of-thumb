## Pass by reference vs copy

When should you pass a reference (pointer), and when should you use pass by value?

> [!HINT]
> pass by reference if you want to mutate the data, otherwise pass a copy

Performance-wise, this one is almost impossible to give general advice for. If your struct (or nested structs)
are very big (it depends on the types of fields too), copying will become slower.
But if you have many more pointers, you increase GC pressure and your program will
spend more time on waiting on memory pointer lookup.

References (pointers) vs copied values is way more complicated,
and there is tons of resources on this topic, great one is
[this article](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go) by Dave Cheney.
