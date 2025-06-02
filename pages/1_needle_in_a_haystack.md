## Needle in a haystack

When is it more efficient to convert a _slice_ into a _map_ for locating an element `x` within the set `A` (x ∈ A)?

> **TL;DR**: use `map` when `len(haystack) > 100 && len(needles) > 100`

![needle in a haystack graph](assets/BenchmarkNeedleInAHaystack.png "Find element in set benchmark")
Depending on size of the _haystack_ (size) and number of _needles_ (iterations), this will differ:

|                                                  | assets/BenchmarkNeedleInAHaystack-map.txt |     | assets/BenchmarkNeedleInAHaystack-slice.txt |     |           |              |
| ------------------------------------------------ | ----------------------------------------- | --- | ------------------------------------------- | --- | --------- | ------------ |
|                                                  | sec/op                                    | CI  | sec/op                                      | CI  | vs base   | P            |
| NeedleInAHaystack/size=10_iterations=10-8        | 2.0195000000000003e-07                    | 0%  | 5.218e-08                                   | 0%  | -74.16%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=100-8       | 6.858000000000001e-07                     | 0%  | 3.5580000000000005e-07                      | 0%  | -48.12%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=500-8       | 2.849e-06                                 | 0%  | 1.5585e-06                                  | 0%  | -45.30%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10_iterations=1000-8      | 5.527500000000001e-06                     | 0%  | 3.475e-06                                   | 0%  | -37.13%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=10-8       | 1.1250000000000002e-06                    | 0%  | 2.54e-07                                    | 0%  | -77.42%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=100-8      | 1.6680000000000002e-06                    | 0%  | 2.1815e-06                                  | 0%  | +30.79%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=500-8      | 3.8165e-06                                | 0%  | 1.0569e-05                                  | 0%  | +176.93%  | p=0.000 n=10 |
| NeedleInAHaystack/size=100_iterations=1000-8     | 6.5995e-06                                | 0%  | 2.2282500000000002e-05                      | 0%  | +237.64%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=10-8       | 6.55e-06                                  | 0%  | 1.1410000000000002e-06                      | 1%  | -82.58%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=100-8      | 7.164500000000001e-06                     | 0%  | 1.1229500000000001e-05                      | 1%  | +56.74%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=500-8      | 9.369000000000001e-06                     | 0%  | 5.4998000000000004e-05                      | 1%  | +487.02%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500_iterations=1000-8     | 1.20535e-05                               | 0%  | 0.00011133700000000001                      | 1%  | +823.69%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=10-8      | 1.51545e-05                               | 1%  | 2.164e-06                                   | 1%  | -85.72%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=100-8     | 1.58995e-05                               | 0%  | 2.08445e-05                                 | 1%  | +31.10%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=500-8     | 1.7502999999999998e-05                    | 1%  | 0.000105146                                 | 1%  | +500.73%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000_iterations=1000-8    | 2.06655e-05                               | 0%  | 0.0002072555                                | 1%  | +902.91%  | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=10-8      | 7.94045e-05                               | 0%  | 1.01825e-05                                 | 1%  | -87.18%   | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=100-8     | 7.982400000000001e-05                     | 0%  | 0.000102442                                 | 1%  | +28.33%   | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=500-8     | 8.2366e-05                                | 0%  | 0.0005082955                                | 3%  | +517.12%  | p=0.000 n=10 |
| NeedleInAHaystack/size=5000_iterations=1000-8    | 8.38165e-05                               | 0%  | 0.0010191010000000001                       | 1%  | +1115.87% | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=10-8     | 0.00016355150000000001                    | 0%  | 2.01805e-05                                 | 1%  | -87.66%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=100-8    | 0.000164117                               | 0%  | 0.00020134900000000003                      | 1%  | +22.69%   | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=500-8    | 0.00016748800000000003                    | 0%  | 0.0010075785                                | 2%  | +501.58%  | p=0.000 n=10 |
| NeedleInAHaystack/size=10000_iterations=1000-8   | 0.00016740450000000002                    | 0%  | 0.0020248345000000003                       | 4%  | +1109.55% | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=10-8     | 0.000853139                               | 0%  | 0.00010159900000000001                      | 0%  | -88.09%   | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=100-8    | 0.000859315                               | 0%  | 0.0010305635000000001                       | 2%  | +19.93%   | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=500-8    | 0.0008550655000000001                     | 0%  | 0.005007204                                 | 4%  | +485.59%  | p=0.000 n=10 |
| NeedleInAHaystack/size=50000_iterations=1000-8   | 0.0008593165000000001                     | 0%  | 0.010244186500000002                        | 6%  | +1092.13% | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=10-8    | 0.0016933565000000002                     | 1%  | 0.00020292850000000001                      | 1%  | -88.02%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=100-8   | 0.0016954315                              | 0%  | 0.002017895                                 | 3%  | +19.02%   | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=500-8   | 0.0016948945000000001                     | 0%  | 0.0100737875                                | 9%  | +494.36%  | p=0.000 n=10 |
| NeedleInAHaystack/size=100000_iterations=1000-8  | 0.001696308                               | 0%  | 0.020756502500000003                        | 7%  | +1123.63% | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=10-8    | 0.015493154500000002                      | 3%  | 0.0010208700000000001                       | 3%  | -93.41%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=100-8   | 0.015549726000000002                      | 3%  | 0.010098937                                 | 8%  | -35.05%   | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=500-8   | 0.015741979                               | 2%  | 0.052639357                                 | 7%  | +234.39%  | p=0.000 n=10 |
| NeedleInAHaystack/size=500000_iterations=1000-8  | 0.015602500500000002                      | 4%  | 0.09976974450000001                         | 21% | +539.45%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=10-8   | 0.0418112905                              | 2%  | 0.0020738275                                | 2%  | -95.04%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=100-8  | 0.0415576295                              | 1%  | 0.020849485                                 | 8%  | -49.83%   | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=500-8  | 0.0417492095                              | 2%  | 0.0975659975                                | 12% | +133.70%  | p=0.000 n=10 |
| NeedleInAHaystack/size=1000000_iterations=1000-8 | 0.041954636000000003                      | 1%  | 0.1997639465                                | 26% | +376.14%  | p=0.000 n=10 |
| geomean                                          | 0.00015047703447270417                    |     | 0.00019927621720009707                      |     | +32.43%   |              |
