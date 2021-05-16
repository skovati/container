# gocon
Extremely simple implementation of a rootless Linux container for the educational purpose of learning about namespaces and cgroups.

Inspired and heavily derived from [Liz Rice's talk](https://www.youtube.com/watch?v=_TsSmSu57Zo)

Includes an [Alpine Linux](https://www.alpinelinux.org/) rootfs to use.

## installation
```bash
git clone https://github.com/skovati/gocon
cd gocon
make install
```

## usage
```bash
gocon run echo hello from alpine!
```

```bash
gocon run /bin/sh
```
