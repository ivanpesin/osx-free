# osx-free

Equivalent of 'free' command for macos.

## Example

```
$ free
              total        used        free      appmem       wired   compressed (ratio)
Mem:          8.00G       5.63G       2.37G       2.15G       1.45G   2.37G -> 908.09M (62%)
+/- Cache:                4.49G       3.51G     |mempressure:   30%, normal
Swap:         2.00G     349.50M       1.66G     | swap usage:   17%
free -m
              total        used        free      appmem       wired   compressed (ratio)
Mem:           8192        5888        2304        2323        1482   2419 -> 907 (62%)
+/- Cache:                 4720        3472     |mempressure:   30%, normal
Swap:          2048         350        1698     | swap usage:   17%
```

## Screenshot

![screenshot of free vs Activity Monitor](https://github.com/ivanpesin/osx-free/blob/master/screenshot.png?raw=true)

## Options

```
$ free --help
Usage of free:
  -b	show output in bytes
  -g	show output in gigabytes
  -h	show human-readable output
  -k	show output in kilobytes
  -m	show output in megabytes
  ```
