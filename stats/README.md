The .txt files here are meant to be analyzed with the go benchstat tool.
As of v1.0.0,
that tool generates the following output:
```
goos: linux
goarch: amd64
pkg: github.com/rsnorthscope/tattle
cpu: Intel(R) Core(TM) i5-4570 CPU @ 3.20GHz
                      │    v0.1.3    │               v1.0.0                │
                      │    sec/op    │    sec/op     vs base               │
01CallOverhead-4        0.3155n ± 4%   0.2965n ± 6%  -6.02% (p=0.001 n=30)
02TatErrorsOnly-4        2.240n ± 0%    2.241n ± 0%       ~ (p=0.706 n=30)
03AddTaintCheck-4        2.238n ± 0%    2.240n ± 0%       ~ (p=0.224 n=30)
04StandardTemplate-4     7.517n ± 1%    7.533n ± 0%       ~ (p=0.923 n=30)
05TemplateLogVsLogf-4    5.590n ± 0%    5.604n ± 0%  +0.24% (p=0.021 n=30)
10HandCoded-4            4.830n ± 1%    4.824n ± 0%       ~ (p=0.796 n=30)
geomean                  2.617n         2.592n       -0.96%
```
