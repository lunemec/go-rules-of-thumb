## Pass by reference vs copy

When should you pass a reference (pointer), and when should you use pass by value?

> [!TIP]
> use pointers when you need mutation  
> for read-only data, start with the simpler API and measure  
> this benchmark does not show a reliable universal size cutoff

In this benchmark, passing a pointer wins for all tested struct sizes, and the gap grows as the copied array gets larger. That is still a narrow microbenchmark, so the safe rule is not "always use pointers", but "measure once copying large values shows up in a profile".

References (pointers) vs copied values are still way more complicated than one synthetic test can capture, and there is tons of resources on this topic. A great one is
[this article](https://dave.cheney.net/2017/04/29/there-is-no-pass-by-reference-in-go) by Dave Cheney.
