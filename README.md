# osx-free

Equivalent of 'free' command for macos.

## Example

```
$ free
              total        used        free      appmem       wired   compressed (ratio)
Mem:          8.00G       7.98G      16.29M       2.92G       1.66G   5.12G -> 2.12G (58%)
+/- Cache:                6.71G       1.29G
Swap:      1024.00M     170.00M     854.00M
$ free -m
              total        used        free      appmem       wired   compressed (ratio)
Mem:           8192        8175          17        3000        1701   5225 -> 2167 (58%)
+/- Cache:                 6876        1316
Swap:          1024         170         854
```

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
