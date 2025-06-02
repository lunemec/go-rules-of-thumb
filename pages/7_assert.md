## Negative space programming (a.k.a. asserts everywhere)

What is the cost of adding `assert`? Does it make any significant impact?

> **TL;DR**: Use asserts whenever possible to improve reliability of your software. The cost is almost non-existent.

| Type         | N statements | ns/op        |     |
| ------------ | ------------ | ------------ | --- |
| no assert    | 1            | 0.3453 ns/op |
| assert       | 1            | 0.4979 ns/op |
| assert       | 5            | 1.791 ns/op  |
| defer assert | 1            | 2.411 ns/op  |

Read more:

- <https://spinroot.com/gerard/pdf/P10.pdf>
- <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety>
- <https://alfasin.com/2017/12/21/negative-space-and-how-does-it-apply-to-coding/>
